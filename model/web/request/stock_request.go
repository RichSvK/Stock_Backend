package request

type SearchStockQuery struct {
	Code string `form:"code" validate:"required,alpha,min=1,max=4"`
}
