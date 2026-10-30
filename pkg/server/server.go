package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/Mind-Informatica-srl/restapi/pkg/actions"
	"github.com/Mind-Informatica-srl/restapi/pkg/controllers"
	"github.com/Mind-Informatica-srl/restapi/pkg/logger"
	"github.com/gorilla/mux"
)

// shutdownTimeout è il tempo massimo concesso alle richieste in corso per
// completarsi dopo la ricezione di SIGTERM/SIGINT, prima di forzare la chiusura.
// Deve restare sotto il terminationGracePeriodSeconds del pod Kubernetes.
const shutdownTimeout = 25 * time.Second

// RestApiServer represent a server and provid usefull functions to configure it
type RestApiServer struct {
	mux.Router
	BasePath           string
	jwtHandler         func(next http.Handler) http.Handler
	authHandler        func(next http.Handler, authorizations []string) http.Handler
	authUserContextKey interface{}
	errorTranslator    func(string) string
	beforeRequest      func(next http.Handler) http.Handler
	afterRequest       func(req *http.Request) error
}

// NewRestApiServer instantiate a new RestApiServer
//
// router *mux.Router,
//
// basePath string,
//
// jwtHandler func(next http.Handler) http.Handler,
//
// authHandler func(next http.Handler, authorizations []string) http.Handler,
//
// authUserContextKey interface{},
//
// errorTranslator func(string) string: funzione per tradurre errori per il client
func NewRestApiServer(
	router *mux.Router,
	basePath string,
	jwtHandler func(next http.Handler) http.Handler,
	authHandler func(next http.Handler, authorizations []string) http.Handler,
	authUserContextKey interface{},
	errorTranslator func(string) string,
	beforeRequest func(next http.Handler) http.Handler,
	afterRequest func(req *http.Request) error,
) RestApiServer {
	router.Use(requestLoggingMiddleware(), requestCorsMiddleware())
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusOK)
	})
	if errorTranslator == nil {
		errorTranslator = defaultErrorTranslator
	}
	return RestApiServer{
		Router:             *router,
		BasePath:           basePath,
		jwtHandler:         jwtHandler,
		authHandler:        authHandler,
		authUserContextKey: authUserContextKey,
		errorTranslator:    errorTranslator,
		afterRequest:       afterRequest,
		beforeRequest:      beforeRequest,
	}
}

// RegisterAction register an handler to the relative path on the server
func (s *RestApiServer) RegisterAction(basePath string, action actions.AbstractAction) *mux.Route {
	var handler http.Handler = actionHandler(action, s.errorTranslator, s.afterRequest)
	if !action.IsSkipAuth() {
		handler = s.jwtHandler(s.authHandler(handler, action.GetAuthorizations()))
	}
	if s.beforeRequest != nil {
		handler = s.beforeRequest(handler)
	}
	return s.Handle(s.BasePath+basePath+action.GetPath(), handler).Methods(action.GetMethod())
}

// RegisterController register all the actions in the controller on the server
func (s *RestApiServer) RegisterController(controller *controllers.Controller) []*mux.Route {
	var routes []*mux.Route
	for _, a := range controller.Actions {
		routes = append(routes, s.RegisterAction(controller.Path, a))
	}
	return routes
}

// Serve start the server and blocks until it stops.
//
// On SIGTERM/SIGINT (the signal sent by Kubernetes when terminating a pod)
// it stops accepting new connections and waits up to shutdownTimeout for
// in-flight requests to complete before returning, so the process can exit
// cleanly (exit code 0) instead of being killed mid-request.
func (s RestApiServer) Serve(listenAddresses string) error {
	logger.Log().Info("Starting web server", "listenAddresses", listenAddresses)

	httpServer := &http.Server{
		Addr:    listenAddresses,
		Handler: &s,
	}

	serveErr := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serveErr <- err
		}
		close(serveErr)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sigCh)

	select {
	case err := <-serveErr:
		return err
	case sig := <-sigCh:
		logger.Log().Info("shutdown signal received, draining connections", "signal", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			return err
		}
		return <-serveErr
	}
}

// RequestLoggingMiddleware is a middleware logging each request arrive to web server
func requestLoggingMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req)
		})
	}
}

// RequestCorsMiddleware is a middleware enabling cors each request arrive to web server
func requestCorsMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Access-Control-Allow-Origin, X-Client-Version")
			// Check if the request is a preflight request
			if req.Method == http.MethodOptions && req.Header.Get("Access-Control-Request-Method") != "" {
				// Handle preflight request
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

func actionHandler(
	action actions.AbstractAction,
	errTranslator func(string) string,
	afterRequest func(req *http.Request) error,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tw := &trackingResponseWriter{ResponseWriter: w}
		// The action runs in a closure so that a panic is recovered here and
		// afterRequest is executed anyway. Without the recover net/http would
		// catch the panic and close the connection with no response at all,
		// so the only observable change is that the client gets a 500
		// instead of a network error
		func() {
			defer func() {
				if v := recover(); v != nil {
					if v == http.ErrAbortHandler {
						// intentional abort: let net/http handle it
						panic(v)
					}
					log := logger.Log()
					if log == nil {
						log = slog.Default()
					}
					log.Error("panic in action handler",
						"panic", fmt.Sprint(v),
						"method", r.Method,
						"path", r.URL.Path,
						"stack", string(debug.Stack()),
					)
					if !tw.wroteHeader {
						http.Error(tw, errTranslator(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
					}
				}
			}()
			if err := action.Serve(tw, r); err != nil {
				// use the innermost ActionError, stopping at a nil one
				for {
					if actionErr, ok := err.Err.(*actions.ActionError); ok && actionErr != nil {
						err = actionErr
					} else {
						break
					}
				}
				status := err.Status
				if status < 100 || status > 999 {
					status = http.StatusInternalServerError
				}
				http.Error(tw, errTranslator(err.Error()), status)
			}
		}()
		if afterRequest != nil {
			if err := afterRequest(r); err != nil {
				http.Error(w, errTranslator(err.Error()), http.StatusInternalServerError)
			}
		}
	})
}

// trackingResponseWriter records whether the response headers have already
// been sent, so that after a panic an error is written only if it is still
// possible
type trackingResponseWriter struct {
	http.ResponseWriter
	wroteHeader bool
}

func (w *trackingResponseWriter) WriteHeader(code int) {
	// informational headers (1xx except 101) don't commit the response
	if code < 100 || code > 199 || code == http.StatusSwitchingProtocols {
		w.wroteHeader = true
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *trackingResponseWriter) Write(b []byte) (int, error) {
	w.wroteHeader = true
	return w.ResponseWriter.Write(b)
}

// ReadFrom keeps the io.ReaderFrom optimization (sendfile) of the underlying writer
func (w *trackingResponseWriter) ReadFrom(src io.Reader) (int64, error) {
	w.wroteHeader = true
	if rf, ok := w.ResponseWriter.(io.ReaderFrom); ok {
		return rf.ReadFrom(src)
	}
	return io.Copy(struct{ io.Writer }{w.ResponseWriter}, src)
}

func (w *trackingResponseWriter) Flush() {
	w.wroteHeader = true
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack is kept for callers that type-assert http.Hijacker instead of using http.ResponseController
func (w *trackingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	conn, rw, err := http.NewResponseController(w.ResponseWriter).Hijack()
	if err == nil {
		w.wroteHeader = true
	}
	return conn, rw, err
}

// Unwrap allows http.ResponseController to reach the underlying writer
func (w *trackingResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func defaultErrorTranslator(value string) string {
	return value
}
