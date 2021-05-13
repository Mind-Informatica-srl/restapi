package actions

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mind-Informatica-srl/restapi/internal/testutils"
	"github.com/gorilla/mux"
)

func TestDBUpdateAction(t *testing.T) {
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

	db, mock := testutils.SetupTestForGorm()
	delegate := testutils.SimpleObjectWithIdDelegate{DB: db}

	query := `UPDATE "simple_object_with_ids" SET "nome"=\$1,"cognome"=\$2 WHERE "id" = \$3`
	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs("John", "Doe", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	action := DBUpdateAction{
		Delegate: delegate,
	}

	action.Serve(responseWriter, request)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Log(err)
		t.Fail()
	}
	if responseWriter.Code != http.StatusOK {
		t.Logf("Wrong response code: %d", responseWriter.Code)
		t.Fail()
	}
	expectedBody := "{\"ID\":1,\"Nome\":\"John\",\"Cognome\":\"Doe\"}\n"
	if body, err := io.ReadAll(responseWriter.Body); err != nil {
		t.Log(err)
		t.Fail()
	} else if strBody := string(body); strBody != expectedBody {
		t.Logf("Wrong response body: %s. Expected: %s", strBody, expectedBody)
		t.Fail()
	}
}
