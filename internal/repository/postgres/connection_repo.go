package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/BrajK111/sw-services/internal/domain/dto"
	"github.com/BrajK111/sw-services/internal/repository/mapper"
)

type SewerageRepository interface {
	CreateConnection(
		connection dto.SewerageConnection,
		applicationNumber string,
		id string,
	) error

	SearchConnection(
		tenantID string,
		applicationNumber string,
	) ([]dto.SewerageConnection, error)

	UpdateConnection(
		connection dto.SewerageConnection,
	) error
}

type sewerageRepository struct {
	db *sql.DB
}

func NewSewerageRepository(
	db *sql.DB,
) SewerageRepository {
	return &sewerageRepository{
		db: db,
	}
}

func (r *sewerageRepository) CreateConnection(
	connection dto.SewerageConnection,
	applicationNumber string,
	id string,
) error {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	connectionHoldersJSON, _ :=
		json.Marshal(connection.ConnectionHolders)

	plumberInfoJSON, _ :=
		json.Marshal(connection.PlumberInfo)

	roadCuttingInfoJSON, _ :=
		json.Marshal(connection.RoadCuttingInfo)

	query := `
INSERT INTO eg_sw_connection (
	id,
	tenantid,
	applicationno,
	property_id,
	connectiontype,
	roadtype,
	roadcuttingarea,
	applicationstatus,
	status,
	channel,
	connection_holders,
	plumber_info,
	road_cutting_info,
	createdby,
	lastmodifiedby,
	createdtime,
	lastmodifiedtime
)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9,
	$10,
	$11,
	$12,
	$13,
	$14,
	$15,
	$16,
	$17
)
`

	createdBy := "system"
	if connection.AuditDetails != nil {
		createdBy = connection.AuditDetails.CreatedBy
	}

	currentTime := time.Now().UnixMilli()

	_, err := r.db.ExecContext(
		ctx,
		query,
		id,
		connection.TenantID,
		applicationNumber,
		connection.PropertyID,
		connection.ConnectionType,
		connection.RoadType,
		connection.RoadCuttingArea,
		connection.ApplicationStatus,
		connection.ApplicationStatus, // status
		connection.Channel,
		connectionHoldersJSON,
		plumberInfoJSON,
		roadCuttingInfoJSON,
		createdBy,
		createdBy,
		currentTime,
		currentTime,
	)

	return err
}

func (r *sewerageRepository) SearchConnection(
	tenantID string,
	applicationNumber string,
) ([]dto.SewerageConnection, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	query := `
SELECT
	tenantid,
	property_id,
	applicationno,
	connectiontype,
	roadtype,
	roadcuttingarea,
	applicationstatus,
	channel,
	connection_holders,
	plumber_info,
	road_cutting_info,
	connectionno
FROM eg_sw_connection
WHERE 1=1
`

	var args []interface{}
	argPosition := 1

	if tenantID != "" {
		query += "\nAND tenantid = $" + fmt.Sprint(argPosition)
		args = append(args, tenantID)
		argPosition++
	}

	if applicationNumber != "" {
		query += "\nAND applicationno = $" + fmt.Sprint(argPosition)
		args = append(args, applicationNumber)
		argPosition++
	}

	rows, err := r.db.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []dto.SewerageConnection

	for rows.Next() {

		var connection dto.SewerageConnection

		var connectionHoldersData []byte
		var plumberInfoData []byte
		var roadCuttingInfoData []byte
		var connectionNo sql.NullString

		err := rows.Scan(
			&connection.TenantID,
			&connection.PropertyID,
			&connection.ApplicationNumber,
			&connection.ConnectionType,
			&connection.RoadType,
			&connection.RoadCuttingArea,
			&connection.ApplicationStatus,
			&connection.Channel,
			&connectionHoldersData,
			&plumberInfoData,
			&roadCuttingInfoData,
			&connectionNo,
		)
		if err != nil {
			return nil, err
		}

		if connectionNo.Valid {
			connection.ConnectionNo = connectionNo.String
		}
		connection.ApplicationNo = connection.ApplicationNumber

		connection.ConnectionHolders, err =
			mapper.ParseConnectionHolders(
				connectionHoldersData,
			)
		if err != nil {
			return nil, err
		}

		connection.PlumberInfo, err =
			mapper.ParsePlumberInfo(
				plumberInfoData,
			)
		if err != nil {
			return nil, err
		}

		connection.RoadCuttingInfo, err =
			mapper.ParseRoadCuttingInfo(
				roadCuttingInfoData,
			)
		if err != nil {
			return nil, err
		}

		connection.ProcessInstance =
			dto.ProcessInstance{
				Action: "INITIATE",
			}

		connections =
			append(connections, connection)
	}

	return connections, nil
}

func (r *sewerageRepository) UpdateConnection(
	connection dto.SewerageConnection,
) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	connectionHoldersJSON, _ := json.Marshal(connection.ConnectionHolders)
	plumberInfoJSON, _ := json.Marshal(connection.PlumberInfo)
	roadCuttingInfoJSON, _ := json.Marshal(connection.RoadCuttingInfo)

	appNo := connection.ApplicationNo
	if appNo == "" {
		appNo = connection.ApplicationNumber
	}

	query := `
UPDATE eg_sw_connection SET
	connectiontype = $1,
	applicationstatus = $2,
	connectionno = $3,
	connection_holders = $4,
	plumber_info = $5,
	road_cutting_info = $6,
	lastmodifiedby = $7,
	lastmodifiedtime = $8
WHERE applicationno = $9
`

	lastModifiedBy := "system"
	if connection.AuditDetails != nil {
		lastModifiedBy = connection.AuditDetails.LastModifiedBy
	}

	currentTime := time.Now().UnixMilli()

	_, err := r.db.ExecContext(
		ctx,
		query,
		connection.ConnectionType,
		connection.ApplicationStatus,
		connection.ConnectionNo,
		connectionHoldersJSON,
		plumberInfoJSON,
		roadCuttingInfoJSON,
		lastModifiedBy,
		currentTime,
		appNo,
	)

	return err
}
