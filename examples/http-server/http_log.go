package main

import (
	"log"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

// LogHTTP is a simple logging middleware
func LogHTTP(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := statusWriter{ResponseWriter: w}
		handler.ServeHTTP(&sw, r)
		duration := time.Now().Sub(start)

		log.Printf("%s - \"%s %s %s\" \"%s\" \"%s\" %d %s %dB",
			r.RemoteAddr,
			r.Method,
			r.RequestURI,
			r.Proto,
			r.Referer(),
			r.UserAgent(),
			sw.status,
			duration,
			sw.length,
		)
	}
}
