package routers

import (
	"tryout-runner/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/run",
			beego.NSInclude(
				&controllers.RunController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
