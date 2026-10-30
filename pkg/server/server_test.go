package server

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mind-Informatica-srl/restapi/pkg/actions"
)

func newTestAction(f actions.ActionFunc) *actions.Action {
	return &actions.Action{ActionFunc: f, Path: "/test", Method: http.MethodGet}
}

// serveTest runs the action through actionHandler and reports whether afterRequest was called
func serveTest(t *testing.T, f actions.ActionFunc) (*httptest.ResponseRecorder, bool) {
	t.Helper()
	afterCalled := false
	handler := actionHandler(newTestAction(f), defaultErrorTranslator, func(req *http.Request) error {
		afterCalled = true
		return nil
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/test", nil))
	return rec, afterCalled
}

func assertResponse(t *testing.T, rec *httptest.ResponseRecorder, status int, body string) {
	t.Helper()
	if rec.Code != status {
		t.Errorf("status = %d, want %d", rec.Code, status)
	}
	if got := rec.Body.String(); got != body {
		t.Errorf("body = %q, want %q", got, body)
	}
}

func assertTextPlain(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("Content-Type = %q, want text/plain", ct)
	}
}

func TestActionHandlerNilErr(t *testing.T) {
	rec, after := serveTest(t, func(w http.ResponseWriter, r *http.Request) *actions.ActionError {
		return &actions.ActionError{Err: nil, Status: http.StatusBadRequest}
	})
	assertResponse(t, rec, http.StatusBadRequest, "Bad Request\n")
	assertTextPlain(t, rec)
	if !after {
		t.Error("afterRequest not called")
	}
}

func TestActionHandlerZeroStatus(t *testing.T) {
	rec, _ := serveTest(t, func(w http.ResponseWriter, r *http.Request) *actions.ActionError {
		return &actions.ActionError{Err: errors.New("x"), Status: 0}
	})
	assertResponse(t, rec, http.StatusInternalServerError, "x\n")
	assertTextPlain(t, rec)
}

func TestActionHandlerNestedActionError(t *testing.T) {
	rec, _ := serveTest(t, func(w http.ResponseWriter, r *http.Request) *actions.ActionError {
		return &actions.ActionError{
			Err: &actions.ActionError{
				Err:    &actions.ActionError{Err: errors.New("inner"), Status: http.StatusConflict},
				Status: http.StatusBadRequest,
			},
			Status: http.StatusInternalServerError,
		}
	})
	assertResponse(t, rec, http.StatusConflict, "inner\n")
}

func TestActionHandlerNestedNilActionError(t *testing.T) {
	rec, _ := serveTest(t, func(w http.ResponseWriter, r *http.Request) *actions.ActionError {
		var inner *actions.ActionError
		return &actions.ActionError{Err: inner, Status: http.StatusNotFound}
	})
	assertResponse(t, rec, http.StatusNotFound, "Not Found\n")
}

func TestActionHandlerPanic(t *testing.T) {
	rec, after := serveTest(t, func(w http.ResponseWriter, r *http.Request) *actions.ActionError {
		panic("boom")
	})
	assertResponse(t, rec, http.StatusInternalServerError, "Internal Server Error\n")
	assertTextPlain(t, rec)
	if !after {
		t.Error("afterRequest not called after panic")
	}
}

func TestActionHandlerPanicAfterWrite(t *testing.T) {
	afterCalled := false
	handler := actionHandler(newTestAction(func(w http.ResponseWriter, r *http.Request) *actions.ActionError {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		panic("boom")
	}), defaultErrorTranslator, func(req *http.Request) error {
		afterCalled = true
		return nil
	})
	var serverLog bytes.Buffer
	srv := httptest.NewUnstartedServer(handler)
	srv.Config.ErrorLog = log.New(&serverLog, "", 0)
	srv.Start()

	resp, err := http.Get(srv.URL)
	if err != nil {
		srv.Close()
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	// Close waits for the handler to return, so serverLog and afterCalled are safe to read
	srv.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if string(body) != "ok" {
		t.Errorf("body = %q, want %q", body, "ok")
	}
	if strings.Contains(serverLog.String(), "superfluous") {
		t.Errorf("unexpected server log: %s", serverLog.String())
	}
	if !afterCalled {
		t.Error("afterRequest not called after panic")
	}
}

func TestActionHandlerAbortHandler(t *testing.T) {
	defer func() {
		if v := recover(); v != http.ErrAbortHandler {
			t.Errorf("recovered %v, want http.ErrAbortHandler", v)
		}
	}()
	serveTest(t, func(w http.ResponseWriter, r *http.Request) *actions.ActionError {
		panic(http.ErrAbortHandler)
	})
	t.Error("panic not propagated")
}

func TestActionHandlerResponseController(t *testing.T) {
	rec, _ := serveTest(t, func(w http.ResponseWriter, r *http.Request) *actions.ActionError {
		w.Write([]byte("chunk"))
		if err := http.NewResponseController(w).Flush(); err != nil {
			t.Errorf("ResponseController.Flush: %v", err)
		}
		w.(http.Flusher).Flush()
		return nil
	})
	assertResponse(t, rec, http.StatusOK, "chunk")
	if !rec.Flushed {
		t.Error("response not flushed")
	}
}
