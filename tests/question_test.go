package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	c "github.com/kodability/tryout-runner/controllers"
	"github.com/kodability/tryout-runner/db"
	m "github.com/kodability/tryout-runner/models"
	_ "github.com/kodability/tryout-runner/routers"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func deleteQuestionsAndDescAndCodes() {
	conn := db.Conn
	conn.Unscoped().Delete(m.Question{})
	conn.Unscoped().Delete(m.QuestionDescription{})
	conn.Unscoped().Delete(m.QuestionCode{})
}

func TestPostQuestion(t *testing.T) {
	enDesc := c.QuestionLocaleDesc{
		LocaleID: "en",
		Title:    "test",
		Desc:     "test description",
	}
	koDesc := c.QuestionLocaleDesc{
		LocaleID: "ko",
		Title:    "테스트",
		Desc:     "테스트 설명",
	}
	javaCode := c.QuestionLangCode{
		Lang:     "java",
		InitCode: "public class Main {}",
		TestCode: "public class MainTest {}",
	}
	javascriptCode := c.QuestionLangCode{
		Lang:     "javascript",
		InitCode: "const fn = function() {}",
		TestCode: "const fnTest = function() {}",
	}

	Convey("POST question", t, func() {
		deleteQuestionsAndDescAndCodes()
		body := c.QuestionPostBody{
			Level:         1,
			EstimatedTime: 30,
			Tags:          "tree,sort",
			Demo:          true,
			Desc:          []c.QuestionLocaleDesc{enDesc, koDesc},
			Codes:         []c.QuestionLangCode{javaCode, javascriptCode},
		}
		req, rw, _ := makePostJSON("/api/v1/question", &body)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 201", func() {
			So(rw.Code, ShouldEqual, http.StatusCreated)
		})
		Convey("Inserted Question", func() {
			var questions []m.Question
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
			var questionDescriptions []m.QuestionDescription
			db.Conn.Find(&questionDescriptions)
			So(questionDescriptions, ShouldHaveLength, 2)
		})
		Convey("Inserted QuestionCode", func() {
			var questionCodes []m.QuestionCode
			db.Conn.Find(&questionCodes)
			So(questionCodes, ShouldHaveLength, 2)
		})
	})

	Convey("POST question with empty desc", t, func() {
		deleteQuestionsAndDescAndCodes()
		body := c.QuestionPostBody{
			Level:         1,
			EstimatedTime: 30,
			Tags:          "tree,sort",
			Demo:          true,
			Desc:          []c.QuestionLocaleDesc{},
			Codes:         []c.QuestionLangCode{javaCode, javascriptCode},
		}
		req, rw, _ := makePostJSON("/api/v1/question", &body)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 400", func() {
			So(rw.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("POST question with empty code", t, func() {
		deleteQuestionsAndDescAndCodes()
		body := c.QuestionPostBody{
			Level:         1,
			EstimatedTime: 30,
			Tags:          "tree,sort",
			Demo:          true,
			Desc:          []c.QuestionLocaleDesc{enDesc},
			Codes:         []c.QuestionLangCode{},
		}
		req, rw, _ := makePostJSON("/api/v1/question", &body)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 201", func() {
			So(rw.Code, ShouldEqual, http.StatusCreated)
		})
		Convey("Inserted QuestionDescription", func() {
			var questionDescriptions []m.QuestionDescription
			db.Conn.Find(&questionDescriptions)
			So(questionDescriptions, ShouldHaveLength, 1)
		})
	})

	Convey("POST question : code insert error", t, func() {
		deleteQuestionsAndDescAndCodes()
		body := c.QuestionPostBody{
			Level:         1,
			EstimatedTime: 30,
			Tags:          "tree,sort",
			Demo:          true,
			Desc:          []c.QuestionLocaleDesc{enDesc},
			Codes:         []c.QuestionLangCode{javaCode, javaCode},
		}
		req, rw, _ := makePostJSON("/api/v1/question", &body)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 500", func() {
			So(rw.Code, ShouldEqual, http.StatusInternalServerError)
		})
		Convey("Empty QuestionCode and QuestionDesc", func() {
			var questionCodes []m.QuestionCode
			var questionDescriptions []m.QuestionDescription
			db.Conn.Find(&questionCodes)
			db.Conn.Find(&questionDescriptions)
			So(questionCodes, ShouldHaveLength, 0)
			So(questionDescriptions, ShouldHaveLength, 0)
		})
	})

	Convey("POST question : duplicated desc insert error", t, func() {
		deleteQuestionsAndDescAndCodes()
		body := c.QuestionPostBody{
			Level:         1,
			EstimatedTime: 30,
			Tags:          "tree,sort",
			Demo:          true,
			Desc:          []c.QuestionLocaleDesc{enDesc, enDesc},
			Codes:         []c.QuestionLangCode{},
		}
		req, rw, _ := makePostJSON("/api/v1/question", &body)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 500", func() {
			So(rw.Code, ShouldEqual, http.StatusInternalServerError)
		})
		Convey("Empty QuestionDesc", func() {
			var questionDescriptions []m.QuestionDescription
			db.Conn.Find(&questionDescriptions)
			So(questionDescriptions, ShouldHaveLength, 0)
		})
	})
}

func TestGetQuestionByID(t *testing.T) {
	deleteQuestionsAndDescAndCodes()

	question := m.Question{}
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
			var result m.Question
			json.Unmarshal(rw.Body.Bytes(), &result)
			So(result.ID, ShouldEqual, question.ID)
		})
	})
}

