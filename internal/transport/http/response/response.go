package response

import (
	"encoding/json"
	"net/http"

	customerror "github.com/BrajK111/sw-services/internal/domain/errors"
)

func WriteJSON(
	w http.ResponseWriter,
	status int,
	data interface{},
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {

	errorResponse := customerror.ErrorResponse{
		Code:    code,
		Message: message,
	}

	WriteJSON(
		w,
		status,
		errorResponse,
	)
}
