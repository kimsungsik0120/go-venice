package dto

type BalanceRequest struct {
	Amount string `json:"amount"`
	Symbol string `json:"symbol"`
}
