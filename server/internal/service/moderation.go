package service

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"blog/server/internal/model"
	"gorm.io/gorm"
)

const (
	defaultModerationBanThreshold = 5
	moderationBanThresholdKey     = "moderation_ban_threshold"
	autoBanReasonPrefix           = "自动风控封禁"
)

var moderationThresholdState = struct {
	mu        sync.RWMutex
	threshold int
}{
	threshold: defaultModerationBanThreshold,
}

type ModerationError struct {
	Field string
	Word  string
}

func (e *ModerationError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%s 包含违禁词，已禁止提交：%s", e.Field, e.Word)
}

func IsModerationError(err error) bool {
	var target *ModerationError
	return errors.As(err, &target)
}

type moderationService struct {
	mu           sync.RWMutex
	builtinWords []string
	customWords  []string
	words        []string
}

var sensitiveChecker = newModerationService()

func newModerationService() *moderationService {
	wordSet := map[string]struct{}{}
	addWords := func(items ...string) {
		for _, item := range items {
			normalized := normalizeModerationText(item)
			if normalized == "" {
				continue
			}
			wordSet[normalized] = struct{}{}
		}
	}

	addWords(
		"色情", "黄图", "黄片", "成人视频", "成人直播", "约炮", "嫖娼", "招嫖", "卖淫", "援交", "裸聊",
		"赌博", "博彩", "赌盘", "彩票计划群", "赌场", "上分", "下分", "代下注",
		"毒品", "冰毒", "海洛因", "可卡因", "大麻", "k粉", "笑气", "致幻剂",
		"办证", "假证", "洗钱", "跑分", "电诈", "诈骗", "刷单", "套现", "代开发票", "验证码平台", "黑产",
		"枪支", "手枪", "步枪", "炸药", "爆炸物", "买凶", "雇凶", "投毒", "恐怖组织",
		"法轮功", "台独", "港独", "疆独", "藏独", "反共", "颠覆国家政权",
		"加微信", "加vx", "加wx", "点击链接", "扫码进群", "永久地址", "备用网址", "商务合作请联系",
		"porn", "adultvideo", "escort", "yuepao", "dubo", "bocai", "caipiao", "bingdu", "daima", "ddos",
	)

	words := make([]string, 0, len(wordSet))
	for word := range wordSet {
		words = append(words, word)
	}
	sort.Slice(words, func(i, j int) bool {
		if len(words[i]) == len(words[j]) {
			return words[i] < words[j]
		}
		return len(words[i]) > len(words[j])
	})

	return &moderationService{
		builtinWords: words,
		words:        append([]string(nil), words...),
	}
}

