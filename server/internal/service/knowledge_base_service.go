package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"blog/server/internal/model"
	"gorm.io/gorm"
)

var (
	ErrKBNotFound    = errors.New("知识库不存在")
	ErrKBAccessDenied = errors.New("无权访问该知识库")
	ErrDocNotFound    = errors.New("文档不存在")
	ErrKBNameEmpty    = errors.New("知识库名称不能为空")
	ErrDocContentEmpty = errors.New("文档内容不能为空")
	ErrKBNoAPIKey     = errors.New("请先配置 API Key")
)

const (
	chunkSize    = 500
	chunkOverlap = 100
	maxQueryDocs = 5
	embeddingDim = 1536
)

type KnowledgeBaseService struct {
	db       *gorm.DB
	qdrant   *QdrantClient
	embedding EmbeddingProvider
}

func NewKnowledgeBaseService(db *gorm.DB, qdrantAddr, qdrantKey string, embedding EmbeddingProvider) *KnowledgeBaseService {
	qdrant := NewQdrantClient(qdrantAddr, qdrantKey)
	return &KnowledgeBaseService{
		db:        db,
		qdrant:    qdrant,
		embedding: embedding,
	}
}

func (s *KnowledgeBaseService) collectionName(kbID uint) string {
	return fmt.Sprintf("kb_%d", kbID)
}

func (s *KnowledgeBaseService) Create(userID uint, name, description string) (*model.KnowledgeBase, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrKBNameEmpty
	}

	kb := model.KnowledgeBase{
		UserID:      userID,
		Name:        name,
		Description: strings.TrimSpace(description),
	}
	if err := s.db.Create(&kb).Error; err != nil {
		return nil, err
	}

	if err := s.qdrant.HealthCheck(); err != nil {
		s.db.Delete(&kb)
		return nil, fmt.Errorf("知识库服务(Qdrant)连接失败，请确认服务已启动: %w", err)
	}
	if err := s.qdrant.EnsureCollection(s.collectionName(kb.ID), embeddingDim); err != nil {
		s.db.Delete(&kb)
		return nil, fmt.Errorf("创建向量集合失败: %w", err)
	}

	return &kb, nil
}

func (s *KnowledgeBaseService) List(userID uint) ([]model.KnowledgeBase, error) {
	var kbs []model.KnowledgeBase
	if err := s.db.Where("user_id = ?", userID).Order("updated_at desc").Find(&kbs).Error; err != nil {
		return nil, err
	}
	return kbs, nil
}

func (s *KnowledgeBaseService) GetByID(kbID, userID uint) (*model.KnowledgeBase, error) {
	var kb model.KnowledgeBase
	if err := s.db.First(&kb, kbID).Error; err != nil {
		return nil, ErrKBNotFound
	}
	if kb.UserID != userID {
		return nil, ErrKBAccessDenied
	}
	return &kb, nil
}

func (s *KnowledgeBaseService) Delete(kbID, userID uint) error {
	kb, err := s.GetByID(kbID, userID)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("knowledge_base_id = ?", kbID).Delete(&model.KnowledgeDocument{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(kb).Error; err != nil {
			return err
		}
		return s.qdrant.DeleteCollection(s.collectionName(kb.ID))
	})
}

type AddDocumentOpts struct {
	Title      string
	Content    string
	IsPublic   bool
	IsMarkdown bool
	CategoryID *uint
	TagNames   []string
}

