package database

import (
	"fmt"
	"log"

	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
)

// DB is a database handle representing a pool of zero or more underlying connections.
// It's safe to concurrent use by multiple goroutines.
var DB *gorm.DB

// Connect connect to database.
// After successful connection, it's free to use 'Conn' variable
func Connect() (*gorm.DB, error) {
	dialect := beego.AppConfig.String("db.dialect")
	if dialect == "sqlite3" {
		filename := beego.AppConfig.String("db.filename")
		db, err := newSqlite3(filename)
		DB = db
		return db, err
	}

	return nil, fmt.Errorf("Unknown dialect: %s", dialect)
}

func newSqlite3(filename string) (*gorm.DB, error) {
	log.Printf("Connecting to Sqlite3 DB : %s\n", filename)
	return gorm.Open("sqlite3", filename)
}

func newMysql(dsn string) (*gorm.DB, error) {
	log.Printf("Connecting to Mysql DB : %s\n", dsn)
	return gorm.Open("mysql", dsn)
}
