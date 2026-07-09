package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	swerrors "github.com/egovernments/sw-services-go/internal/domain/errors"
)

// WriteError maps a service-layer error to the DIGIT standard error envelope.
// CustomException maps to 400 (matches Java's validation-failure contract for
// this service); anything else is an unexpected 500.
func WriteError(c *gin.Context, err error) {
	if custom, ok := err.(*swerrors.CustomException); ok {
		c.JSON(http.StatusBadRequest, custom.Response())
		return
	}
	c.JSON(http.StatusInternalServerError, swerrors.New("EG_SW_UNKNOWN_ERROR", err.Error()).Response())
}
