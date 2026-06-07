package models

import "time"

type PasswordReset struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	Email     string    `gorm:"type:varchar(100);not null;index" json:"-"`
	CodeHash  string    `gorm:"type:varchar(255);not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"-"`
	Used      bool      `gorm:"default:false" json:"-"`
	Attempts  int       `gorm:"default:0" json:"-"`
	CreatedAt time.Time `json:"-"`
}
