package rate

type rateRequest struct {
	CurrencyID uint    `json:"currency_id" binding:"required"`
	Rate       float64 `json:"rate" binding:"required"`
	Date       string  `json:"date" binding:"required"`
}

type rateMost7DataPointRequest struct {
	CurrencyID uint `json:"currency_id" binding:"required"`
}
