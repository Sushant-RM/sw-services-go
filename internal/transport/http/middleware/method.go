package middleware

import (
	"net/http"

	"github.com/BrajK111/sw-services/internal/transport/http/response"
)

func RequireMethod(method string, next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != method {

			response.WriteError(
				w,
				http.StatusMethodNotAllowed,
				"METHOD_NOT_ALLOWED",
				"invalid HTTP method",
			)

			return
		}

		next.ServeHTTP(w, r)
	}
}
