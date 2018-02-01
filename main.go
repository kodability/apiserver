package main

import (
	"log"
	"tryout-runner/db"
	"tryout-runner/models"
	_ "tryout-runner/routers"

	"github.com/astaxie/beego"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func migrate() {
	db.Conn.AutoMigrate(&models.Question{},
		&models.QuestionDesc{},
		&models.QuestionCode{},
		&models.QuestionTag{},
	)
}

func main() {
	conn, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to connect to DB", err)
	}

	migrate()

	defer conn.Close()

	beego.Run()
}
