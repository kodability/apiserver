package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
)

// RunController defines operations about tryout run.
type RunController struct {
	beego.Controller
}

type RunBody struct {
	QuestionID int
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
	c.Ctx.Output.SetStatus(201)

	// TODO: run a tryout
	c.Data["json"] = map[string]string{"msg": "not implemented"}
	c.ServeJSON()
}
