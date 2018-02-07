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
	c.Ctx.Output.SetStatus(201)

	// Insert question
	conn := db.Conn
	question := models.Question{
		Level:         body.Level,
		EstimatedTime: body.EstimatedTime,
		Tags:          body.Tags,
		Demo:          body.Demo,
	}
	conn.Create(&question)

	// Insert QuestionDesc
	// TODO

	// Insert QuestionCode
	// TODO

	c.Data["json"] = map[string]string{"msg": "not implemented"}
	c.ServeJSON()
}
