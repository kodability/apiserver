package controllers

import (
	"encoding/json"
	"tryout-runner/db"
	"tryout-runner/models"

	"github.com/astaxie/beego"
)

type QuestionController struct {
	beego.Controller
}

type QuestionLocaleDesc struct {
	LocaleID string
	Title    string
	Desc     string
}

type QuestionLangCode struct {
	Lang     string
	InitCode string
	TestCode string
}

type QuestionBody struct {
	Desc          []QuestionLocaleDesc
	Codes         []QuestionLangCode
	Level         int
	EstimatedTime int
	Tags          string
	Demo          bool
}

func (c *QuestionController) Post() {
	var body QuestionBody
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)

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
