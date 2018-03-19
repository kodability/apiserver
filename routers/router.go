package routers

import (
	beecontext "github.com/astaxie/beego/context"

	"github.com/astaxie/beego"
	"github.com/kodability/apiserver/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	qc := &controllers.QuestionController{}
	rc := &controllers.RunController{}

	ns := beego.NewNamespace("/api/v1",
		beego.NSCond(func(ctx *beecontext.Context) bool {
			return true
		}),
		beego.NSNamespace("/question",
			beego.NSRouter("/", qc, "post:AddQuestion"),
			beego.NSRouter("/:id:int", qc, "get:GetQuestionByID;delete:DeleteQuestionByID;put:UpdateQuestion"),
			beego.NSRouter("/:id:int/code", qc, "post:AddQuestionCode"),
			beego.NSRouter("/:id:int/code/:lang", qc, "get:GetQuestionCode;delete:DeleteQuestionCode;put:UpdateQuestionCode"),
		),
		beego.NSNamespace("/run",
			beego.NSRouter("/", rc, "post:RunTryout"),
		),
	)
	beego.AddNamespace(ns)
}
