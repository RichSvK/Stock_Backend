package entity

type Link struct {
	Name        string `json:"web_name" gorm:"type:VARCHAR(50);column:web_name;not null"`
	URL         string `json:"web_url_link" gorm:"type:TEXT;;column:url_link;not null"`
	Image       string `json:"web_image" gorm:"type:TEXT;column:web_image;not null"`
	Description string `json:"web_description" gorm:"type:TEXT;column:web_description;not null"`
	CategoryID  int    `json:"category_id" gorm:"type:INT;not null"`
}

// Make table name to "link"
func (link *Link) TableName() string {
	return "link"
}
