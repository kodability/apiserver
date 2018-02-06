package database

import "tryout-runner/models"

// AutoMigrate migrates database
func AutoMigrate() {
	DB.AutoMigrate(
		&models.Question{},
		&models.QuestionDesc{},
		&models.QuestionCode{},
		&models.QuestionTag{},
	)
}
