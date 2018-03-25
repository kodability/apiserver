package tests

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/kodability/apiserver/services/run"

	c "github.com/kodability/apiserver/controllers"
	"github.com/kodability/apiserver/db"
	m "github.com/kodability/apiserver/models"
	_ "github.com/kodability/apiserver/routers"

	"github.com/astaxie/beego"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	. "github.com/smartystreets/goconvey/convey"
)

func deleteTryouts() {
	conn := db.Conn
	conn.Unscoped().Delete(m.Tryout{})
}

func createGroovyQuestionCode() m.QuestionCode {
	return m.QuestionCode{
		Lang:     "groovy",
		InitCode: "",
		TestCode: `
import static org.testng.AssertJUnit.*
import org.testng.annotations.*

class TestExample {
	@Test(timeOut= 1000L)
	void test1() {
		assertEquals(55, new Example().sum(1, 10))
	}
	@Test
	void test2() {
		assertEquals(1, new Example().sum(1, 1))
	}
}`,
	}
}

func TestRun(t *testing.T) {
	deleteQuestionsAndDescAndCodes()

	if beego.AppConfig.String("tryout.runner") == "mock" {
		c.SetTryoutRunner(&run.TryoutMockRunner{
			Err: nil,
			Result: &run.JUnitReport{
				Tests:       2,
				ElapsedTime: 1,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test1"},
					run.JUnitTestcaseResult{Name: "test1"},
				},
			},
		})
	}

	// Create a question
	groovyQuestionCode := createGroovyQuestionCode()
	question := m.Question{
		Codes: []m.QuestionCode{groovyQuestionCode},
	}
	db.Conn.Create(&question)

	Convey("When POST run", t, func() {
		deleteTryouts()

		runBody := c.RunBody{
			QuestionID: question.ID,
			Lang:       groovyQuestionCode.Lang,
			Code: `
			class Example {
				int sum(int from, int to) {
					(from + to) * (to - from + 1) / 2
				}
			}`,
		}
		req, rw, _ := makePostJSON("/api/v1/run", runBody)
		beego.BeeApp.Handlers.ServeHTTP(rw, req)

		Convey("Then StatusCode = 201 & TryoutResult response", func() {
			log.Println(rw)
			So(rw.Code, ShouldEqual, 201)
			var result m.TryoutResult
			json.Unmarshal(rw.Body.Bytes(), &result)
			So(map[string]interface{}{
				"TestCount":    result.TestCount,
				"ErrorCount":   result.ErrorCount,
				"FailureCount": result.FailureCount,
				"ErrorMsg":     "",
			}, ShouldResemble, map[string]interface{}{
				"TestCount":    2,
				"ErrorCount":   0,
				"FailureCount": 0,
				"ErrorMsg":     "",
			})
		})
		Convey("Then Tryout inserted", func() {
			var tryouts []m.Tryout
			db.Conn.Find(&tryouts)
			So(tryouts, ShouldHaveLength, 1)

			tryout := tryouts[0]
			So(map[string]interface{}{
				"QuestionID": tryout.QuestionID,
				"Lang":       tryout.Lang,
				"Code":       tryout.Code,
			}, ShouldResemble, map[string]interface{}{
				"QuestionID": question.ID,
				"Lang":       groovyQuestionCode.Lang,
				"Code":       runBody.Code,
			})
		})
	})
}
