package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bbombardella/arcana-oracle/internal/cards"
	"github.com/bbombardella/arcana-oracle/internal/prompts"
	"github.com/bbombardella/arcana-oracle/internal/scaleway"
	"github.com/bbombardella/arcana-oracle/internal/types"
)

type AstroHandler struct {
	scw *scaleway.Client
}

func NewAstroHandler(scw *scaleway.Client) *AstroHandler {
	return &AstroHandler{scw: scw}
}

func (h *AstroHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.AstroRequest
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

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	prompt := prompts.BuildAstroPrompt(req)
	_, _ = h.scw.Stream(r.Context(), prompts.SystemPrompt(req.Lang), prompt, w)
}
