package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ascarter/requestid"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

// Auth 認証（dbはフェイク）
func Auth(db string) (fn func(http.Handler) http.Handler) {
	fn = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token != "admin" {
				respondError(w, http.StatusUnauthorized, fmt.Errorf("利用権限がありません"))
				return
			}
			h.ServeHTTP(w, r)
		})
	}
	return
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		t := t2.Sub(t1)
		reqID, ok := requestid.FromContext(r.Context())
		if !ok {
			reqID = uuid.New().String()
		}
		logger.Infof("request_id %s req_time %s req_time_nsec %v", reqID, t.String(), t.Nanoseconds())
	})
}

var devCORS = cors.New(cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
})

var stagingCORS = cors.New(cors.Options{
	AllowedOrigins:   []string{"staging.com"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
})

var productionCORS = cors.New(cors.Options{
	AllowedOrigins:   []string{"production.com"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
})

// CorsConfig CORSの設定を環境別に持っている
var CorsConfig = map[string]*cors.Cors{
	"develop":    devCORS,
	"staging":    stagingCORS,
	"production": productionCORS,
}
