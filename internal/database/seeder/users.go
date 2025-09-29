package seeder

import (
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database/model"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	users := []model.User{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
	}
	for _, u := range users {
		if err := db.Create(&u).Error; err != nil {
			return err
		}
	}
	return nil
}
