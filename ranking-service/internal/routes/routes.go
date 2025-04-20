package routes

import (
	"net/http"
	"video-realtime-ranking/ranking-service/internal/handler/resful"
)

type Routes struct {
	serverMux      *http.ServeMux
	rankingHandler *resful.Handler
}

func NewRouter(serverMux *http.ServeMux,
	rankingHandler *resful.Handler,
) *Routes {
	return &Routes{
		serverMux:      serverMux,
		rankingHandler: rankingHandler,
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
	r.serverMux.HandleFunc("/top-k", methodHandlerFunc(http.MethodGet, r.rankingHandler.GetTopK))
	return r.serverMux
}
