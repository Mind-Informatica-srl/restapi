package actions

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mind-Informatica-srl/restapi/internal/testutils"
)

func TestDBINsertAction(t *testing.T) {
	content, err := testutils.SomeSimpleObject()
	if err != nil {
		panic(err)
	}

	request := httptest.NewRequest("POST", "/insert", bytes.NewReader(content))
	responseWriter := httptest.NewRecorder()
	delegate := testutils.SimpleObjectDelegate
	mock := testutils.SetupTestForGorm(&delegate)

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "simple_objects" .*`).WithArgs("John", "Doe").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	action := DBInsertAction{
		Delegate: delegate,
	}

	if err := action.Serve(responseWriter, request); err != nil {
		t.Log(err)
		t.Fail()
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Log(err)
		t.Fail()
	}
	if responseWriter.Code != http.StatusCreated {
		t.Logf("Wrong response code: %d", responseWriter.Code)
	}
}
