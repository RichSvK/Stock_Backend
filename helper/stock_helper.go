package helper

import "backend/model/web/response"

func ToSearchStockResponses(stock []string) response.SearchStockResponse {
	return response.SearchStockResponse{
		Data: stock,
	}
}
