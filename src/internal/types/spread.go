package types

type SpreadRequest struct {
	Cards      []CardInfo `json:"cards"`
	SpreadSize int        `json:"spreadSize"`
	Lang       string     `json:"lang"`
}
