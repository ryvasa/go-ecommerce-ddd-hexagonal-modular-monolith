package persistence

type UserModel struct {
	ID    string `gorm:"primaryKey;size:36"`
	Email string `gorm:"uniqueIndex;size:255"`
}

func (UserModel) TableName() string {
	return "users"
}
