package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"tryout-runner/database"
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

func main() {
	initLogs()

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to DB", err)
	}
	db.AutoMigrate()

	defer db.Close()

	beego.Run()
}
