package restapi

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mind-Informatica-srl/restapi/internal/testutils"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func TestDBUpdateAction(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/simple_object_with_id.json")
	if err != nil {
		panic(err)
	}
	request := httptest.NewRequest("PUT", "/getone/1", bytes.NewReader(content))

	responseWriter := httptest.NewRecorder()

	vars := map[string]string{
		"id": "1",
	}

	request = mux.SetURLVars(request, vars)

	connectionPool, mock := testutils.SetupTestForGorm()

	query := `UPDATE "simple_object_with_ids" SET "nome"=\$1,"cognome"=\$2 WHERE "id" = \$3`
	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs("John", "Doe", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	action := DBUpdateAction{
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
		PKVerificator: func(element interface{}, pk interface{}) *PKNotVerifiedError {
			e := element.(*SimpleObjectWithId)
			id := pk.(int)
			if e.ID != id {
				err := NewPKNotVerifiedError(element, pk)
				return &err
			}
			return nil
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
	expectedBody := "{\"ID\":1,\"Nome\":\"John\",\"Cognome\":\"Doe\"}\n"
	if body, err := io.ReadAll(responseWriter.Body); err != nil {
		t.Log(err)
		t.Fail()
	} else if strBody := string(body); strBody != expectedBody {
		t.Logf("Wrong response body: %s. Expected: %s", strBody, expectedBody)
		t.Fail()
	}
}
