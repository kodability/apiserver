package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/kodability/tryout-runner/db"
	"github.com/kodability/tryout-runner/models"

	"github.com/astaxie/beego"
)

// QuestionLocaleDesc defines question title and descriptions for given locale
type QuestionLocaleDesc struct {
	LocaleID string
	Title    string
	Desc     string
}

// QuestionLangCode defines code for given language.
type QuestionLangCode struct {
	Lang     string
	InitCode string
	TestCode string
}

// QuestionBody is a body struct for Post
type QuestionBody struct {
	Desc          []QuestionLocaleDesc
	Codes         []QuestionLangCode
	Level         int
	EstimatedTime int
	Tags          string
	Demo          bool
}

// =============================================================================

// QuestionController defines questions http requests
type QuestionController struct {
	beego.Controller
}

// Post a new question
func (c *QuestionController) Post() {
	var body QuestionBody
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	// Validate
	if len(body.Desc) == 0 {
		badRequest(&c.Controller, "desc cannot be empty")
		return
	}

	conn := db.Conn
	tx := conn.Begin()

	// Insert question
	question := models.Question{
		Level:         body.Level,
		EstimatedTime: body.EstimatedTime,
		Tags:          body.Tags,
		Demo:          body.Demo,
	}
	if err := tx.Create(&question).Error; err != nil {
		tx.Rollback()
		internalServerError(&c.Controller, err.Error())
		return
	}

	// Insert QuestionDesc
	for _, desc := range body.Desc {
		questionDesc := models.QuestionDescription{
			QuestionID:  question.ID,
			LocaleID:    desc.LocaleID,
			Title:       desc.Title,
			Description: desc.Desc,
		}
		if err := tx.Create(&questionDesc).Error; err != nil {
			tx.Rollback()
			internalServerError(&c.Controller, err.Error())
			return
		}
	}

	// Insert QuestionCode
	for _, code := range body.Codes {
		questionCode := models.QuestionCode{
			QuestionID: question.ID,
			Lang:       code.Lang,
			InitCode:   code.InitCode,
			TestCode:   code.TestCode,
		}
		if err := tx.Create(&questionCode).Error; err != nil {
			tx.Rollback()
			internalServerError(&c.Controller, err.Error())
			return
		}
	}

	tx.Commit()
	jsonCreated(&c.Controller, nil)
}

// =============================================================================

// QuestionIDController defines single question specific http requests
type QuestionIDController struct {
	beego.Controller
}

// Get question by ID
func (c *QuestionIDController) Get() {
	id := c.Ctx.Input.Param(":id")

	conn := db.Conn

	var question models.Question
	if err := conn.Where("ID = ?", id).First(&question).Error; err != nil {
		badRequest(&c.Controller, fmt.Sprintf("Question not found. ID=%v", id))
		return
	}

	jsonOK(&c.Controller, question)
}
