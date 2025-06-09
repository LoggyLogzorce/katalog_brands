package models

type Category struct {
	ID           uint64 `gorm:"primaryKey" json:"category_id"`
	Name         string `json:"name"`
	Photo        string `json:"photo"`
	ProductCount int    `json:"product_count"`
}

// TableName задаёт имя таблицы для Category
func (Category) TableName() string {
	return "categories"
}
