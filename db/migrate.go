package db

import . "github.com/kodability/tryout-runner/models"

// AutoMigrate migrates database
func AutoMigrate() {
	Conn.AutoMigrate(
		&Question{},
		&QuestionDescription{},
		&QuestionCode{},
		&QuestionTag{},
		&Tryout{},
	)
}
