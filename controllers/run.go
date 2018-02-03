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
	code       string
}

func (c *RunController) Post() {
	var body RunBody
	json.Unmarshal(c.Ctx.Input.RequestBody, &body)
	c.Data["json"] = map[string]string{"msg": "not implemented"}
	c.ServeJSON()
}
