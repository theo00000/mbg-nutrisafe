package models

import "time"

type Teacher struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SchoolID  uint      `gorm:"not null;index" json:"school_id"`
	Name      string    `gorm:"type:varchar(150);not null" json:"name"`
	NIP       string    `gorm:"type:varchar(30);not null" json:"nip"`
	Position  string    `gorm:"type:varchar(100);not null" json:"position"`
	Gender    string    `gorm:"type:varchar(20);not null" json:"gender"`
	Address   string    `gorm:"type:text" json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
