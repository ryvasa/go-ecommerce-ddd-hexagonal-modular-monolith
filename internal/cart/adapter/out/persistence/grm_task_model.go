package persistence

type CartModel struct {
	ID     string `gorm:"primaryKey;size:36"`
	UserID string `gorm:"size:36"`
	Title  string `gorm:"size:255"`
	Done   bool
}

func (CartModel) TableName() string {
	return "carts"
}