func TestDeleteQuestionByID(t *testing.T) {
	deleteQuestionsAndDescAndCodes()

	question := m.Question{}
	db.Conn.Create(&question)

	koDesc := m.QuestionDescription{
		QuestionID: question.ID,
		LocaleID:   "ko",
		Title:      "테스트",
	}
	javaCode := m.QuestionCode{
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

	question := m.Question{}
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
			"level":         2,
			"estimatedTime": 10,
			"tags":          "A,B",
			"demo":          true,
		}
		req, rw, _ := makePutJSON(fmt.Sprintf("/api/v1/question/%v", question.ID), body)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 200", func() {
			So(rw.Code, ShouldEqual, http.StatusOK)

			var actual m.Question
			db.Conn.Where("id = ?", question.ID).First(&actual)
			So(actual.Level, ShouldEqual, body["level"])
			So(actual.EstimatedTime, ShouldEqual, body["estimatedTime"])
			So(actual.Tags, ShouldEqual, body["tags"])
			So(actual.Demo, ShouldEqual, body["demo"])
		})
	})
}

func TestGetQuestionCode(t *testing.T) {
	deleteQuestionsAndDescAndCodes()

	javaCode := m.QuestionCode{
		QuestionID: 1,
		Lang:       "java",
		InitCode:   "public class Main {}",
		TestCode:   "public class MainTest {}",
	}
	db.Conn.Create(&javaCode)

	Convey("GET : invalid lang", t, func() {
		req, rw, _ := makeGet("/api/v1/question/9999/code/cpp")
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 400", func() {
			So(rw.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("GET : get QuestionCode", t, func() {
		req, rw, _ := makeGet("/api/v1/question/1/code/java")
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 200", func() {
			So(rw.Code, ShouldEqual, http.StatusOK)

			var actual m.QuestionCode
			db.Conn.Where("question_id = ? AND lang= ?", javaCode.QuestionID, javaCode.Lang).First(&actual)
			So(actual.QuestionID, ShouldEqual, javaCode.QuestionID)
			So(actual.Lang, ShouldEqual, javaCode.Lang)
			So(actual.InitCode, ShouldEqual, javaCode.InitCode)
			So(actual.TestCode, ShouldEqual, javaCode.TestCode)
		})
	})
}

func TestDeleteQuestionCode(t *testing.T) {
	deleteQuestionsAndDescAndCodes()

	javaCode := m.QuestionCode{
		QuestionID: 1,
		Lang:       "java",
		InitCode:   "public class Main {}",
		TestCode:   "public class MainTest {}",
	}
	db.Conn.Create(&javaCode)

	Convey("DELETE : invalid lang", t, func() {
		req, rw, _ := makeDelete("/api/v1/question/1/code/cpp")
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 400", func() {
			So(rw.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("DELETE : delete QuestionCode", t, func() {
		req, rw, _ := makeDelete("/api/v1/question/1/code/java")
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 204", func() {
			So(rw.Code, ShouldEqual, http.StatusNoContent)

			var questionCodes []m.QuestionCode
			var count uint
			db.Conn.Where("question_id = ? AND lang= ?", javaCode.QuestionID, javaCode.Lang).Find(&questionCodes).Count(&count)
			So(count, ShouldEqual, 0)
		})
	})
}

func TestPostQuestionCode(t *testing.T) {
	deleteQuestionsAndDescAndCodes()

	// Add a question
	question := m.Question{}
	db.Conn.Create(&question)

	Convey("POST : invalid questionID", t, func() {
		req, rw, _ := makePostJSON("/api/v1/question/9999/code", nil)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 400", func() {
			So(rw.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("POST : add QuestionCode", t, func() {
		questionID := question.ID
		body := map[string]interface{}{
			"lang":     "java",
			"initCode": "init code",
			"testCode": "test code",
		}
		req, rw, _ := makePostJSON(fmt.Sprintf("/api/v1/question/%d/code", questionID), body)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 201", func() {
			So(rw.Code, ShouldEqual, http.StatusCreated)

			var actual m.QuestionCode
			db.Conn.Where("question_id = ? AND lang= ?", questionID, body["lang"]).First(&actual)
			So(actual.InitCode, ShouldEqual, body["initCode"])
			So(actual.TestCode, ShouldEqual, body["testCode"])
		})
	})
}

func TestPutQuestionCode(t *testing.T) {
	deleteQuestionsAndDescAndCodes()

	javaCode := m.QuestionCode{
		QuestionID: 1,
		Lang:       "java",
		InitCode:   "public class Main {}",
		TestCode:   "public class MainTest {}",
	}
	db.Conn.Create(&javaCode)

	Convey("PUT : invalid lang", t, func() {
		req, rw, _ := makePutJSON("/api/v1/question/1/code/cpp", nil)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 400", func() {
			So(rw.Code, ShouldEqual, http.StatusBadRequest)
		})
	})

	Convey("PUT : update QuestionCode", t, func() {
		body := map[string]interface{}{
			"initCode": "init code",
			"testCode": "test code",
		}
		req, rw, _ := makePutJSON("/api/v1/question/1/code/java", body)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("StatusCode = 200", func() {
			So(rw.Code, ShouldEqual, http.StatusOK)

			var actual m.QuestionCode
			db.Conn.Where("question_id = ? AND lang= ?", javaCode.QuestionID, javaCode.Lang).First(&actual)
			So(actual.InitCode, ShouldEqual, body["initCode"])
			So(actual.TestCode, ShouldEqual, body["testCode"])
		})
	})
}
