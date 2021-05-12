package restapi

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mind-Informatica-srl/restapi/internal/testutils"
	"gorm.io/gorm"
)

type SimpleObject struct {
	Nome    string `json:"nome" gorm:"primaryKey"`
	Cognome string `json:"cognome"`
}

func TestDBINsertAction(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/simple_object.json")
	if err != nil {
		panic(err)
	}

	request := httptest.NewRequest("POST", "/insert", bytes.NewReader(content))
	responseWriter := httptest.NewRecorder()

	connectionPool, mock := testutils.SetupTestForGorm()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "simple_objects" .*`).WithArgs("John", "Doe").WillReturnResult(sqlmock.NewResult(1, 1))

	action := DBInsertAction{
		ObjectCreator: func() interface{} {
			return &SimpleObject{}
		},
		DBProvider: func() *gorm.DB {
			return connectionPool
		},
	}

	action.ServeHTTP(responseWriter, request)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Log(err)
		t.Fail()
	}
	if responseWriter.Code != http.StatusCreated {
		t.Logf("Wrong response code: %d", responseWriter.Code)
	}
}
