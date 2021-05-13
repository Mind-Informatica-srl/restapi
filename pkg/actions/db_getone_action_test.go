package actions

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mind-Informatica-srl/restapi/internal/testutils"
	"github.com/gorilla/mux"
)

func TestDBIGetOneAction(t *testing.T) {

	request := httptest.NewRequest("GET", "/getone/mario", nil)

	responseWriter := httptest.NewRecorder()

	vars := map[string]string{
		"id": "1",
	}

	request = mux.SetURLVars(request, vars)
	db, mock := testutils.SetupTestForGorm()
	delegate := testutils.SimpleObjectWithIdDelegate{DB: db}

	rows := sqlmock.NewRows([]string{"id", "nome", "cognome"}).
		AddRow("1", "mario", "rossi").AddRow("2", "paolo", "bianchi")
	query := `SELECT \* FROM "simple_object_with_ids" WHERE "simple_object_with_ids"."id" = \$1 ORDER BY "simple_object_with_ids"."id" LIMIT 1`
	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	action := DBGetOneAction{
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
	expectedBody := "{\"ID\":1,\"Nome\":\"mario\",\"Cognome\":\"rossi\"}\n"
	if body, err := io.ReadAll(responseWriter.Body); err != nil {
		t.Log(err)
		t.Fail()
	} else if strBody := string(body); strBody != expectedBody {
		t.Logf("Wrong response body: %s. Expected: %s", strBody, expectedBody)
		t.Fail()
	}
}
