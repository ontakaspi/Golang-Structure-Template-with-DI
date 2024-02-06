package _testTools

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ztrue/tracerr"
	"go.mongodb.org/mongo-driver/mongo"
	"golang-structure-template-with-di/app/helper"
	"golang-structure-template-with-di/libraries/httpResponse"
	"golang-structure-template-with-di/libraries/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	logger1 "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type errorHandler struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
}
type Expected struct {
	TestName           string
	ExpectedStatusCode int
	ExpectedStatusBody string
	RequestBody        interface{}
	Error              bool
	Params             []gin.Param
	Query              url.Values
	AuthToken          string
	BeforeExecuteFunc  func()
	AfterExecuteFunc   func()
	FunctionController func(c *gin.Context)
}

var DBMock *gorm.DB
var MongoDBMock *mongo.Client
var MongoCollectionMock *mongo.Collection

func StructToIKeyAndValue(dataStruct interface{}) ([]string, []driver.Value) {
	var dataMap []driver.Value
	var dataMapKey []string
	var myMap map[string]interface{}
	data, _ := json.Marshal(dataStruct)
	json.Unmarshal(data, &myMap)
	for key, value := range myMap {
		if (key == "CreatedAt" || key == "UpdatedAt" || key == "DeletedAt") && value != nil {
			//parse string to time
			value, _ := time.Parse(time.RFC3339, value.(string))
			dataMap = append(dataMap, value)
			dataMapKey = append(dataMapKey, helper.ToSnakeCase(key))
		} else {
			dataMap = append(dataMap, value)
			dataMapKey = append(dataMapKey, helper.ToSnakeCase(key))
		}

	}
	return dataMapKey, dataMap
}

func KeyAndValueToString(dataMapKey []string, dataMap []driver.Value) (string, string) {

	var keyString string
	var valueString string
	for i, key := range dataMapKey {
		if i == 0 {
			keyString = key
			valueString = "?"
		} else {
			keyString = "\"" + keyString + "," + key
			valueString = valueString + ",?"
		}
	}
	return keyString, valueString

}

func SetupDatabase() sqlmock.Sqlmock {

	var (
		db  *sql.DB
		err error
	)
	db, Mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
		WithoutReturning:     true,
	})
	DBMock, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger1.Default.LogMode(logger1.Silent),
	})
	if err != nil {
		panic(err)
	}

	return Mock
}

//func SetupDatabaseMongo(t *testing.T) *mtest.T {
//
//	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
//	defer mt.Close()
//	MongoCollectionMock = mt.Coll
//	MongoDBMock = mt.Client
//	return mt
//}

func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	os.Setenv("UNIT_TEST", "1")
	gin.SetMode(gin.ReleaseMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

func NewHtppRecorder(method string, params gin.Params, u url.Values, token string, content interface{}) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c := GetTestGinContext(w)
	c.Request.Method = method
	c.Request.Header.Set("Content-Type", "application/json")
	if token != "" {
		c.Request.Header.Set("Authorization", "Bearer "+token)
	}
	// set path params
	c.Params = params
	// set query params
	c.Request.URL.RawQuery = u.Encode()
	// set request body
	if content != nil {
		jsonbytes, err := json.Marshal(content)
		if err != nil {
			panic(err)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonbytes))
	}

	return w, c
}

