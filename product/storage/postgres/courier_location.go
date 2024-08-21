package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	pb "gitlab.com/bahodirova/product/genproto/product"
	"googlemaps.github.io/maps"
)

type CourierLocationManager struct {
	db         *sql.DB
	mapsClient *maps.Client
}

func NewCourierLocationManager(db *sql.DB, mapsAPIKey string) *CourierLocationManager {
	client, err := maps.NewClient(maps.WithAPIKey(mapsAPIKey))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &CourierLocationManager{db: db, mapsClient: client}
}

func (c *CourierLocationManager) CreateCourierLocation(req *pb.CreateCourierLocationRequest) (*pb.Empty, error) {
	id := uuid.NewString()

	var status string
	statusQuery := `SELECT status FROM courier WHERE id = $1`
	err := c.db.QueryRow(statusQuery, req.CourierId).Scan(&status)
	if err != nil {
		log.Println("Error while checking courier status", err)
		return nil, err
	}

	if status != "unemployed" {
		log.Println("Courier status is not unemployed")
		return nil, fmt.Errorf("courier status is not unemployed")
	}

	routes, err := c.getDirections(req.Latitude, req.Longitude)
	if err != nil {
		log.Println("Error getting directions from Google Maps", err)
		return nil, err
	}

	query := `
		INSERT INTO 
			courier_locations
				(id, 
				courier_id, 
				latitude, 
				longitude, 
				status,
				start_time,
				end_time)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7)`
	_, err = c.db.Exec(
		query,
		id,
		req.CourierId,
		req.Latitude,
		req.Longitude,
		req.Status,
		req.StartTime,
		req.EndTime)

	if err != nil {
		log.Println("Error while creating courier location", err)
		return nil, err
	}

	log.Printf("Routes: %+v\n", routes)

	return &pb.Empty{}, nil
}


func (c *CourierLocationManager) GetCourierLocation(req *pb.GetById) (*pb.CourierLocation, error) {

	query := `
		SELECT
			id,
			courier_id,
			latitude,
			longitude,
			status,
			start_time,
			end_time
		FROM
			courier_locations
		WHERE
			id = $1`

	row := c.db.QueryRow(query, req.Id)
	var location pb.CourierLocation
	err := row.Scan(
		&location.Id,
		&location.CourierId,
		&location.Latitude,
		&location.Longitude,
		&location.Status,
		&location.StartTime,
		&location.EndTime,
	)
	if err != nil {
		log.Println("no rows result set")
		return nil, err
	}
	return &location, nil
}

func (c *CourierLocationManager) GetAllCourierLocations(req *pb.GetAllCourierLocationsReq) (*pb.GetAllCourierLocationsRes, error) {

	query := `
		SELECT
			id,
			courier_id,
			latitude,
			longitude,
			status,
			start_time,
			end_time
		FROM
			courier_locations`

	var args []interface{}
	argCount := 1
	filters := []string{}

	if req.CourierId != "" {
		filters = append(filters, fmt.Sprintf("courier_id = $%d", argCount))
		args = append(args, req.CourierId)
		argCount++
	}
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}
	if req.Filter != nil {
		if req.Filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT $%d", argCount)
			args = append(args, req.Filter.Limit)
			argCount++
		}
		if req.Filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, req.Filter.Offset)
			argCount++
		}
	}
	rows, err := c.db.Query(query, args...)
	if err != nil {
		log.Println("no rows result set")
		return nil, err
	}
	defer rows.Close()
	var locations []*pb.CourierLocation
	for rows.Next() {
		var location pb.CourierLocation
		err := rows.Scan(
			&location.Id,
			&location.CourierId,
			&location.Latitude,
			&location.Longitude,
			&location.Status,
			&location.StartTime,
			&location.EndTime,
		)
		if err != nil {
			log.Println("no rows result set")
			return nil, err
		}
		locations = append(locations, &location)
	}
	return &pb.GetAllCourierLocationsRes{CourierLocations: locations}, nil
}

