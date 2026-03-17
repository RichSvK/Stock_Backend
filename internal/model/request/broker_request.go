package request

type CreateBrokerRequest struct {
	Code string `json:"code" validate:"required,alpha,len=2"`
	Name string `json:"name" validate:"required"`
}

type UpdateBrokerRequest struct {
	Code string `json:"code" validate:"required,alpha,len=2"`
	Name string `json:"name" validate:"required"`
}