func CatchPanic(t *testing.T, w *httptest.ResponseRecorder, ScenarioExpected Expected, err interface{}) {
	unitTest, _ := strconv.Atoi(os.Getenv("UNIT_TEST"))
	var errData interface{}
	var httpCode int

	switch err.(type) {
	case map[string]interface{}:
		recoverErr := err.(map[string]interface{})
		httpCode = recoverErr["httpCode"].(int)
		errData = recoverErr["error"]
		var errorData error
		if recoverErr["errorData"] != nil {
			errorData = recoverErr["errorData"].(error)
		} else {
			errorData = nil
		}

		if httpCode >= 500 && unitTest != 1 {
			logger.SetLogFileAndConsole(logger.LogData{
				Message: "Unexpected Error",
				CustomFields: logrus.Fields{
					"data": errData,
				},
				Level: "ERROR",
			})
			if errorData != nil {
				dataErrr := tracerr.Wrap(errorData)
				tracerr.PrintSourceColor(dataErrr)
			}

		}
	case string:
		httpCode = http.StatusInternalServerError
		errData = errorHandler{
			Status:  "error",
			Message: err,
		}
		if httpCode >= 500 && unitTest != 1 {
			logger.SetLogFileAndConsole(logger.LogData{
				Message: "Unexpected Error",
				CustomFields: logrus.Fields{
					"data": errData,
				},
				Level: "ERROR",
			})
			dataErrr := tracerr.New(err.(string))
			tracerr.PrintSourceColor(dataErrr)
		}
	default:
		httpCode = http.StatusInternalServerError
		errData = errorHandler{
			Status:  "error",
			Message: err.(error).Error(),
		}
		if httpCode >= 500 && unitTest != 1 {
			logger.SetLogFileAndConsole(logger.LogData{
				Message: "Unexpected Error",
				CustomFields: logrus.Fields{
					"data": errData,
				},
				Level: "ERROR",
			})

			dataErrr := tracerr.Wrap(err.(error))
			tracerr.PrintSourceColor(dataErrr)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write([]byte(fmt.Sprintf("%v", errData)))

	assert.EqualValues(t, ScenarioExpected.ExpectedStatusCode, w.Code)
	assert.EqualValues(t, ScenarioExpected.ExpectedStatusBody, w.Body.String())

}

func CatchPanicUniversal(err interface{}) (message string) {
	switch err.(type) {
	case map[string]interface{}:
		var message1 string
		recoverErr := err.(map[string]interface{})
		errData := recoverErr["error"]
		switch errData.(type) {
		case string:
			message1 = errData.(string)
		case error:
			message1 = errData.(error).Error()
		case []string:
			message1 = strings.Join(errData.([]string), ",")
		case httpResponse.ErrorResp:
			message1 = errData.(httpResponse.ErrorResp).Message
		case httpResponse.ValidationErrorResp:
			erros := errData.(httpResponse.ValidationErrorResp).Errors
			var errors []string
			for _, v := range erros {
				errors = append(errors, v)
			}
			message1 = strings.Join(errors, ",")
		}
		return message1
	case string:
		return err.(string)
	default:
		return err.(error).Error()
	}

}

func RunControllerTest(t *testing.T, testName string, method string,
	expectedScenario Expected) bool {
	return t.Run(testName, func(t *testing.T) {

		//run the func
		if expectedScenario.BeforeExecuteFunc != nil {
			expectedScenario.BeforeExecuteFunc()
		}
		var w *httptest.ResponseRecorder
		var c *gin.Context
		var AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

		// check if method is allowed
		if !utils.Contains(AllowedMethods, method) {
			t.Errorf("Method not allowed")
		}
		w, c = NewHtppRecorder(method, expectedScenario.Params, expectedScenario.Query,
			expectedScenario.AuthToken, expectedScenario.RequestBody)

		if expectedScenario.Error {
			defer func() {
				if err := recover(); err != nil {
					CatchPanic(t, w, expectedScenario, err)
					if expectedScenario.AfterExecuteFunc != nil {
						expectedScenario.AfterExecuteFunc()
					}
				}
			}()
			if expectedScenario.FunctionController != nil {
				expectedScenario.FunctionController(c)
			}
		} else {
			defer func() {
				if err := recover(); err != nil {
					t.Errorf("Error: %v", err)
				}
			}()
			if expectedScenario.FunctionController != nil {
				expectedScenario.FunctionController(c)
			}
			assert.Equal(t, expectedScenario.ExpectedStatusCode, w.Code)
			assert.Equal(t, expectedScenario.ExpectedStatusBody, w.Body.String())
			if expectedScenario.AfterExecuteFunc != nil {
				expectedScenario.AfterExecuteFunc()
			}
		}
	})
}
