package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
)

func setStatus(c *beego.Controller, statusCode int) {
	c.Ctx.Output.SetStatus(statusCode)
}

func setStatusOK(c *beego.Controller) {
	setStatus(c, http.StatusOK)
}

func jsonCreated(c *beego.Controller, value interface{}) {
	c.Data["json"] = value
	setStatus(c, http.StatusCreated)
	c.ServeJSON()
}

func jsonNoContent(c *beego.Controller, value interface{}) {
	setStatus(c, http.StatusNoContent)
}

func badRequest(c *beego.Controller, msg string) {
	c.CustomAbort(http.StatusBadRequest, msg)
}

func internalServerError(c *beego.Controller, msg string) {
	c.CustomAbort(http.StatusInternalServerError, msg)
}

func notImplemented(c *beego.Controller) {
	c.CustomAbort(501, "not implemented")
}
