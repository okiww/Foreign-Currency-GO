package currency

type currencyResponse struct {
	ID    uint            `json:"id"`
	From  string          `json:"from"`
	To    string          `json:"to"`
	Rates []ratesRepsonse `json:"rates"`
}

type ratesRepsonse struct {
	ID   uint   `json:"id"`
	Date string `json:"date"`
	Rate string `json:"rate"`
}
