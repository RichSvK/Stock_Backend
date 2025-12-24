package response

type SearchStockResponse struct {
	Message string   `json:"message"`
	Data    []string `json:"data"`
}
