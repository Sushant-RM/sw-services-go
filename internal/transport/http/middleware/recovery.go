package middleware

import (
	"log"
	"net/http"

	"github.com/BrajK111/sw-services/internal/transport/http/response"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {

				log.Printf("PANIC RECOVERED: %v", err)

				response.WriteError(
					w,
					http.StatusInternalServerError,
					"INTERNAL_ERROR",
					"internal server error",
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
