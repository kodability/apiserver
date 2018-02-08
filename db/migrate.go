package db

import "tryout-runner/models"

// AutoMigrate migrates database
func AutoMigrate() {
	Conn.AutoMigrate(
		&models.Question{},
		&models.QuestionDescription{},
		&models.QuestionCode{},
		&models.QuestionTag{},
	)
}
