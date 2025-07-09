package response

type ScriptlessResponse struct {
	Code               string  `json:"code"`
	FirstShare         uint64  `json:"first_share"`
	SecondShare        uint64  `json:"second_share"`
	FirstListedShares  uint64  `json:"first_listed_share"`
	SecondListedShares uint64  `json:"second_listed_share"`
	Change             int64   `json:"change"`
	ChangePercentage   float64 `json:"change_percentage"`
}