package tests

import (
	"fmt"
	"testing"

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
	void test2() {
		assertEquals(1, new Example().sum(1, 1))
	}
}`,
	}
}

func TestRun(t *testing.T) {
	deleteQuestionsAndDescAndCodes()

	// Create a question
	groovyQuestionCode := createGroovyQuestionCode()
	question := m.Question{
		Codes: []m.QuestionCode{groovyQuestionCode},
	}
	db.Conn.Create(&question)

	Convey("POST run\n", t, func() {
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
		r, w, _ := makePostJSON("/api/v1/run", runBody)
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		fmt.Printf("%v", w)

		Convey("StatusCode = 201", func() {
			So(w.Code, ShouldEqual, 201)
		})
		Convey("Inserted Tryout", func() {
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
