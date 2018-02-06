package routers

import (
	"tryout-runner/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/v1/run", &controllers.RunController{})
}