func (s *KnowledgeBaseService) AddDocument(kbID, userID uint, opts AddDocumentOpts) (*model.KnowledgeDocument, error) {
	if _, err := s.GetByID(kbID, userID); err != nil {
		return nil, err
	}

	opts.Content = strings.TrimSpace(opts.Content)
	if opts.Content == "" {
		return nil, ErrDocContentEmpty
	}
	opts.Title = strings.TrimSpace(opts.Title)

	now := time.Now()
	doc := model.KnowledgeDocument{
		KnowledgeBaseID: kbID,
		UserID:          userID,
		Title:           opts.Title,
		Content:         opts.Content,
		SourceType:      "manual",
		IsPublic:        opts.IsPublic,
		IsMarkdown:      opts.IsMarkdown,
		CategoryID:      opts.CategoryID,
	}

	if opts.IsPublic {
		doc.PublishedAt = &now
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&doc).Error; err != nil {
			return err
		}

		if len(opts.TagNames) > 0 {
			tags, err := s.ensureTags(tx, opts.TagNames)
			if err != nil {
				return err
			}
			if err := tx.Model(&doc).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}

		return tx.Model(&model.KnowledgeBase{}).Where("id = ?", kbID).
			UpdateColumn("doc_count", gorm.Expr("doc_count + 1")).Error
	})
	if err != nil {
		return nil, err
	}

	go s.tryVectorize(kbID, doc, opts.Content, opts.Title)

	doc.Tags = nil
	return &doc, nil
}

func (s *KnowledgeBaseService) UpdateDocument(docID, kbID, userID uint, opts AddDocumentOpts) (*model.KnowledgeDocument, error) {
	var doc model.KnowledgeDocument
	if err := s.db.Where("id = ? AND knowledge_base_id = ? AND user_id = ?", docID, kbID, userID).First(&doc).Error; err != nil {
		return nil, ErrDocNotFound
	}

	now := time.Now()
	updates := map[string]interface{}{
		"title":       opts.Title,
		"content":     opts.Content,
		"is_public":   opts.IsPublic,
		"is_markdown": opts.IsMarkdown,
		"category_id": opts.CategoryID,
	}

	if opts.IsPublic && !doc.IsPublic {
		updates["published_at"] = &now
	} else if !opts.IsPublic && doc.IsPublic {
		updates["published_at"] = nil
	}

	if err := s.db.Model(&doc).Updates(updates).Error; err != nil {
		return nil, err
	}

	if len(opts.TagNames) > 0 {
		tags, err := s.ensureTags(s.db, opts.TagNames)
		if err != nil {
			return nil, err
		}
		if err := s.db.Model(&doc).Association("Tags").Replace(tags); err != nil {
			return nil, err
		}
	}

	go s.tryVectorize(kbID, doc, opts.Content, opts.Title)

	var updated model.KnowledgeDocument
	s.db.Preload("Category").Preload("Tags").First(&updated, doc.ID)
	return &updated, nil
}

func (s *KnowledgeBaseService) ListDocuments(kbID, userID uint) ([]model.KnowledgeDocument, error) {
	if _, err := s.GetByID(kbID, userID); err != nil {
		return nil, err
	}

	var docs []model.KnowledgeDocument
	if err := s.db.Where("knowledge_base_id = ?", kbID).
		Preload("Category").Preload("Tags").
		Order("created_at desc").Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (s *KnowledgeBaseService) GetPublicNote(noteID uint) (*model.KnowledgeDocument, error) {
	var doc model.KnowledgeDocument
	if err := s.db.Where("id = ? AND is_public = ?", noteID, true).
		Preload("User").Preload("Category").Preload("Tags").
		First(&doc).Error; err != nil {
		return nil, ErrDocNotFound
	}
	s.db.Model(&doc).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	return &doc, nil
}

func (s *KnowledgeBaseService) ensureTags(tx *gorm.DB, names []string) ([]model.Tag, error) {
	var tags []model.Tag
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		var tag model.Tag
		if err := tx.Where("name = ?", name).FirstOrCreate(&tag, model.Tag{Name: name}).Error; err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (s *KnowledgeBaseService) tryVectorize(kbID uint, doc model.KnowledgeDocument, content, title string) {
	if s.embedding == nil {
		return
	}
	if err := s.qdrant.HealthCheck(); err != nil {
		return
	}

	chunks := s.chunkText(content)
	vectors, err := s.embedding.Embed(chunks)
	if err != nil {
		return
	}

	points := make([]QdrantPoint, len(chunks))
	for i := range chunks {
		points[i] = QdrantPoint{
			ID:     uint64(doc.ID)*10000 + uint64(i),
			Vector: vectors[i],
			Payload: map[string]interface{}{
				"kb_id":  kbID,
				"doc_id": doc.ID,
				"title":  title,
				"chunk":  i,
				"text":   chunks[i],
			},
		}
	}

	if err := s.qdrant.UpsertPoints(s.collectionName(kbID), points); err != nil {
		return
	}

	s.db.Model(&doc).UpdateColumn("chunk_count", len(chunks))
}

func (s *KnowledgeBaseService) DeleteDocument(kbID, docID, userID uint) error {
	if _, err := s.GetByID(kbID, userID); err != nil {
		return err
	}

	var doc model.KnowledgeDocument
	if err := s.db.Where("id = ? AND knowledge_base_id = ?", docID, kbID).First(&doc).Error; err != nil {
		return ErrDocNotFound
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&doc).Error; err != nil {
			return err
		}

		ids := make([]uint64, doc.ChunkCount)
		for i := 0; i < doc.ChunkCount; i++ {
			ids[i] = uint64(doc.ID)*10000 + uint64(i)
		}
		if err := s.qdrant.DeletePoints(s.collectionName(kbID), ids); err != nil {
			return err
		}

		return tx.Model(&model.KnowledgeBase{}).Where("id = ?", kbID).
			UpdateColumn("doc_count", gorm.Expr("doc_count - 1")).Error
	})
}

