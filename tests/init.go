package tests

import (
	"log"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/kodability/apiserver/db"
	m "github.com/kodability/apiserver/models"
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

func deleteQuestionsAndDescAndCodes() {
	conn := db.Conn
	conn.Unscoped().Delete(m.Question{})
	conn.Unscoped().Delete(m.QuestionDescription{})
	conn.Unscoped().Delete(m.QuestionCode{})
}
