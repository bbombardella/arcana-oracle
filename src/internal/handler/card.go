package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/bbombardella/arcana-oracle/internal/cache"
	"github.com/bbombardella/arcana-oracle/internal/cards"
	"github.com/bbombardella/arcana-oracle/internal/prompts"
	"github.com/bbombardella/arcana-oracle/internal/scaleway"
	"github.com/bbombardella/arcana-oracle/internal/types"
)

type CardHandler struct {
	scw   *scaleway.Client
	cache *cache.DynamoDBCache
}

func NewCardHandler(scw *scaleway.Client, c *cache.DynamoDBCache) *CardHandler {
	return &CardHandler{scw: scw, cache: c}
}

func (h *CardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.CardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	if !cards.Valid(req.Card.Id) {
		http.Error(w, "unknown card id", http.StatusBadRequest)
		return
	}
	if req.Lang == "" {
		req.Lang = prompts.DefaultLang
	}
	if !prompts.ValidLang(req.Lang) {
		http.Error(w, "unsupported language", http.StatusBadRequest)
		return
	}

	// DynamoDB cache lookup
	if cached, hit, err := h.cache.Get(r.Context(), req.Card.Id, req.Card.Reversed, req.Lang); err == nil && hit {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Cache", "HIT")
		_, _ = io.WriteString(w, cached)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Cache", "MISS")

	prompt := prompts.BuildCardPrompt(req)
	full, err := h.scw.Stream(r.Context(), prompts.SystemPrompt(req.Lang), prompt, w)
	if err != nil || full == "" {
		return
	}

	// Store in cache asynchronously — response is already streaming to the client
	go func() {
		_ = h.cache.Set(context.Background(), req.Card.Id, req.Card.Reversed, req.Lang, full)
	}()
}
