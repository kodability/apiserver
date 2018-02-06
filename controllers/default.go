package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "kodability"
	c.Data["Email"] = "admin@kodability.com"
	c.TplName = "index.tpl"
}
