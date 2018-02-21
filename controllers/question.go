package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/kodability/tryout-runner/db"
	"github.com/kodability/tryout-runner/models"

	"github.com/astaxie/beego"
)

// =============================================================================

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

	// QuestionDescription
	var descriptions []models.QuestionDescription
	for _, desc := range body.Desc {
		questionDesc := models.QuestionDescription{
			LocaleID:    desc.LocaleID,
			Title:       desc.Title,
			Description: desc.Desc,
		}
		descriptions = append(descriptions, questionDesc)
	}

	// QuestionCode
	var codes []models.QuestionCode
	for _, code := range body.Codes {
		questionCode := models.QuestionCode{
			Lang:     code.Lang,
			InitCode: code.InitCode,
			TestCode: code.TestCode,
		}
		codes = append(codes, questionCode)
	}

	// Insert question
	question := models.Question{
		Level:         body.Level,
		EstimatedTime: body.EstimatedTime,
		Desctiptions:  descriptions,
		Codes:         codes,
		Tags:          body.Tags,
		Demo:          body.Demo,
	}
	if err := tx.Create(&question).Error; err != nil {
		tx.Rollback()
		internalServerError(&c.Controller, err.Error())
		return
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

	// Find Question
	var question models.Question
	if err := conn.Where("id = ?", id).First(&question).Error; err != nil {
		badRequest(&c.Controller, fmt.Sprintf("Question not found. id=%v", id))
		return
	}

	jsonOK(&c.Controller, question)
}

// Delete question by ID
func (c *QuestionIDController) Delete() {
	id := c.Ctx.Input.Param(":id")

	conn := db.Conn
	tx := conn.Begin()

	// Delete Question
	if err := tx.Unscoped().Where("id = ?", id).Delete(models.Question{}).Error; err != nil {
		tx.Rollback()
		internalServerError(&c.Controller, fmt.Sprintf("Failed to delete Question: %v", err.Error()))
		return
	}

	// Delete QuestionDescription
	if err := tx.Unscoped().Where("question_id = ?", id).Delete(models.QuestionDescription{}).Error; err != nil {
		tx.Rollback()
		internalServerError(&c.Controller, fmt.Sprintf("Failed to delete QuestionDescription: %v", err.Error()))
		return
	}

	// Delete QuestionCode
	if err := tx.Unscoped().Where("question_id = ?", id).Delete(models.QuestionCode{}).Error; err != nil {
		tx.Rollback()
		internalServerError(&c.Controller, fmt.Sprintf("Failed to delete QuestionCode: %v", err.Error()))
		return
	}

	tx.Commit()

	noContent(&c.Controller, nil)
}

// Put updates a question by ID
func (c *QuestionIDController) Put() {
	id := c.Ctx.Input.Param(":id")

	var body map[string]interface{}
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

	conn := db.Conn

	// Find Question
	var question models.Question
	if err := conn.Where("id = ?", id).First(&question).Error; err != nil {
		badRequest(&c.Controller, fmt.Sprintf("Question not found. id=%v", id))
		return
	}

	if len(body) > 0 {
		conn.Model(&question).Updates(body)
	}

	setStatusOK(&c.Controller)
}
