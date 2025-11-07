package models

import (
	"time"
	// "gorm.io/gorm"
)

type Expense struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Title    string    `gorm:"not null" json:"title"`
	Category string    `gorm:"not null" json:"category"`
	Type     string    `gorm:"not null;default:'Debet'" json:"type"`
	Amount   float64   `gorm:"not null" json:"amount"`
	Date     time.Time `gorm:"not null" json:"date"`
}
