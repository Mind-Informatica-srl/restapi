package actions

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Mind-Informatica-srl/restapi/pkg/testutils"
)

func TestDBGetAllAction(t *testing.T) {
	content, err := testutils.SomeSimpleObject()
	if err != nil {
		panic(err)
	}

	request := httptest.NewRequest("GET", "/getall", bytes.NewReader(content))
	responseWriter := httptest.NewRecorder()
	db, mock := testutils.SetupTestForGorm()
	delegate := testutils.SimpleObjectDelegate{DB: db}

	rows := sqlmock.NewRows([]string{"nome", "cognome"}).
		AddRow("mario", "rossi").AddRow("paolo", "bianchi")
	query := `SELECT \* FROM "simple_objects"`
	mock.ExpectQuery(query).WillReturnRows(rows)

	action := DBGetAllAction{
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
	expectedBody := "[{\"nome\":\"mario\",\"cognome\":\"rossi\"},{\"nome\":\"paolo\",\"cognome\":\"bianchi\"}]\n"
	if body, err := io.ReadAll(responseWriter.Body); err != nil {
		t.Log(err)
		t.Fail()
	} else if strBody := string(body); strBody != expectedBody {
		t.Logf("Wrong response body: %s. Expected: %s", strBody, expectedBody)
		t.Fail()
	}

}

func TestDBIGetAllActionWithQueryParams(t *testing.T) {
	content, err := testutils.SomeSimpleObject()
	if err != nil {
		panic(err)
	}

	request := httptest.NewRequest("GET", "/getall?q=nome.equal=mario", bytes.NewReader(content))
	responseWriter := httptest.NewRecorder()

	db, mock := testutils.SetupTestForGorm()
	delegate := testutils.SimpleObjectDelegate{DB: db}

	rows := sqlmock.NewRows([]string{"nome", "cognome"}).
		AddRow("mario", "rossi")
	query := `SELECT \* FROM "simple_objects" WHERE lower\(nome\) = \$1`
	mock.ExpectQuery(query).WithArgs("lower('mario')").WillReturnRows(rows)

	action := DBGetAllAction{
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
	expectedBody := "[{\"nome\":\"mario\",\"cognome\":\"rossi\"}]\n"
	if body, err := io.ReadAll(responseWriter.Body); err != nil {
		t.Log(err)
		t.Fail()
	} else if strBody := string(body); strBody != expectedBody {
		t.Logf("Wrong response body: %s. Expected: %s", strBody, expectedBody)
		t.Fail()
	}
}
