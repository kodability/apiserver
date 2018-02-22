package routers

import (
	beecontext "github.com/astaxie/beego/context"

	. "github.com/kodability/tryout-runner/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &MainController{})

	qc := &QuestionController{}
	rc := &RunController{}

	ns := beego.NewNamespace("/api/v1",
		beego.NSCond(func(ctx *beecontext.Context) bool {
			return true
		}),
		beego.NSNamespace("/question",
			beego.NSRouter("/", qc, "post:AddQuestion"),
			beego.NSRouter("/:id:int", qc, "get:GetQuestionByID;delete:DeleteQuestionByID;put:UpdateQuestion"),
		),
		beego.NSNamespace("/run",
			beego.NSRouter("/", rc, "post:RunTryout"),
		),
	)
	beego.AddNamespace(ns)
}
