package restapi

import (
	"bytes"
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

func TestDBDeleteAction(t *testing.T) {
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

	query := `DELETE FROM "simple_object_with_ids" WHERE "simple_object_with_ids"."id" = \$1`
	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	action := DBDeleteAction{
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
		PKAssigner: func(element interface{}, pk interface{}) {
			e := element.(*SimpleObjectWithId)
			id := pk.(int)
			e.ID = id
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
}
