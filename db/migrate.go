package db

import . "github.com/kodability/apiserver/models"

// AutoMigrate migrates database
func AutoMigrate() {
	Conn.AutoMigrate(
		&Question{},
		&QuestionDescription{},
		&QuestionCode{},
		&QuestionTag{},
		&Tryout{},
		&TryoutResult{},
	)
}
