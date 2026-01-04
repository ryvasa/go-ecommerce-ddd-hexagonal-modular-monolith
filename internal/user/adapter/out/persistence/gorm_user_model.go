package persistence

import (
	"time"

	"gorm.io/gorm"
)

type (
	UserModel struct {
		gorm.Model
		ID              string     `gorm:"primaryKey;size:36"`
		Email           string     `gorm:"uniqueIndex;size:255;not null"`
		EmailVerifiedAt *time.Time `gorm:"type:datetime"`
		Password        string     `gorm:"size:255;not null"`
		Username        string     `gorm:"size:255;not null"`
		FirstName       string     `gorm:"size:255;not null"`
		LastName        string     `gorm:"size:255;not null"`
		Phone           *string    `gorm:"size:255"`
		BirthDate       *time.Time `gorm:"type:datetime"`
	}
)

func (UserModel) TableName() string {
	return "users"
}
