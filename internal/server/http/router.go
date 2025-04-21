package http

import (
	"log"
	"net/http"
	"test-workmate/internal/server/handler"

	_ "test-workmate/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/gorilla/mux"
)

func InitRouter(taskHandler *handler.TaskHandler) *mux.Router {
	r := mux.NewRouter()
	r.Use(LoggingMiddleware)

	r.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		log.Printf("%s %s - %d", r.Method, r.RequestURI, lrw.statusCode)
	})
}
