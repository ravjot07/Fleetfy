package vehicles

import (
	"context"
	"encoding/json"
	
	"fleetfy/backend/db"
	"fleetfy/backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type VehicleRequest struct {
	VehicleNo   string `json:"vehicle_no"`
	VehicleType string `json:"vehicle_type"`
}

type VehicleResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Vehicle models.Vehicle `json:"vehicle,omitempty"`
}

type VehiclesResponse struct {
	Success  bool             `json:"success"`
	Message  string           `json:"message"`
	Vehicles []models.Vehicle `json:"vehicles,omitempty"`
}

// GetVehicleInfo retrieves vehicle information by vehicle number
func GetVehicleInfoById(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	vehicleNo := params["vehicleNo"]
	if vehicleNo == "" {
		http.Error(w, "Missing Vehicle Number parameter", http.StatusBadRequest)
		return
	}

	pool := db.GetPool()

	// Query to fetch vehicle by vehicle_no
	query := `SELECT id, vehicle_no, vehicle_type, latitude, longitude, busy, created_at 
	          FROM vehicles WHERE vehicle_no = $1`

	var vehicle models.Vehicle

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := pool.QueryRow(ctx, query, vehicleNo).Scan(
		&vehicle.ID,
		&vehicle.VehicleNo,
		&vehicle.VehicleType,
		&vehicle.Latitude,
		&vehicle.Longitude,
		&vehicle.Busy,
		&vehicle.CreatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(VehicleResponse{Success: false, Message: "Vehicle not found"})
			return
		}
		http.Error(w, "Failed to fetch vehicle", http.StatusInternalServerError)
		return
	}

	// Return the vehicle details
	json.NewEncoder(w).Encode(VehicleResponse{
		Success: true,
		Message: "Vehicle retrieved successfully",
		Vehicle: vehicle,
	})
}

// GetAllVehicles retrieves all vehicles from the database
func GetAllVehicles(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	pool := db.GetPool()

	// Query to fetch all vehicles
	query := `SELECT id, vehicle_no, vehicle_type, latitude, longitude, busy, created_at FROM vehicles`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Failed to fetch vehicles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var vehicles []models.Vehicle

	for rows.Next() {
		var vehicle models.Vehicle
		err := rows.Scan(
			&vehicle.ID,
			&vehicle.VehicleNo,
			&vehicle.VehicleType,
			&vehicle.Latitude,
			&vehicle.Longitude,
			&vehicle.Busy,
			&vehicle.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Failed to decode vehicle", http.StatusInternalServerError)
			return
		}
		vehicles = append(vehicles, vehicle)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Cursor error while fetching vehicles", http.StatusInternalServerError)
		return
	}

	// Return the list of vehicles
	json.NewEncoder(w).Encode(VehiclesResponse{
		Success:  true,
		Message:  "Vehicles retrieved successfully",
		Vehicles: vehicles,
	})
}

// AddVehicleHandler adds a new vehicle to the database
func AddVehicleHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var vehicle models.Vehicle
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&vehicle)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	pool := db.GetPool()

	// Check if the vehicle already exists
	checkQuery := `SELECT COUNT(*) FROM vehicles WHERE vehicle_no = $1`
	var count int
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.QueryRow(ctx, checkQuery, vehicle.VehicleNo).Scan(&count)
	if err != nil {
		http.Error(w, "Failed to check existing vehicle", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		// Vehicle already exists
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Vehicle already exists"})
		return
	}

	// Set default values if not provided
	if vehicle.Coordinates.Latitude == 0 && vehicle.Coordinates.Longitude == 0 {
		vehicle.Coordinates = models.Coordinates{Latitude: 28.612894, Longitude: 77.216721} // Default coordinates
	}
	vehicle.Busy = false // Default busy status

	// Insert vehicle into the database
	insertQuery := `INSERT INTO vehicles 
	                (vehicle_no, vehicle_type, latitude, longitude, busy, created_at) 
	                VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	var vehicleID int64
	createdAt := time.Now()

	err = pool.QueryRow(ctx, insertQuery,
		vehicle.VehicleNo,
		vehicle.VehicleType,
		vehicle.Coordinates.Latitude,
		vehicle.Coordinates.Longitude,
		vehicle.Busy,
		createdAt,
	).Scan(&vehicle.ID, &vehicle.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to add vehicle", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(VehicleResponse{
		Success: true,
		Message: "Vehicle added successfully",
		Vehicle: vehicle,
	})
}

// Struct for the request payload to update location
type LocationUpdate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UID       string  `json:"uid"`
}

// UpdateVehicleLocation updates the location of a vehicle by its vehicle number
func UpdateVehicleLocation(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Extract vehicle_no from the URL
	vars := mux.Vars(r)
	vehicleNo := vars["vehicle_no"]
	if vehicleNo == "" {
		http.Error(w, "Missing vehicle_no parameter", http.StatusBadRequest)
		return
	}

	// Decode the request body
	var locationUpdate LocationUpdate
	err := json.NewDecoder(r.Body).Decode(&locationUpdate)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	pool := db.GetPool()

	// Convert UID to int64 if necessary (assuming UID is part of the request, adjust based on your logic)
	var uidInt int64
	if locationUpdate.UID != "" {
		uidInt, err = strconv.ParseInt(locationUpdate.UID, 10, 64)
		if err != nil {
			http.Error(w, "Invalid UID format", http.StatusBadRequest)
			return
		}
	}

	// Update the vehicle's coordinates
	updateQuery := `UPDATE vehicles SET latitude = $1, longitude = $2 WHERE vehicle_no = $3`
	result, err := pool.Exec(context.Background(), updateQuery, locationUpdate.Latitude, locationUpdate.Longitude, vehicleNo)
	if err != nil {
		http.Error(w, "Failed to update vehicle location", http.StatusInternalServerError)
		return
	}

	// Check if the vehicle was found and updated
	if result.RowsAffected() == 0 {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	// Return success response
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Vehicle location updated successfully",
	})
}

// GetVehicleCoordsHandler retrieves the coordinates of a vehicle by its vehicle number
func GetVehicleCoordsHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Get vehicle_no from the URL parameters
	vars := mux.Vars(r)
	vehicleNo := vars["vehicle_no"]

	pool := db.GetPool()

	// Query to fetch vehicle coordinates by vehicle_no
	query := `SELECT latitude, longitude FROM vehicles WHERE vehicle_no = $1`

	var coords models.Coordinates

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := pool.QueryRow(ctx, query, vehicleNo).Scan(&coords.Latitude, &coords.Longitude)
	if err != nil {
		if err.Error() == "no rows in result set" {
			http.Error(w, "Vehicle not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch vehicle coordinates", http.StatusInternalServerError)
		return
	}

	// Return the vehicle coordinates as JSON
	json.NewEncoder(w).Encode(coords)
}
