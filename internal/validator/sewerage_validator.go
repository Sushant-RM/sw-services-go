package validator

import (
	"errors"

	"github.com/BrajK111/sw-services/internal/domain/dto"
)

func ValidateCreateConnection(
	request dto.CreateSewerageConnectionRequest,
) error {

	connection := request.SewerageConnection

	if connection.TenantID == "" {
		return errors.New("tenantId is required")
	}

	if connection.PropertyID == "" {
		return errors.New("propertyId is required")
	}

	return nil
}
