package restapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mind-Informatica-srl/restapi/internal/testutils"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type SimpleObjectWithId struct {
	ID      int
	Nome    string
	Cognome string
}

func TestDBIGetOneAction(t *testing.T) {

	request := httptest.NewRequest("GET", "/getone/mario", nil)

	responseWriter := httptest.NewRecorder()

	vars := map[string]string{
		"id": "1",
	}

	request = mux.SetURLVars(request, vars)

	connectionPool, mock := testutils.SetupTestForGorm()

	rows := sqlmock.NewRows([]string{"id", "nome", "cognome"}).
		AddRow("1", "mario", "rossi").AddRow("2", "paolo", "bianchi")
	query := `SELECT \* FROM "simple_object_with_ids" WHERE "simple_object_with_ids"."id" = \$1 ORDER BY "simple_object_with_ids"."id" LIMIT 1`
	mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

	action := DBGetOneAction{
		ObjectCreator: func() interface{} {
			return &SimpleObjectWithId{}
		},
		DBProvider: func() *gorm.DB {
			return connectionPool
		},
		PKExtractor: func(r *http.Request) (interface{}, error) {
			vars := mux.Vars(r)
			if pk, err := strconv.Atoi(vars["id"]); err != nil {
				return nil, err
			} else {
				return pk, nil
			}
		},
	}

	action.ServeHTTP(responseWriter, request)
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
