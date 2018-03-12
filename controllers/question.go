package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/kodability/tryout-runner/db"
	m "github.com/kodability/tryout-runner/models"
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

// QuestionPostBody is a body struct for Post
type QuestionPostBody struct {
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

// AddQuestion add a new question
func (c *QuestionController) AddQuestion() {
	var body QuestionPostBody
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		badRequest(&c.Controller, err.Error())
		return
	}

	// Validate
	if len(body.Desc) == 0 {
		badRequest(&c.Controller, "desc cannot be empty")
		return
	}

	conn := db.Conn
	tx := conn.Begin()

	// QuestionDescription
	var descriptions []m.QuestionDescription
	for _, desc := range body.Desc {
		questionDesc := m.QuestionDescription{
			LocaleID:    desc.LocaleID,
			Title:       desc.Title,
			Description: desc.Desc,
		}
		descriptions = append(descriptions, questionDesc)
	}

	// QuestionCode
	var codes []m.QuestionCode
	for _, code := range body.Codes {
		questionCode := m.QuestionCode{
			Lang:     code.Lang,
			InitCode: code.InitCode,
			TestCode: code.TestCode,
		}
		codes = append(codes, questionCode)
	}

	// Insert question
	question := m.Question{
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

// GetQuestionByID returns a question by ID
func (c *QuestionController) GetQuestionByID() {
	id := c.Ctx.Input.Param(":id")

	conn := db.Conn

	// Find Question
	var question m.Question
	if err := conn.Where("id = ?", id).First(&question).Error; err != nil {
		badRequest(&c.Controller, fmt.Sprintf("Question not found. id=%v", id))
		return
	}

	jsonOK(&c.Controller, question)
}

// DeleteQuestionByID deletes a question by ID
func (c *QuestionController) DeleteQuestionByID() {
	id := c.Ctx.Input.Param(":id")

	var err error
	conn := db.Conn
	tx := conn.Begin()

	// Delete Question
	err = tx.Unscoped().Where("id = ?", id).Delete(m.Question{}).Error
	if err != nil {
		tx.Rollback()
		internalServerError(&c.Controller, fmt.Sprintf("Failed to delete Question: %v", err.Error()))
		return
	}

	// Delete QuestionDescription
	err = tx.Unscoped().Where("question_id = ?", id).Delete(m.QuestionDescription{}).Error
	if err != nil {
		tx.Rollback()
		internalServerError(&c.Controller, fmt.Sprintf("Failed to delete QuestionDescription: %v", err.Error()))
		return
	}

	// Delete QuestionCode
	err = tx.Unscoped().Where("question_id = ?", id).Delete(m.QuestionCode{}).Error
	if err != nil {
		tx.Rollback()
		internalServerError(&c.Controller, fmt.Sprintf("Failed to delete QuestionCode: %v", err.Error()))
		return
	}

	tx.Commit()

	noContent(&c.Controller, nil)
}

type QuestionPutBody struct {
	Level         *int
	EstimatedTime *int
	Tags          *string
	Demo          *bool
}

// UpdateQuestion updates a question
func (c *QuestionController) UpdateQuestion() {
	conn := db.Conn
	id := c.Ctx.Input.Param(":id")

	var body QuestionPutBody
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &body); err != nil {
		badRequest(&c.Controller, "Failed to parse body")
		return
	}

	// Find Question
	var question m.Question
	if err := conn.Where("id = ?", id).First(&question).Error; err != nil {
		badRequest(&c.Controller, fmt.Sprintf("Question not found. id=%v", id))
		return
	}

	// Update fields
	if body.Level != nil {
		question.Level = *body.Level
	}
	if body.EstimatedTime != nil {
		question.EstimatedTime = *body.EstimatedTime
	}
	if body.Tags != nil {
		question.Tags = *body.Tags
	}
	if body.Demo != nil {
		question.Demo = *body.Demo
	}
	conn.Model(&question).Updates(body)

	setStatusOK(&c.Controller)
}

// GetQuestionCode returns a question code
func (c *QuestionController) GetQuestionCode() {
	conn := db.Conn
	questionID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	lang := c.Ctx.Input.Param(":lang")

	var err error
	var questionCode m.QuestionCode
	err = conn.Where("question_id = ? AND lang = ?", questionID, lang).First(&questionCode).Error
	if err != nil {
		badRequest(&c.Controller, fmt.Sprintf("QuestionCode not found: %v", err.Error()))
		return
	}

	jsonOK(&c.Controller, questionCode)
}

// AddQuestionCode adds a new question code
func (c *QuestionController) AddQuestionCode() {
	var err error

	conn := db.Conn
	questionID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))

	// Find question
	var question m.Question
	err = conn.Where("id = ?", questionID).First(&question).Error
	if err != nil {
		badRequest(&c.Controller, fmt.Sprintf("Question not found. id=%v", questionID))
		return
	}

	// Parse body
	var questionCode m.QuestionCode
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &questionCode)
	if err != nil {
		badRequest(&c.Controller, fmt.Sprintf("Failed to parse body : %v", err.Error()))
		return
	}
	questionCode.QuestionID = uint(questionID)

	// Insert
	err = conn.Create(&questionCode).Error
	if err != nil {
		internalServerError(&c.Controller, fmt.Sprintf("Failed to add QuestionCode: %v", err.Error()))
		return
	}

	jsonCreated(&c.Controller, questionCode)
}

// UpdateQuestionCode updates a question code
func (c *QuestionController) UpdateQuestionCode() {
	var err error

	conn := db.Conn
	questionID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	lang := c.Ctx.Input.Param(":lang")

	var body map[string]interface{}
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	if err != nil {
		badRequest(&c.Controller, fmt.Sprintf("Failed to parse body : %v", err.Error()))
		return
	}

	var questionCode m.QuestionCode
	err = conn.Where("question_id = ? AND lang = ?", questionID, lang).First(&questionCode).Error
	if err != nil {
		badRequest(&c.Controller, fmt.Sprintf("QuestionCode not found: %v", err.Error()))
		return
	}

	// Update question code
	err = conn.Model(&questionCode).Where("question_id = ? AND lang = ?", questionID, lang).Updates(body).Error
	if err != nil {
		internalServerError(&c.Controller, fmt.Sprintf("Failed to update QuestionCode: %v", err.Error()))
		return
	}

	setStatusOK(&c.Controller)
}

// DeleteQuestionCode deletes a question code
func (c *QuestionController) DeleteQuestionCode() {
	var err error
	conn := db.Conn
	questionID, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	lang := c.Ctx.Input.Param(":lang")

	// Find matching QuestionCode
	var questionCode m.QuestionCode
	err = conn.Where("question_id = ? AND lang = ?", questionID, lang).First(&questionCode).Error
	if err != nil {
		badRequest(&c.Controller, fmt.Sprintf("QuestionCode not found (questionId=%v, lang=%v). %v", questionID, lang, err.Error()))
		return
	}

	// Delete QuestionCode
	err = conn.Delete(&questionCode).Error
	if err != nil {
		internalServerError(&c.Controller, fmt.Sprintf("Failed to delete QuestionCode: %v", err.Error()))
		return
	}

	noContent(&c.Controller, nil)
}
