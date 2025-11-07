package database

import (
	"expense-tracker-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("expenses.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto-migrate the Expense model
	err = DB.AutoMigrate(&models.Expense{})
	if err != nil {
		return err
	}

	return nil
}

