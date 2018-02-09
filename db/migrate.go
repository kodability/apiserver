package db

import "github.com/kodability/tryout-runner/models"

// AutoMigrate migrates database
func AutoMigrate() {
	Conn.AutoMigrate(
		&models.Question{},
		&models.QuestionDescription{},
		&models.QuestionCode{},
		&models.QuestionTag{},
	)
}
