package query_filter

type GetLinkQuery struct {
	Name       string `form:"name" validate:"omitempty"`
	CategoryID string `form:"category_id" validate:"omitempty,numeric"`
}