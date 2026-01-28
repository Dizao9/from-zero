package transport

import (
	"log"
	"net/http"
	"time"
)

func LoggindMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[GET] /users - %d ms", time.Since(now).Milliseconds())
	})
}
