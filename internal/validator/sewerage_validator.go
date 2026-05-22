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

	if connection.ConnectionType == "" {
		return errors.New("connectionType is required")
	}

	if connection.RoadType == "" {
		return errors.New("roadType is required")
	}

	if connection.ApplicationStatus == "" {
		return errors.New("applicationStatus is required")
	}

	return nil
}
