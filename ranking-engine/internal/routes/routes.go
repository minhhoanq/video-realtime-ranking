package routes

import (
	"net/http"
	"video-realtime-ranking/ranking-engine/internal/service"
)

type Routes struct {
	serverMux *http.ServeMux
	service   service.InteractionService
}

func NewRouter(serverMux *http.ServeMux,
	service service.InteractionService) *Routes {
	return &Routes{
		serverMux: serverMux,
		service:   service,
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
	r.serverMux.HandleFunc("/interaction", methodHandlerFunc(http.MethodPost, r.service.CreateInteraction))
	return r.serverMux
}
