package models

import "time"

type PasswordReset struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"type:varchar(100);not null;index" json:"email"`
	CodeHash  string    `gorm:"type:varchar(255);not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	Attempts  int       `gorm:"default:0" json:"attempts"`
	CreatedAt time.Time `json:"created_at"`
}
