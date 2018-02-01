package main

import (
	"log"
	"tryout-runner/db"
	_ "tryout-runner/routers"

	"github.com/astaxie/beego"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	conn, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to connect to DB", err)
	}

	defer conn.Close()

	beego.Run()
}
