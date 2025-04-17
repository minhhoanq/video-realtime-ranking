package routes

import (
	"net/http"
)

type Routes struct {
	serverMux *http.ServeMux
}

func NewRouter(serverMux *http.ServeMux) *Routes {
	return &Routes{
		serverMux: serverMux,
	}
}

func methodHandlerFunc(method string, handler http.HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			handler(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (r *Routes) SetupRouter() http.Handler {
	r.serverMux.HandleFunc("/health-check", methodHandlerFunc(http.MethodGet, healthCheck))
	return r.serverMux
}
