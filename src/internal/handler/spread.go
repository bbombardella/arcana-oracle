package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bbombardella/arcana-oracle/internal/cards"
	"github.com/bbombardella/arcana-oracle/internal/prompts"
	"github.com/bbombardella/arcana-oracle/internal/scaleway"
	"github.com/bbombardella/arcana-oracle/internal/types"
)

type SpreadHandler struct {
	scw *scaleway.Client
}

func NewSpreadHandler(scw *scaleway.Client) *SpreadHandler {
	return &SpreadHandler{scw: scw}
}

func (h *SpreadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.SpreadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	for _, c := range req.Cards {
		if !cards.Valid(c.Id) {
			http.Error(w, "unknown card id: "+c.Id, http.StatusBadRequest)
			return
		}
	}
	if req.Lang == "" {
		req.Lang = prompts.DefaultLang
	}
	if !prompts.ValidLang(req.Lang) {
		http.Error(w, "unsupported language", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	prompt := prompts.BuildSpreadPrompt(req)
	_, _ = h.scw.Stream(r.Context(), prompts.SystemPrompt(req.Lang), prompt, w)
}
