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
	Code       string
}

// Run a tryout
func (c *RunController) RunTryout() {
	var body RunBody
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	c.Ctx.Output.SetStatus(201)

	// TODO: run a tryout
	c.Data["json"] = map[string]string{"msg": "not implemented"}
	c.ServeJSON()
}
