package entity

type Category struct {
	ID   int    `gorm:"column:id;type:INT;primaryKey"`
	Name string `gorm:"column:name;type:VARCHAR(100);not null"`

	// Relationship
	Web []Link `gorm:"foreignKey:category_id;references:id"`
}

// Make table name to "category"
func (category *Category) TableName() string {
	return "category"
}