func normalizeModerationText(text string) string {
	var builder strings.Builder
	builder.Grow(len(text))

	for _, r := range strings.ToLower(strings.TrimSpace(text)) {
		r = normalizeModerationRune(r)
		switch {
		case unicode.Is(unicode.Han, r):
			builder.WriteRune(r)
		case unicode.IsLetter(r), unicode.IsDigit(r):
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

func normalizeModerationRune(r rune) rune {
	if r >= 0xFF01 && r <= 0xFF5E {
		r -= 0xFEE0
	}

	switch r {
	case '0':
		return 'o'
	case '1':
		return 'i'
	case '3':
		return 'e'
	case '4':
		return 'a'
	case '5':
		return 's'
	case '6', '9':
		return 'g'
	case '7':
		return 't'
	case '8':
		return 'b'
	case '@':
		return 'a'
	case '$':
		return 's'
	default:
		return r
	}
}

func (m *moderationService) findSensitiveWord(text string) string {
	m.mu.RLock()
	words := m.words
	m.mu.RUnlock()

	normalized := normalizeModerationText(text)
	if normalized == "" {
		return ""
	}

	for _, word := range words {
		count := strings.Count(normalized, word)
		if count == 0 {
			continue
		}
		if count > 1 {
			return word
		}
		if isStandaloneWord(normalized, word) {
			return word
		}
	}
	return ""
}

func isStandaloneWord(text, word string) bool {
	tr := []rune(text)
	wr := []rune(word)
	if len(wr) == 0 {
		return false
	}
	for i := 0; i <= len(tr)-len(wr); i++ {
		match := true
		for j := 0; j < len(wr); j++ {
			if tr[i+j] != wr[j] {
				match = false
				break
			}
		}
		if !match {
			continue
		}
		beforeOK := i == 0 || !isWordRune(tr[i-1])
		afterOK := i+len(wr) >= len(tr) || !isWordRune(tr[i+len(wr)])
		if beforeOK && afterOK {
			return true
		}
	}
	return false
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func (m *moderationService) setCustomWords(words []string) {
	normalizedSet := make(map[string]struct{}, len(words))
	customWords := make([]string, 0, len(words))
	for _, word := range words {
		normalized := normalizeModerationText(word)
		if normalized == "" {
			continue
		}
		if _, exists := normalizedSet[normalized]; exists {
			continue
		}
		normalizedSet[normalized] = struct{}{}
		customWords = append(customWords, normalized)
	}

	mergedSet := make(map[string]struct{}, len(m.builtinWords)+len(customWords))
	merged := make([]string, 0, len(m.builtinWords)+len(customWords))
	for _, word := range m.builtinWords {
		if _, exists := mergedSet[word]; exists {
			continue
		}
		mergedSet[word] = struct{}{}
		merged = append(merged, word)
	}
	for _, word := range customWords {
		if _, exists := mergedSet[word]; exists {
			continue
		}
		mergedSet[word] = struct{}{}
		merged = append(merged, word)
	}

	sort.Slice(merged, func(i, j int) bool {
		if len(merged[i]) == len(merged[j]) {
			return merged[i] < merged[j]
		}
		return len(merged[i]) > len(merged[j])
	})

	m.mu.Lock()
	m.customWords = customWords
	m.words = merged
	m.mu.Unlock()
}

func LoadSensitiveWords(db *gorm.DB) error {
	var items []model.SensitiveWord
	if err := db.Where("is_enabled = ?", true).Find(&items).Error; err != nil {
		return err
	}

	words := make([]string, 0, len(items))
	for _, item := range items {
		words = append(words, item.Word)
	}

	sensitiveChecker.setCustomWords(words)
	return nil
}

func LoadModerationSettings(db *gorm.DB) error {
	var item model.SystemSetting
	err := db.Where(&model.SystemSetting{Key: moderationBanThresholdKey}).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		item = model.SystemSetting{
			Key:   moderationBanThresholdKey,
			Value: strconv.Itoa(defaultModerationBanThreshold),
			Note:  "24 小时内触发多少次违禁词后自动封禁普通用户",
		}
		if createErr := db.Create(&item).Error; createErr != nil {
			return createErr
		}
		setModerationBanThreshold(defaultModerationBanThreshold)
		return nil
	}
	if err != nil {
		return err
	}

	threshold, convErr := strconv.Atoi(strings.TrimSpace(item.Value))
	if convErr != nil || threshold <= 0 {
		threshold = defaultModerationBanThreshold
	}
	setModerationBanThreshold(threshold)
	return nil
}

func setModerationBanThreshold(value int) {
	if value <= 0 {
		value = defaultModerationBanThreshold
	}
	moderationThresholdState.mu.Lock()
	moderationThresholdState.threshold = value
	moderationThresholdState.mu.Unlock()
}

func currentModerationBanThreshold() int {
	moderationThresholdState.mu.RLock()
	defer moderationThresholdState.mu.RUnlock()
	return moderationThresholdState.threshold
}

func validateModerationField(field, value string) error {
	if word := sensitiveChecker.findSensitiveWord(value); word != "" {
		return &ModerationError{Field: field, Word: word}
	}
	return nil
}

func validateModerationList(field string, values []string) error {
	for _, value := range values {
		if err := validateModerationField(field, value); err != nil {
			return err
		}
	}
	return nil
}

func moderationWord(err error) string {
	var target *ModerationError
	if errors.As(err, &target) {
		return target.Word
	}
	return ""
}

func createModerationHit(db *gorm.DB, userID uint, scene, field, value string, err error) {
	if db == nil || userID == 0 || !IsModerationError(err) {
		return
	}

	snippet := strings.TrimSpace(value)
	if len([]rune(snippet)) > 120 {
		snippet = string([]rune(snippet)[:120]) + "..."
	}

	hit := model.ModerationHit{
		UserID:      userID,
		Scene:       scene,
		Field:       field,
		MatchedWord: moderationWord(err),
		Snippet:     snippet,
	}
	if db.Create(&hit).Error != nil {
		return
	}

	notifyAdminsOfModerationHit(db, userID, scene, field, hit.MatchedWord)
	applyModerationRiskControl(db, userID)
}

func notifyAdminsOfModerationHit(db *gorm.DB, userID uint, scene, field, word string) {
	var user model.User
	if err := db.Select("id", "username", "role", "status").First(&user, userID).Error; err != nil {
		return
	}

	var admins []model.User
	if err := db.Select("id").Where("role = ? AND status = ?", model.RoleAdmin, model.UserActive).Find(&admins).Error; err != nil {
		return
	}

	title := "违禁词拦截提醒"
	content := fmt.Sprintf("用户 %s 在 %s 的 %s 命中违禁词：%s", user.Username, moderationSceneLabel(scene), field, word)
	for _, admin := range admins {
		_ = createNotification(db, NotificationCreateInput{
			UserID:  admin.ID,
			Title:   title,
			Content: content,
			Type:    model.NotificationTypeModeration,
			Payload: map[string]any{
				"userId": user.ID,
				"scene":  scene,
				"field":  field,
				"word":   word,
			},
		})
	}
}

func applyModerationRiskControl(db *gorm.DB, userID uint) {
	var user model.User
	if err := db.Select("id", "username", "role", "status").First(&user, userID).Error; err != nil {
		return
	}
	if user.Role == model.RoleAdmin || user.Status == model.UserBanned {
		return
	}

	var count int64
	if err := db.Model(&model.ModerationHit{}).
		Where("user_id = ? AND created_at >= ?", userID, time.Now().Add(-24*time.Hour)).
		Count(&count).Error; err != nil {
		return
	}

	threshold := currentModerationBanThreshold()
	if count < int64(threshold) {
		return
	}

	banReason := fmt.Sprintf("%s：24 小时内触发 %d 次违禁词拦截", autoBanReasonPrefix, count)
	if err := db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]any{
		"status":     model.UserBanned,
		"ban_reason": banReason,
	}).Error; err != nil {
		return
	}

	var admins []model.User
	if err := db.Select("id").Where("role = ? AND status = ?", model.RoleAdmin, model.UserActive).Find(&admins).Error; err != nil {
		return
	}

	title := "用户已自动风控封禁"
	content := fmt.Sprintf("用户 %s 在 24 小时内累计触发 %d 次违禁词拦截，达到阈值 %d，系统已自动封禁。", user.Username, count, threshold)
	for _, admin := range admins {
		_ = createNotification(db, NotificationCreateInput{
			UserID:  admin.ID,
			Title:   title,
			Content: content,
			Type:    model.NotificationTypeModeration,
			Payload: map[string]any{
				"userId":    user.ID,
				"hitCount":  count,
				"threshold": threshold,
			},
		})
	}
}

func moderationSceneLabel(scene string) string {
	switch scene {
	case "register":
		return "注册"
	case "profile":
		return "个人资料"
	case "article_create":
		return "文章创建"
	case "article_update":
		return "文章编辑"
	case "article_submit":
		return "文章提交"
	case "article_review":
		return "文章审核"
	case "project_create":
		return "项目创建"
	case "project_update":
		return "项目编辑"
	case "project_submit":
		return "项目提交"
	case "project_review":
		return "项目审核"
	case "comment_create":
		return "评论发布"
	default:
		return scene
	}
}
