package rate

type rateResponse struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
	Rate string `json:"date" binding:"required"`
	Avg  string `json:"7_day_avg" binding:"required"`
}
