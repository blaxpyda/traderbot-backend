package auth

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                 string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AlpacaUserID       string    `gorm:"uniqueIndex;not null" json:"alpaca_user_id"` 
	Email              string    `json:"email,omitempty"`
	AccountType        string    `json:"account_type"`
	AlpacaAccessToken  string    `gorm:"column:alpaca_access_token" json:"-"`
	AlpacaRefreshToken string    `gorm:"column:alpaca_refresh_token" json:"-"`
	TokenExpiry        time.Time `json:"-"`
	CreatedAt          time.Time `json:"-"`
	UpdatedAt          time.Time `json:"-"`
}

// Table name
func (User) TableName() string {
	return "users"
}

// Automigrate on startup
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}