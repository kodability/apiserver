package tests

import (
	"encoding/json"
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

func TestJUnitReport(t *testing.T) {
	Convey("go", t, func() {
		Convey("test_go_ok", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_go_ok.xml", true)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       1,
				ElapsedTime: 0.003,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "TestSum", Time: 0.002, Error: ""},
				},
			})
		})
		Convey("test_go_failure", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_go_failure.xml", true)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       1,
				Failures:    1,
				ElapsedTime: 0,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "[build failed]", Time: 0, Error: "Failed"},
				},
			})
		})
	})

	Convey("groovy", t, func() {
		Convey("test_groovy_ok", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_groovy_ok.xml", false)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       2,
				ElapsedTime: 0.087,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test2", Time: 0.002, Error: ""},
					run.JUnitTestcaseResult{Name: "test1", Time: 0.085, Error: ""},
				},
			})
		})
		Convey("test_groovy_error", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_groovy_error.xml", false)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       2,
				Errors:      1,
				ElapsedTime: 0.109,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test2", Time: 0, Error: "expected:<2> but was:<1>"},
					run.JUnitTestcaseResult{Name: "test1", Time: 0.109, Error: ""},
				},
			})
		})
		Convey("test_groovy_failure", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_groovy_failure.xml", false)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       2,
				Failures:    2,
				ElapsedTime: 0.072,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test1", Time: 0.072,
						Error: "Cannot cast object 'foo' with class 'java.lang.String' to class 'int'",
					},
					run.JUnitTestcaseResult{Name: "test2", Time: 0,
						Error: "Cannot cast object 'foo' with class 'java.lang.String' to class 'int'",
					},
				},
			})
		})
	})

	Convey("java", t, func() {
		Convey("test_java_ok", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_java_ok.xml", false)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       2,
				ElapsedTime: 0.004,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test1", Time: 0.003, Error: ""},
					run.JUnitTestcaseResult{Name: "test2", Time: 0.001, Error: ""},
				},
			})
		})
		Convey("test_java_error", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_java_error.xml", false)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       2,
				Errors:      2,
				ElapsedTime: 0.004,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test2", Time: 0.001, Error: "/ by zero"},
					run.JUnitTestcaseResult{Name: "test1", Time: 0.003, Error: "/ by zero"},
				},
			})
		})
		Convey("test_java_failure", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_java_failure.xml", false)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       2,
				Failures:    2,
				ElapsedTime: 0.005,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test2", Time: 0.001, Error: "expected:<1> but was:<0>"},
					run.JUnitTestcaseResult{Name: "test1", Time: 0.004, Error: "expected:<55> but was:<5>"},
				},
			})
		})
	})

	Convey("javascript", t, func() {
		Convey("test_javascript_ok", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_javascript_ok.xml", true)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       2,
				ElapsedTime: 0.003,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test1", Time: 0.001, Error: ""},
					run.JUnitTestcaseResult{Name: "test2", Time: 0.002, Error: ""},
				},
			})
		})
		Convey("test_javascript_failure", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_javascript_failure.xml", true)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       2,
				Failures:    2,
				ElapsedTime: 0.003,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test1", Time: 0.002,
						Error: "AssertionError [ERR_ASSERTION]: Infinity == 55\n    at Context.<anonymous> (test/exampleSpec.js:7:16)",
					},
					run.JUnitTestcaseResult{Name: "test2", Time: 0.001,
						Error: "AssertionError [ERR_ASSERTION]: Infinity == 1\n    at Context.<anonymous> (test/exampleSpec.js:10:16)",
					},
				},
			})
		})
	})

	Convey("python", t, func() {
		Convey("test_python_ok", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_python_ok.xml", false)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       1,
				ElapsedTime: 0.001,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test_sum", Time: 0.001, Error: ""},
				},
			})
		})
		Convey("test_python_error", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_python_error.xml", false)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       1,
				Errors:      1,
				ElapsedTime: 0.005,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test_sum", Time: 0.005, Error: "division by zero"},
				},
			})
		})
		Convey("test_python_failure", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_python_failure.xml", false)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       1,
				Failures:    1,
				ElapsedTime: 0.003,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test_sum", Time: 0.003, Error: "551 != 55.0"},
				},
			})
		})
	})

	Convey("scala", t, func() {
		Convey("test_scala_ok", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_scala_ok.xml", false)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       1,
				ElapsedTime: 0.009,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test1", Time: 0.009, Error: ""},
				},
			})
		})
		Convey("test_scala_error", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_scala_error.xml", false)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       1,
				Errors:      1,
				ElapsedTime: 0.008,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test1", Time: 0.008, Error: "/ by zero"},
				},
			})
		})
		Convey("test_scala_failure", func() {
			report, err := run.ReadJunitReportFile("./tests/resources/test_scala_failure.xml", false)
			So(err, ShouldBeNil)
			So(report, ShouldResemble, &run.JUnitReport{
				Tests:       1,
				Failures:    1,
				ElapsedTime: 0.014,
				TestResults: []run.JUnitTestcaseResult{
					run.JUnitTestcaseResult{Name: "test1", Time: 0.014, Error: "expected:<12> but was:<1>"},
				},
			})
		})
	})
}
