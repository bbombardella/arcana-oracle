package types

type AstroRequest struct {
	Sign SignInfo  `json:"sign"`
	Card CardInfo  `json:"card"`
	Lang string    `json:"lang"`
}

type SignInfo struct {
	Name    string `json:"name"`
	Element string `json:"element"`
}
