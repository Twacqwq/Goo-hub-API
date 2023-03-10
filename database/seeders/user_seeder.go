package seeders

import (
	"fmt"
	"thub/database/factories"
	"thub/pkg/console"
	"thub/pkg/logger"
	"thub/pkg/seed"

	"gorm.io/gorm"
)

func init() {
	seed.Add("SeedUsersTable", func(db *gorm.DB) {
		users := factories.MakeUsers(10)

		result := db.Table("users").Create(&users)

		if err := result.Error; err != nil {
			logger.LogIf(err)
			return
		}

		console.Success(fmt.Sprintf("Table [%v] %v rows seeded", result.Statement.Table, result.Statement.RowsAffected))
	})
}
