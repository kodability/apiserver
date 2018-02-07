package test

import (
	"testing"
	"tryout-runner/controllers"
	_ "tryout-runner/routers"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPostQuestion(t *testing.T) {
	body := controllers.QuestionBody{
		Level:         1,
		EstimatedTime: 30,
		Tags:          "tree,sort",
		Demo:          true,
	}
	req, rw, _ := makePostJSON("/api/v1/question", &body)
	beego.BeeApp.Handlers.ServeHTTP(rw, req)

	Convey("POST question\n", t, func() {
		Convey("StatusCode = 201", func() {
			So(rw.Code, ShouldEqual, 201)
		})
	})
}
