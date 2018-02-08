package controllers

import "github.com/astaxie/beego"

func setStatus(c *beego.Controller, statusCode int) {
	c.Ctx.Output.SetStatus(statusCode)
}

func setStatusOK(c *beego.Controller) {
	setStatus(c, 200)
}

func jsonCreated(c *beego.Controller, value interface{}) {
	c.Data["json"] = value
	setStatus(c, 201)
	c.ServeJSON()
}

func internalServerError(c *beego.Controller, msg string) {
	c.CustomAbort(500, msg)
}

func notImplemented(c *beego.Controller) {
	c.CustomAbort(501, "not implemented")
}
