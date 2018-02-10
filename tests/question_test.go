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

func deleteQuestions() {
	conn := db.Conn
	conn.Delete(models.Question{})
	conn.Delete(models.QuestionDescription{})
	conn.Delete(models.QuestionCode{})
}

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

	Convey("POST question", t, func() {
		deleteQuestions()
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

		Convey("StatusCode = 201", func() {
			So(rw.Code, ShouldEqual, 201)
		})
		Convey("Inserted Question", func() {
			var questions []models.Question
			db.Conn.Find(&questions)
			So(questions, ShouldHaveLength, 1)

			question := questions[0]
			So(map[string]interface{}{
				"Level":         question.Level,
				"EstimatedTime": question.EstimatedTime,
				"Tags":          question.Tags,
			}, ShouldResemble, map[string]interface{}{
				"Level":         body.Level,
				"EstimatedTime": body.EstimatedTime,
				"Tags":          body.Tags,
			})
		})
		Convey("Inserted QuestionDescription", func() {
			var questionDescriptions []models.QuestionDescription
			db.Conn.Find(&questionDescriptions)
			So(questionDescriptions, ShouldHaveLength, 2)
		})
		Convey("Inserted QuestionCode", func() {
			var questionCodes []models.QuestionCode
			db.Conn.Find(&questionCodes)
			So(questionCodes, ShouldHaveLength, 2)
		})
	})

	Convey("POST question with empty desc", t, func() {
		deleteQuestions()
		body := controllers.QuestionBody{
			Level:         1,
			EstimatedTime: 30,
			Tags:          "tree,sort",
			Demo:          true,
			Desc:          []controllers.QuestionLocaleDesc{},
			Codes:         []controllers.QuestionLangCode{javaCode, javascriptCode},
		}
		req, rw, _ := makePostJSON("/api/v1/question", &body)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 400", func() {
			So(rw.Code, ShouldEqual, 400)
		})
	})
}
