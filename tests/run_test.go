package tests

import (
	"testing"

	"github.com/kodability/tryout-runner/controllers"
	_ "github.com/kodability/tryout-runner/routers"

	"github.com/astaxie/beego"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRun(t *testing.T) {
	runBody := controllers.RunBody{QuestionID: 1, Code: ""}
	r, w, _ := makePostJSON("/api/v1/run", runBody)
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	Convey("POST run\n", t, func() {
		Convey("StatusCode = 201", func() {
			So(w.Code, ShouldEqual, 201)
		})
	})
}