type QueryResult struct {
	Answer  string   `json:"answer"`
	Sources []Source `json:"sources"`
}

type Source struct {
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

type LLMProvider interface {
	Chat(systemPrompt, userMessage string) (string, error)
	Name() string
}

func (s *KnowledgeBaseService) Query(kbID, userID uint, question string, llm LLMProvider) (*QueryResult, error) {
	kb, err := s.GetByID(kbID, userID)
	if err != nil {
		return nil, err
	}
	question = strings.TrimSpace(question)
	if question == "" {
		return nil, fmt.Errorf("问题不能为空")
	}

	if s.embedding == nil {
		return nil, ErrKBNoAPIKey
	}

	vectors, err := s.embedding.Embed([]string{question})
	if err != nil {
		return nil, fmt.Errorf("生成查询向量失败: %w", err)
	}

	results, err := s.qdrant.Search(s.collectionName(kbID), vectors[0], maxQueryDocs)
	if err != nil {
		return nil, fmt.Errorf("向量检索失败: %w", err)
	}

	sources := make([]Source, 0, len(results))
	var contextBuilder strings.Builder
	for _, r := range results {
		title, _ := r.Payload["title"].(string)
		text, _ := r.Payload["text"].(string)

		sources = append(sources, Source{
			Title:   title,
			Content: text,
			Score:   r.Score,
		})
		contextBuilder.WriteString(fmt.Sprintf("---\n标题: %s\n内容: %s\n", title, text))
	}

	if len(sources) == 0 {
		return &QueryResult{
			Answer:  "知识库中没有找到相关内容。",
			Sources: sources,
		}, nil
	}

	answer := ""
	if llm != nil {
		systemPrompt := fmt.Sprintf(`你是一个知识库问答助手。请基于以下资料回答用户的问题。
如果资料不足以回答问题，请如实说明。

知识库名称: %s
知识库描述: %s

相关资料:
%s`, kb.Name, kb.Description, contextBuilder.String())

		answer, err = llm.Chat(systemPrompt, question)
		if err != nil {
			return nil, fmt.Errorf("生成回答失败: %w", err)
		}
	} else {
		answer = fmt.Sprintf("找到 %d 条相关结果。请配置 LLM API Key 以获取 AI 回答。", len(sources))
	}

	return &QueryResult{
		Answer:  answer,
		Sources: sources,
	}, nil
}

func (s *KnowledgeBaseService) chunkText(text string) []string {
	runes := []rune(text)
	if len(runes) <= chunkSize {
		return []string{text}
	}

	var chunks []string
	start := 0
	for start < len(runes) {
		end := start + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
		start += chunkSize - chunkOverlap
		if start >= len(runes) {
			break
		}
	}
	return chunks
}
