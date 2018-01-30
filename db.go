package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// NewSqlite3 creates a connection to sqlite3
func NewSqlite3(filename string) (*gorm.DB, error) {
	log.Printf("DB Connection : %s\n", filename)
	return gorm.Open("sqlite3", filename)
}

// NewMysql creates a connection to mysql
func NewMysql(user, password, dbname string) (*gorm.DB, error) {
	return gorm.Open("mysql",
		fmt.Sprintf("%s:%s@/%s?charset=utf8", user, password, dbname))
}
