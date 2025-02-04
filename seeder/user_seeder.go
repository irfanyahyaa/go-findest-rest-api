package seeder

import (
	"errors"
	"fmt"
	"go-findest-rest-api/model"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	users := []model.User{
		{Name: "Daiki Tsuneta"},
		{Name: "Satoru Iguchi"},
		{Name: "Kazuki Arai"},
		{Name: "Yu Seki"},
	}

	for _, user := range users {
		var existing model.User
		if err := db.Where("name = ?", user.Name).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				db.Create(&user)
				fmt.Println("User created:", user.Name)
			}
		}
	}
}
