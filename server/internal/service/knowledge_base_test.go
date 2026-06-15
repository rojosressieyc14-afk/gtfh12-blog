package service

import (
	"strings"
	"testing"
)

func TestChunkTextSmall(t *testing.T) {
	kb := &KnowledgeBaseService{}
	text := "short text"
	chunks := kb.chunkText(text)
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
	if chunks[0] != text {
		t.Fatalf("chunk mismatch: %q", chunks[0])
	}
}

func TestChunkTextLarge(t *testing.T) {
	kb := &KnowledgeBaseService{}

	words := make([]string, 2000)
	for i := range words {
		words[i] = "word"
	}
	text := strings.Join(words, " ")

	chunks := kb.chunkText(text)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}

	for i, chunk := range chunks {
		runes := []rune(chunk)
		if i < len(chunks)-1 && len(runes) > chunkSize+chunkOverlap {
			t.Fatalf("chunk %d too long: %d chars (max %d)", i, len(runes), chunkSize)
		}
	}
}

func TestChunkTextOverlap(t *testing.T) {
	kb := &KnowledgeBaseService{}

	runes := make([]rune, chunkSize*2+50)
	for i := range runes {
		runes[i] = rune('a' + i%26)
	}
	text := string(runes)

	chunks := kb.chunkText(text)
	if len(chunks) < 2 {
		t.Fatalf("expected at least 2 chunks, got %d", len(chunks))
	}

	if len(chunks) > 3 {
		t.Fatalf("expected 2-3 chunks for text of length %d, got %d", len(runes), len(chunks))
	}
}
