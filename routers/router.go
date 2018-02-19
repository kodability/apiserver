package routers

import (
	"github.com/kodability/tryout-runner/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/v1/question", &controllers.QuestionController{})
	beego.Router("/api/v1/question/:id:int", &controllers.QuestionIDController{})
	beego.Router("/api/v1/run", &controllers.RunController{})
}