func (c *CourierLocationManager) UpdateCourierLocation(req *pb.UpdateCourierLocationRequest) (*pb.UpdateCourierLocationResponse, error) {
	var args []interface{}
	var conditions []string

	if req.CourierId != "" && req.CourierId != "string" {
		args = append(args, req.CourierId)
		conditions = append(conditions, fmt.Sprintf("courier_id = $%d", len(args)))
	}

	if req.Latitude != 0 {
		args = append(args, req.Latitude)
		conditions = append(conditions, fmt.Sprintf("latitude = $%d", len(args)))
	}

	if req.Longitude != 0 {
		args = append(args, req.Longitude)
		conditions = append(conditions, fmt.Sprintf("longitude = $%d", len(args)))
	}

	if req.StartTime != "" && req.StartTime != "string" {
		args = append(args, req.StartTime)
		conditions = append(conditions, fmt.Sprintf("start_time = $%d", len(args)))
	}

	if req.EndTime != "" && req.EndTime != "string" {
		args = append(args, req.EndTime)
		conditions = append(conditions, fmt.Sprintf("end_time = $%d", len(args)))
	}
	args = append(args, time.Now())
	conditions = append(conditions, fmt.Sprintf("updated_at = $%d", len(args)))

	query := `UPDATE courier_locations SET ` + strings.Join(conditions, ", ") + ` WHERE id = $` + fmt.Sprintf("%d", len(args)+1)
	args = append(args, req.Id)

	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &pb.UpdateCourierLocationResponse{Success: true, Message: "Courier location updated successfully"}, nil
}

func (c *CourierLocationManager) GetCourierLocationsByTimeRange(req *pb.GetCourierLocationsByTimeRangeReq) (*pb.GetCourierLocationsByTimeRangeRes, error) {
	query := `
	SELECT
		id,
		courier_id,
		latitude,
		longitude,
		status,
		start_time,
		end_time
	FROM
		courier_locations
	WHERE
		1=1
	`

	var args []interface{}
	argCount := 1

	if req.StartTime != "" {
		query += fmt.Sprintf(" AND start_time >= $%d", argCount)
		args = append(args, req.StartTime)
		argCount++
	}

	if req.EndTime != "" {
		query += fmt.Sprintf(" AND end_time <= $%d", argCount)
		args = append(args, req.EndTime)
		argCount++
	}

	if req.Filter != nil {
		if req.Filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT $%d", argCount)
			args = append(args, req.Filter.Limit)
			argCount++
		}
		if req.Filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, req.Filter.Offset)
			argCount++
		}
	}

	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courierLocations := []*pb.CourierLocation{}
	for rows.Next() {
		var courierLocation pb.CourierLocation
		err := rows.Scan(
			&courierLocation.Id,
			&courierLocation.CourierId,
			&courierLocation.Latitude,
			&courierLocation.Longitude,
			&courierLocation.Status,
			&courierLocation.StartTime,
			&courierLocation.EndTime,
		)
		if err != nil {
			log.Println("Error while scanning courier location", err)
			return nil, err
		}
		courierLocations = append(courierLocations, &courierLocation)
	}

	return &pb.GetCourierLocationsByTimeRangeRes{CourierLocations: courierLocations}, nil
}

func (c *CourierLocationManager) UpdateCourierLocationStatus(req *pb.UpdateCourierLocationStatusReq) (*pb.UpdateCourierLocationStatusRes, error) {

	query := `
	UPDATE
		courier_locations
	SET
		status = $1
	WHERE
		id = $2`

	_, err := c.db.Exec(
		query,
		req.Status,
		req.Id)

	if err != nil {
		log.Println("Error while updating courier location status", err)
		return nil, err
	}
	return &pb.UpdateCourierLocationStatusRes{Success: true, Message: "Courier location status updated successfully"}, nil
}

func (c *CourierLocationManager) getDirections(lat, lng float64) ([]maps.Route, error) {
	origin := fmt.Sprintf("%f,%f", lat, lng)
	destination := "some destination address" 

	req := &maps.DirectionsRequest{
		Origin:      origin,
		Destination: destination,
		Mode:        maps.TravelModeDriving,
	}

	routes, _, err := c.mapsClient.Directions(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return routes, nil
}