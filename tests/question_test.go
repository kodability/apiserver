package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/kodability/tryout-runner/controllers"
	"github.com/kodability/tryout-runner/db"
	"github.com/kodability/tryout-runner/models"
	_ "github.com/kodability/tryout-runner/routers"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func deleteQuestionsAndDescAndCodes() {
	conn := db.Conn
	conn.Unscoped().Delete(models.Question{})
	conn.Unscoped().Delete(models.QuestionDescription{})
	conn.Unscoped().Delete(models.QuestionCode{})
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
		deleteQuestionsAndDescAndCodes()
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
			So(rw.Code, ShouldEqual, http.StatusCreated)
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
		deleteQuestionsAndDescAndCodes()
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
			So(rw.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("POST question with empty code", t, func() {
		deleteQuestionsAndDescAndCodes()
		body := controllers.QuestionBody{
			Level:         1,
			EstimatedTime: 30,
			Tags:          "tree,sort",
			Demo:          true,
			Desc:          []controllers.QuestionLocaleDesc{enDesc},
			Codes:         []controllers.QuestionLangCode{},
		}
		req, rw, _ := makePostJSON("/api/v1/question", &body)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 201", func() {
			So(rw.Code, ShouldEqual, http.StatusCreated)
		})
		Convey("Inserted QuestionDescription", func() {
			var questionDescriptions []models.QuestionDescription
			db.Conn.Find(&questionDescriptions)
			So(questionDescriptions, ShouldHaveLength, 1)
		})
	})

	Convey("POST question : code insert error", t, func() {
		deleteQuestionsAndDescAndCodes()
		body := controllers.QuestionBody{
			Level:         1,
			EstimatedTime: 30,
			Tags:          "tree,sort",
			Demo:          true,
			Desc:          []controllers.QuestionLocaleDesc{enDesc},
			Codes:         []controllers.QuestionLangCode{javaCode, javaCode},
		}
		req, rw, _ := makePostJSON("/api/v1/question", &body)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 500", func() {
			So(rw.Code, ShouldEqual, http.StatusInternalServerError)
		})
		Convey("Empty QuestionCode and QuestionDesc", func() {
			var questionCodes []models.QuestionCode
			var questionDescriptions []models.QuestionDescription
			db.Conn.Find(&questionCodes)
			db.Conn.Find(&questionDescriptions)
			So(questionCodes, ShouldHaveLength, 0)
			So(questionDescriptions, ShouldHaveLength, 0)
		})
	})

	Convey("POST question : duplicated desc insert error", t, func() {
		deleteQuestionsAndDescAndCodes()
		body := controllers.QuestionBody{
			Level:         1,
			EstimatedTime: 30,
			Tags:          "tree,sort",
			Demo:          true,
			Desc:          []controllers.QuestionLocaleDesc{enDesc, enDesc},
			Codes:         []controllers.QuestionLangCode{},
		}
		req, rw, _ := makePostJSON("/api/v1/question", &body)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 500", func() {
			So(rw.Code, ShouldEqual, http.StatusInternalServerError)
		})
		Convey("Empty QuestionDesc", func() {
			var questionDescriptions []models.QuestionDescription
			db.Conn.Find(&questionDescriptions)
			So(questionDescriptions, ShouldHaveLength, 0)
		})
	})
}

func TestGetQuestionByID(t *testing.T) {
	deleteQuestionsAndDescAndCodes()

	question := models.Question{}
	db.Conn.Create(&question)

	Convey("GET : invalid ID", t, func() {
		req, rw, _ := makeGet("/api/v1/question/0")
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 400", func() {
			So(rw.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("GET : valid ID", t, func() {
		req, rw, _ := makeGet(fmt.Sprintf("/api/v1/question/%v", question.ID))
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 200", func() {
			So(rw.Code, ShouldEqual, http.StatusOK)
			var result models.Question
			json.Unmarshal(rw.Body.Bytes(), &result)
			So(result.ID, ShouldEqual, question.ID)
		})
	})
}

func TestDeleteQuestionByID(t *testing.T) {
	deleteQuestionsAndDescAndCodes()

	question := models.Question{}
	db.Conn.Create(&question)

	koDesc := models.QuestionDescription{
		QuestionID: question.ID,
		LocaleID:   "ko",
		Title:      "테스트",
	}
	javaCode := models.QuestionCode{
		QuestionID: question.ID,
		Lang:       "java",
	}
	db.Conn.Create(&koDesc)
	db.Conn.Create(&javaCode)

	Convey("DELETE : invalid ID", t, func() {
		req, rw, _ := makeDelete("/api/v1/question/0")
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 204", func() {
			log.Println(rw.Body)
			So(rw.Code, ShouldEqual, http.StatusNoContent)
		})
	})

	Convey("DELETE : valid ID", t, func() {
		req, rw, _ := makeDelete(fmt.Sprintf("/api/v1/question/%v", question.ID))
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 204", func() {
			log.Println(rw.Body)
			So(rw.Code, ShouldEqual, http.StatusNoContent)
		})
	})
}

func TestPutQuestionByID(t *testing.T) {
	deleteQuestionsAndDescAndCodes()

	question := models.Question{}
	db.Conn.Create(&question)

	Convey("PUT : invalid ID", t, func() {
		req, rw, _ := makePutJSON("/api/v1/question/0", nil)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 400", func() {
			So(rw.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("PUT : valid ID", t, func() {
		body := map[string]interface{}{
			"level":          2,
			"estimated_time": 10,
			"tags":           "A,B",
			"demo":           true,
		}
		req, rw, _ := makePutJSON(fmt.Sprintf("/api/v1/question/%v", question.ID), body)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 200", func() {
			So(rw.Code, ShouldEqual, http.StatusOK)

			var actual models.Question
			db.Conn.Where("id = ?", question.ID).First(&actual)
			So(actual.Level, ShouldEqual, body["level"])
			So(actual.EstimatedTime, ShouldEqual, body["estimated_time"])
			So(actual.Tags, ShouldEqual, body["tags"])
			So(actual.Demo, ShouldEqual, body["demo"])
		})
	})
}
