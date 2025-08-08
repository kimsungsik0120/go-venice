package dto

type BalanceResponse struct {
	Amount string `json:"amount"`
	Symbol string `json:"symbol"`
}

type TransactionResponse struct {
	Tx string `json:"tx"`
}
