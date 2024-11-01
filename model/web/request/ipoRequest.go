package request

type Filter struct {
	FilterName  string `json:"filter_name"`
	FilterValue string `json:"filter_value"`
	Symbol      string `json:"symbol"`
	FilterType  string `json:"filter_type"`
}
