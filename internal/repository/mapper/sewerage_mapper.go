package mapper

import (
	"encoding/json"

	"github.com/BrajK111/sw-services/internal/domain/dto"
)

func ParseConnectionHolders(data []byte) ([]dto.ConnectionHolder, error) {

	if len(data) == 0 {
		return []dto.ConnectionHolder{}, nil
	}

	var holders []dto.ConnectionHolder

	err := json.Unmarshal(data, &holders)
	if err != nil {
		return nil, err
	}

	return holders, nil
}

func ParsePlumberInfo(data []byte) ([]dto.PlumberInfo, error) {

	if len(data) == 0 {
		return []dto.PlumberInfo{}, nil
	}

	var plumbers []dto.PlumberInfo

	err := json.Unmarshal(data, &plumbers)
	if err != nil {
		return nil, err
	}

	return plumbers, nil
}

func ParseRoadCuttingInfo(data []byte) ([]dto.RoadCuttingInfo, error) {

	if len(data) == 0 {
		return []dto.RoadCuttingInfo{}, nil
	}

	var roadInfo []dto.RoadCuttingInfo

	err := json.Unmarshal(data, &roadInfo)
	if err != nil {
		return nil, err
	}

	return roadInfo, nil
}
