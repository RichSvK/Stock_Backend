package request

type Filter struct {
	FilterName  string `json:"filter_name" validate:"required"`
	FilterValue string `json:"filter_value" validate:"required"`
	Symbol      string `json:"symbol" validate:"required"`
	FilterType  string `json:"filter_type" validate:"required"`
}