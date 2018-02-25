package tests

import (
	"log"
	"path/filepath"

	"github.com/kodability/tryout-runner/db"

	"github.com/astaxie/beego"
)

func init() {

	if apppath, err := filepath.Abs(".."); err != nil {
		log.Fatal(err.Error())
	} else {
		beego.TestBeegoInit(apppath)
	}

	if _, err := db.Connect(); err != nil {
		log.Fatal("Failed to connect to DB", err)
	}

	db.AutoMigrate()
}
