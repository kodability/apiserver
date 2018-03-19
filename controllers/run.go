package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/kodability/apiserver/db"
	"github.com/kodability/apiserver/models"
)

// RunController defines operations about tryout run.
type RunController struct {
	beego.Controller
}

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

	c.Ctx.Output.SetStatus(201)

	// TODO: run a tryout
	c.Data["json"] = map[string]string{"msg": "not implemented"}
	c.ServeJSON()
}
