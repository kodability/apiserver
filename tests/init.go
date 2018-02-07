package test

import (
	"log"
	"path/filepath"
	"runtime"
	"tryout-runner/db"
	_ "tryout-runner/routers"

	"github.com/astaxie/beego"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)

	conn, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to connect to DB", err)
	}
	conn.LogMode(true)
	db.AutoMigrate()
}
