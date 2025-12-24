package entity

type Link struct {
	ID          int    `gorm:"column:id;type:INT;primaryKey"`
	URL         string `gorm:"type:TEXT;column:url_link"`
	Name        string `gorm:"type:VARCHAR(50);column:web_name;not null"`
	Image       string `gorm:"type:TEXT;column:web_image;not null"`
	Description string `gorm:"type:TEXT;column:web_description;not null"`
	CategoryID  int    `gorm:"type:INT;not null"`
}

// Make table name to "link"
func (link *Link) TableName() string {
	return "link"
}
