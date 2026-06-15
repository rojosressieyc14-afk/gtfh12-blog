package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"blog/server/internal/model"
	"gorm.io/gorm"
)

type DeepSeekConfig interface {
	GetDeepSeekKey() string
	GetDeepSeekURL() string
}

type InterviewService struct {
	db      *gorm.DB
	cfg     DeepSeekConfig
	apiKeyS *ApiKeyService
}

func NewInterviewService(db *gorm.DB, cfg DeepSeekConfig, apiKeyS *ApiKeyService) *InterviewService {
	return &InterviewService{db: db, cfg: cfg, apiKeyS: apiKeyS}
}

type deepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type deepSeekRequest struct {
	Model       string            `json:"model"`
	Messages    []deepSeekMessage `json:"messages"`
	Temperature float64           `json:"temperature,omitempty"`
	MaxTokens   int               `json:"max_tokens,omitempty"`
	Stream      bool              `json:"stream"`
}

type deepSeekChoice struct {
	Message deepSeekMessage `json:"message"`
}

type deepSeekResponse struct {
	Choices []deepSeekChoice `json:"choices"`
}

func (s *InterviewService) StartSession(userID uint, position, resumeText string, totalQuestions int, apiKeyID *uint) (*model.InterviewSession, error) {
	if totalQuestions < 1 {
		totalQuestions = 5
	}
	if totalQuestions > 20 {
		totalQuestions = 20
	}

	if apiKeyID != nil && *apiKeyID > 0 {
		if _, _, _, err := s.apiKeyS.GetDecryptedKey(*apiKeyID, userID); err != nil {
			return nil, fmt.Errorf("API Key 无效: %w", err)
		}
	}

	session := model.InterviewSession{
		UserID:         userID,
		Position:       position,
		ResumeText:     resumeText,
		TotalQuestions: totalQuestions,
		Status:         model.InterviewStatusInProgress,
		ApiKeyID:       apiKeyID,
	}
	if err := s.db.Create(&session).Error; err != nil {
		return nil, err
	}

	question, err := s.generateFirstQuestion(session)
	if err != nil {
		s.db.Delete(&session)
		return nil, err
	}

	round := model.InterviewRound{
		SessionID:   session.ID,
		RoundNumber: 1,
		Question:    question,
	}
	if err := s.db.Create(&round).Error; err != nil {
		s.db.Delete(&session)
		return nil, err
	}

	session.Rounds = []model.InterviewRound{round}
	return &session, nil
}

func (s *InterviewService) SubmitAnswer(sessionID, userID uint, answer string) (*model.InterviewSession, bool, error) {
	var session model.InterviewSession
	if err := s.db.Preload("Rounds", func(db *gorm.DB) *gorm.DB {
		return db.Order("round_number asc")
	}).First(&session, sessionID).Error; err != nil {
		return nil, false, err
	}
	if session.UserID != userID {
		return nil, false, fmt.Errorf("无权操作该面试")
	}
	if session.Status != model.InterviewStatusInProgress {
		return nil, false, fmt.Errorf("面试已结束")
	}

	currentRound := &session.Rounds[len(session.Rounds)-1]
	if currentRound.Answer != "" {
		return nil, false, fmt.Errorf("当前问题已经回答过了")
	}

	currentRound.Answer = answer
	currentRound.UpdatedAt = time.Now()

	score, feedback, err := s.evaluateAnswer(session, *currentRound)
	if err == nil {
		currentRound.Score = score
		currentRound.Feedback = feedback
	}
	s.db.Save(currentRound)

	if len(session.Rounds) >= session.TotalQuestions {
		session.Status = model.InterviewStatusCompleted
		session.UpdatedAt = time.Now()
		s.db.Save(&session)
		return &session, true, nil
	}

	nextQuestion, err := s.generateNextQuestion(session)
	if err != nil {
		return nil, false, err
	}

	nextRound := model.InterviewRound{
		SessionID:   session.ID,
		RoundNumber: len(session.Rounds) + 1,
		Question:    nextQuestion,
	}
	if err := s.db.Create(&nextRound).Error; err != nil {
		return nil, false, err
	}

	session.Rounds = append(session.Rounds, nextRound)
	return &session, false, nil
}

func (s *InterviewService) GetSession(sessionID, userID uint) (*model.InterviewSession, error) {
	var session model.InterviewSession
	if err := s.db.Preload("Rounds", func(db *gorm.DB) *gorm.DB {
		return db.Order("round_number asc")
	}).First(&session, sessionID).Error; err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return nil, fmt.Errorf("无权查看该面试")
	}
	return &session, nil
}

func (s *InterviewService) EndSession(sessionID, userID uint) (*model.InterviewSession, error) {
	var session model.InterviewSession
	if err := s.db.Preload("Rounds", func(db *gorm.DB) *gorm.DB {
		return db.Order("round_number asc")
	}).First(&session, sessionID).Error; err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return nil, fmt.Errorf("无权操作该面试")
	}
	if session.Status == model.InterviewStatusCompleted {
		return &session, nil
	}

	session.Status = model.InterviewStatusCompleted
	session.UpdatedAt = time.Now()
	s.db.Save(&session)
	return &session, nil
}

