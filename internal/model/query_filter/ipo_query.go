package query_filter

type GetIpoQuery struct {
	Code string `form:"code" validate:"alpha,len=4"`
}