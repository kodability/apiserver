package tests

import (
	"log"
	"path/filepath"
	"tryout-runner/db"

	"github.com/astaxie/beego"
)

func init() {
	apppath, _ := filepath.Abs("..")
	beego.TestBeegoInit(apppath)

	if _, err := db.Connect(); err != nil {
		log.Fatal("Failed to connect to DB", err)
	}

	db.AutoMigrate()
}