func (s *InterviewService) buildSystemPrompt(session model.InterviewSession) string {
	var b strings.Builder
	b.WriteString("你是一个专业的面试官。请严格按照以下规则进行面试：\n\n")
	b.WriteString(fmt.Sprintf("## 应聘职位\n%s\n\n", session.Position))

	if session.ResumeText != "" {
		b.WriteString(fmt.Sprintf("## 候选人简历\n%s\n\n", session.ResumeText))
	}

	b.WriteString(`## 面试规则
1. 每次只说一个问题，不要一次性问多个问题
2. 问题要结合实际工作场景，考察候选人的真实能力
3. 涵盖：技术深度、项目经验、问题解决能力、沟通表达
4. 语气专业、友好、有建设性
5. 根据候选人的回答调整后续问题的难度和方向
6. 如果候选人回答优秀，可以深入追问；如果回答一般，可以给提示

## 输出格式
- 只输出面试问题本身，不要加"问题1："之类的前缀
- 不要输出评分或评价（评分在后续步骤进行）
`)
	return b.String()
}

func (s *InterviewService) buildEvaluationPrompt(session model.InterviewSession, round model.InterviewRound) string {
	return fmt.Sprintf(`你是一个专业的面试评估官。请对以下面试回答进行评分和反馈。

## 应聘职位
%s

## 面试问题
%s

## 候选人的回答
%s

请从以下维度评估：
1. 回答的准确性和完整性（0-100分）
2. 沟通表达能力
3. 是否有亮点或待改进之处

## 输出格式（JSON，不要加markdown代码块标记）
{"score": <0-100的整数>, "feedback": "<2-4句中文反馈>"}`, session.Position, round.Question, round.Answer)
}

func (s *InterviewService) buildConversationHistory(session model.InterviewSession) []deepSeekMessage {
	messages := []deepSeekMessage{
		{Role: "system", Content: s.buildSystemPrompt(session)},
	}

	for _, round := range session.Rounds {
		if round.Answer == "" {
			continue
		}
		messages = append(messages, deepSeekMessage{Role: "assistant", Content: round.Question})
		messages = append(messages, deepSeekMessage{Role: "user", Content: round.Answer})
	}

	return messages
}

func (s *InterviewService) resolveCredentials(session model.InterviewSession) (apiKey, baseURL string, err error) {
	if session.ApiKeyID != nil && *session.ApiKeyID > 0 {
		key, _, url, e := s.apiKeyS.GetDecryptedKey(*session.ApiKeyID, session.UserID)
		if e != nil {
			return "", "", e
		}
		return key, url, nil
	}
	return s.cfg.GetDeepSeekKey(), s.cfg.GetDeepSeekURL(), nil
}

func (s *InterviewService) generateFirstQuestion(session model.InterviewSession) (string, error) {
	messages := []deepSeekMessage{
		{Role: "system", Content: s.buildSystemPrompt(session)},
		{Role: "user", Content: fmt.Sprintf("候选人正在面试「%s」职位。请开始面试，问第一个问题。", session.Position)},
	}
	apiKey, baseURL, err := s.resolveCredentials(session)
	if err != nil {
		return "", err
	}
	return s.callDeepSeek(messages, apiKey, baseURL)
}

func (s *InterviewService) generateNextQuestion(session model.InterviewSession) (string, error) {
	messages := s.buildConversationHistory(session)

	messages = append(messages, deepSeekMessage{Role: "user", Content: fmt.Sprintf(
		"基于上面的对话历史，继续面试。这是第 %d 个问题（共 %d 题）。请根据候选人的上一个回答，提出下一个问题。",
		len(session.Rounds)+1, session.TotalQuestions,
	)})
	apiKey, baseURL, err := s.resolveCredentials(session)
	if err != nil {
		return "", err
	}
	return s.callDeepSeek(messages, apiKey, baseURL)
}

func (s *InterviewService) evaluateAnswer(session model.InterviewSession, round model.InterviewRound) (int, string, error) {
	prompt := s.buildEvaluationPrompt(session, round)
	messages := []deepSeekMessage{
		{Role: "system", Content: "你是一个严格的面试评分官。只输出JSON，不要额外文字。"},
		{Role: "user", Content: prompt},
	}

	apiKey, baseURL, err := s.resolveCredentials(session)
	if err != nil {
		return 0, "", err
	}
	result, err := s.callDeepSeek(messages, apiKey, baseURL)
	if err != nil {
		return 0, "", err
	}

	result = cleanJSONResponse(result)
	var eval struct {
		Score    int    `json:"score"`
		Feedback string `json:"feedback"`
	}
	if err := json.Unmarshal([]byte(result), &eval); err != nil {
		return 0, "", err
	}
	if eval.Score < 0 {
		eval.Score = 0
	}
	if eval.Score > 100 {
		eval.Score = 100
	}
	return eval.Score, eval.Feedback, nil
}

func (s *InterviewService) callDeepSeek(messages []deepSeekMessage, apiKey, baseURL string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("API Key 未配置")
	}

	baseURL = strings.TrimRight(baseURL, "/")
	reqBody := deepSeekRequest{
		Model:       "deepseek-chat",
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   1024,
		Stream:      false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("DeepSeek 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DeepSeek 返回错误 %d: %s", resp.StatusCode, string(respBody))
	}

	var dsResp deepSeekResponse
	if err := json.Unmarshal(respBody, &dsResp); err != nil {
		return "", err
	}

	if len(dsResp.Choices) == 0 {
		return "", fmt.Errorf("DeepSeek 返回空结果")
	}

	return strings.TrimSpace(dsResp.Choices[0].Message.Content), nil
}

func cleanJSONResponse(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)
	return text
}
