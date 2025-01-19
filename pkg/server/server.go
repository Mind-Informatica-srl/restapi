package server

import (
	"net/http"

	"github.com/Mind-Informatica-srl/restapi/pkg/actions"
	"github.com/Mind-Informatica-srl/restapi/pkg/controllers"
	"github.com/Mind-Informatica-srl/restapi/pkg/logger"
	"github.com/gorilla/mux"
)

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

// Serve start the server
func (s RestApiServer) Serve(listenAddresses string) error {
	logger.Log().Info("Starting web server", "listenAddresses", listenAddresses)
	return http.ListenAndServe(listenAddresses, &s)
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
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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
		if err := action.Serve(w, r); err != nil {
			logger.Log().Error(err, "server error", "data error", err.Data)
			http.Error(w, errTranslator(err.Error()), err.Status)
		}
		if afterRequest != nil {
			if err := afterRequest(r); err != nil {
				http.Error(w, errTranslator(err.Error()), http.StatusInternalServerError)
			}
		}
	})
}

func defaultErrorTranslator(value string) string {
	return value
}
