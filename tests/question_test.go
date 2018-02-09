package tests

import (
	"testing"

	"github.com/kodability/tryout-runner/controllers"
	"github.com/kodability/tryout-runner/db"
	"github.com/kodability/tryout-runner/models"
	_ "github.com/kodability/tryout-runner/routers"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPostQuestion(t *testing.T) {
	enDesc := controllers.QuestionLocaleDesc{
		LocaleID: "en",
		Title:    "test",
		Desc:     "test description",
	}
	koDesc := controllers.QuestionLocaleDesc{
		LocaleID: "ko",
		Title:    "테스트",
		Desc:     "테스트 설명",
	}
	javaCode := controllers.QuestionLangCode{
		Lang:     "java",
		InitCode: "public class Main {}",
		TestCode: "public class MainTest {}",
	}
	javascriptCode := controllers.QuestionLangCode{
		Lang:     "javascript",
		InitCode: "const fn = function() {}",
		TestCode: "const fnTest = function() {}",
	}

	body := controllers.QuestionBody{
		Level:         1,
		EstimatedTime: 30,
		Tags:          "tree,sort",
		Demo:          true,
		Desc:          []controllers.QuestionLocaleDesc{enDesc, koDesc},
		Codes:         []controllers.QuestionLangCode{javaCode, javascriptCode},
	}
	req, rw, _ := makePostJSON("/api/v1/question", &body)
	beego.BeeApp.Handlers.ServeHTTP(rw, req)

	Convey("POST question\n", t, func() {
		Convey("StatusCode = 201", func() {
			So(rw.Code, ShouldEqual, 201)
		})
		Convey("Inserted Question : 1", func() {
			var actual []models.Question
			db.Conn.Find(&actual)
			So(len(actual), ShouldEqual, 1)
		})
	})
}
