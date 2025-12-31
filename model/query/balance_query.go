package query_model

type BalanceChange struct {
	StockCode         string  `json:"stock_code"`
	CurrentOwnership  float64 `json:"current_ownership"`
	PreviousOwnership float64 `json:"previous_ownership"`
	ChangePercentage  float64 `json:"change_percentage"`
}
