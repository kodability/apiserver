package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"tryout-runner/db"
	"tryout-runner/models"
	_ "tryout-runner/routers"

	"github.com/astaxie/beego"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// initialize log settings
func initLogs() {
	filename := beego.AppConfig.String("log.filename")
	if filename != "" {
		dir := filepath.Dir(filename)
		if err := os.MkdirAll(dir, 0755); err == nil {
			beego.SetLogger("file", fmt.Sprintf(`{"filename":"%s"}`, filename))
		}
	}
}

// migrate database
func migrate() {
	conn := db.Conn
	conn.AutoMigrate(
		&models.Question{},
		&models.QuestionDesc{},
		&models.QuestionCode{},
		&models.QuestionTag{},
	)
}

func main() {
	initLogs()

	conn, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to connect to DB", err)
	}
	db.Conn.LogMode(true)

	migrate()

	defer conn.Close()

	beego.Run()
}
