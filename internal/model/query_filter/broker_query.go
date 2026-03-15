package query_filter

type BrokerQuery struct {
	Code string `form:"code" validate:"alpha,len=2"`
}