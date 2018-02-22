package controllers

import (
	"github.com/astaxie/beego"
)

// MainController response '/'
type MainController struct {
	beego.Controller
}

// Get repoponses 'GET /'
func (c *MainController) Get() {
	c.Data["Website"] = "kodability"
	c.Data["Email"] = "admin@kodability.com"
	c.TplName = "index.tpl"
}
