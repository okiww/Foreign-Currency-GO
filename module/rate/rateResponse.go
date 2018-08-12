package rate

type rateResponse struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
	Rate string `json:"date" binding:"required"`
	Avg  string `json:"7_day_avg" binding:"required"`
}

type rateMost7DataPointResponse struct {
	From      string `json:"from" binding:"required"`
	To        string `json:"to" binding:"required"`
	Average   string `json:"average" binding:"required"`
	Varriance string `json:"varriance" binding:"required"`
	Data      []dataByDateAndRate
}

type dataByDateAndRate struct {
	Date string `json:"date" binding:"required"`
	Rate string `json:"rate" binding:"required"`
}
