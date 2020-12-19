package apiserver

import "net/http"

type responseWriter struct {
	//анонимное поле
	http.ResponseWriter
	code int
}

//WriteHeader write code to response header.
func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
