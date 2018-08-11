package currency

type currencyRequest struct {
	Date string  `json:"date" binding:"required"`
	From string  `json:"from" binding:"required"`
	To   string  `json:"to" binding:"required"`
	Rate float32 `json:"rate" binding:"required"`
}
