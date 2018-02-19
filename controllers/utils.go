package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
)

func setStatus(c *beego.Controller, statusCode int) {
	c.Ctx.Output.SetStatus(statusCode)
}

func jsonWithStatus(c *beego.Controller, value interface{}, statusCode int) {
	c.Data["json"] = value
	setStatus(c, statusCode)
	c.ServeJSON()
}

func setStatusOK(c *beego.Controller) {
	setStatus(c, http.StatusOK)
}

func jsonOK(c *beego.Controller, value interface{}) {
	jsonWithStatus(c, value, http.StatusOK)
}

func jsonCreated(c *beego.Controller, value interface{}) {
	jsonWithStatus(c, value, http.StatusCreated)
}

func noContent(c *beego.Controller, value interface{}) {
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
