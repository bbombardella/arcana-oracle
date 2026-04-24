package types

type CardRequest struct {
	Card     CardInfo      `json:"card"`
	Position *PositionInfo `json:"position,omitempty"`
	Lang     string        `json:"lang"` // ISO 639-1, e.g. "fr"
}

type CardInfo struct {
	Id       string `json:"id"`       // e.g. "major-13", "cups-07"
	Reversed bool   `json:"reversed"`
}

type PositionInfo struct {
	Index      int    `json:"index"`
	Label      string `json:"label"`
	SpreadSize int    `json:"spreadSize"`
}
