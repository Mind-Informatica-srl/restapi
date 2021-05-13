package actions

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mind-Informatica-srl/restapi/internal/testutils"
	"github.com/gorilla/mux"
)

func TestDBDeleteAction(t *testing.T) {
	content, err := testutils.SomeSimpleObjectWithId()
	if err != nil {
		panic(err)
	}
	request := httptest.NewRequest("PUT", "/getone/1", bytes.NewReader(content))

	responseWriter := httptest.NewRecorder()

	vars := map[string]string{
		"id": "1",
	}

	request = mux.SetURLVars(request, vars)

	delegate := testutils.SimpleObjectWithIdDelegate
	mock := testutils.SetupTestForGorm(&delegate)

	query := `DELETE FROM "simple_object_with_ids" WHERE "simple_object_with_ids"."id" = \$1`
	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	action := DBDeleteAction{
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
	if responseWriter.Code != http.StatusOK {
		t.Logf("Wrong response code: %d", responseWriter.Code)
		t.Fail()
	}
}
