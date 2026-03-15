package query_filter

type DateRangeQuery struct {
	StartTime     string
	StartTimeLast string
	EndTime       string
	EndTimeLast   string
}

type ScriptlessChangeQuery struct {
	StartTime string `form:"start_time" validate:"required"`
	EndTime   string `form:"end_time"  validate:"required"`
}