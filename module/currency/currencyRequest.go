package currency

type currencyRequest struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
}

type deleteCurrencyRequest struct {
	CurrencyID uint `json:"currency_id" binding:"required"`
}
