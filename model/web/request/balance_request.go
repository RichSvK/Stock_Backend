package request

type GetBalanceChangeQuery struct {
	Type   string `form:"type" validate:"required"`
	Change string `form:"change" validate:"required,oneof=Increase Decrease"`
	Page   int    `form:"page,default=1" validate:"omitempty,min=1"`
}
