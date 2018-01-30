package main

import (
	"fmt"
	"log"
	_ "tryout-runner/routers"

	"github.com/jinzhu/gorm"

	"github.com/astaxie/beego"
)

// NewDB creates an instance of gorm.DB
func NewDB() (*gorm.DB, error) {
	dialect := beego.AppConfig.String("db.dialect")
	if dialect == "sqlite3" {
		filename := beego.AppConfig.String("db.filename")
		return NewSqlite3(filename)
	}

	return nil, fmt.Errorf("Unknown dialect: %s", dialect)
}

func main() {
	var db *gorm.DB
	db, err := NewDB()
	if err != nil {
		log.Fatal("Failed to connect to DB", err)
	}

	defer db.Close()

	beego.Run()
}
