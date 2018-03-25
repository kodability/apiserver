package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"github.com/kodability/apiserver/db"
	"github.com/kodability/apiserver/models"
	"github.com/kodability/apiserver/services/run"
)

// RunController defines operations about tryout run.
type RunController struct {
	beego.Controller
}

type ITryoutRunner interface {
	Run(lang, code, testCode string) (*run.JUnitReport, error)
}

var tryoutRunner ITryoutRunner

// get ITryoutRunner instance based on 'tryout.runner' configuration.
func getTryoutRunner() ITryoutRunner {
	if tryoutRunner == nil {
		runner := beego.AppConfig.String("tryout.runner")
		if runner == "docker" {
			tryoutRunner = &run.TryoutDockerRunner{
				TempDir:    "/tmp",
				TempPrefix: "kodability-",
			}
		}
	}
	return tryoutRunner
}

// set ITryoutRunner instance. Used for testing.
func SetTryoutRunner(r ITryoutRunner) {
	tryoutRunner = r
}

// RunBody is a body for RunTryout() method.
type RunBody struct {
	QuestionID uint
	Lang       string
	Code       string
}

// RunTryout runs a tryout
func (c *RunController) RunTryout() {
	var body RunBody
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		badRequest(&c.Controller, err.Error())
		return
	}

	var err error
	conn := db.Conn

	// Find QuestionCode
	var questionCode models.QuestionCode
	if err := conn.Where("question_id = ? AND lang = ?", body.QuestionID, body.Lang).First(&questionCode).Error; err != nil {
		badRequest(&c.Controller, fmt.Sprintf("Question not found. id=%v", body.QuestionID))
		return
	}

	// Save Tryout
	tryout := models.Tryout{
		QuestionID: body.QuestionID,
		Lang:       body.Lang,
		Code:       body.Code,
	}
	err = conn.Create(&tryout).Error
	if err != nil {
		internalServerError(&c.Controller, fmt.Sprintf("Failed to add Tryout: %v", err.Error()))
		return
	}

	// Run Tryout
	if tryoutRunner := getTryoutRunner(); tryoutRunner == nil {
		msg := "TryoutRunner is nil"
		createTryoutResultError(conn, tryout.ID, msg)
		internalServerError(&c.Controller, msg)
		return
	}
	report, err := tryoutRunner.Run(tryout.Lang, body.Code, questionCode.TestCode)
	if err != nil {
		msg := fmt.Sprintf("Failed to run Tryout: %v", err.Error())
		createTryoutResultError(conn, tryout.ID, msg)
		internalServerError(&c.Controller, msg)
		return
	}

	// Convert JUnitReport to TryoutResult
	result := report.ToTryoutResult()
	result.TryoutID = tryout.ID
	if err = conn.Create(&result).Error; err != nil {
		internalServerError(&c.Controller, err.Error())
		return
	}

	jsonCreated(&c.Controller, result)
}

func createTryoutResultError(conn *gorm.DB, tryoutID uint, errorMsg string) error {
	result := models.TryoutResult{
		TryoutID: tryoutID,
		ErrorMsg: errorMsg,
	}
	return conn.Create(&result).Error
}
