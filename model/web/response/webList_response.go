package response

import "backend/model/entity"

type GetLinkResponse struct {
	Message string         `json:"message"`
	Data    []LinkResponse `json:"data"`
}

type LinkResponse struct {
	URL         string `json:"web_url_link" gorm:"type:TEXT;column:url_link"`
	Name        string `json:"web_name" gorm:"type:VARCHAR(50);column:web_name;not null"`
	Image       string `json:"web_image" gorm:"type:TEXT;column:web_image;not null"`
	Description string `json:"web_description" gorm:"type:TEXT;column:web_description;not null"`
}

func MapLinksToResponse(links []entity.Link) []LinkResponse {
	result := make([]LinkResponse, 0, len(links))

	for _, link := range links {
		result = append(result, LinkResponse{
			URL:         link.URL,
			Name:        link.Name,
			Image:       link.Image,
			Description: link.Description,
		})
	}

	return result
}
