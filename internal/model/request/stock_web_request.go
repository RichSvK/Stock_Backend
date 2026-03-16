package request

type CreateLinkRequest struct {
	URL         string `json:"url" validate:"required,url"`
	Name        string `json:"name" validate:"required"`
	Image       string `json:"image" validate:"required"`
	Description string `json:"description" validate:"required"`
	Category    int    `json:"category" validate:"required"`
}

type UpdateLinkRequest struct {
	URL         string `json:"url" validate:"omitempty,url"`
	Name        string `json:"name" validate:"required"`
	Image       string `json:"image" validate:"omitempty"`
	Description string `json:"description" validate:"omitempty"`
	Category    int    `json:"category" validate:"omitempty,numeric,gt=0,max=4"`
}