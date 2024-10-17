package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"fleetfy/backend/db"
	"fleetfy/backend/models"
	"fleetfy/backend/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BookingRequest struct {
	UserID        string             `json:"user_id"`
	VehicleType   string             `json:"vehicle_type"`
	PickupCoords  models.Coordinates `json:"pickup_coords"`
	DropoffCoords models.Coordinates `json:"dropoff_coords"`
	Distance      float64            `json:"distance"`
	Cost          float64            `json:"estimatedCost"`
}

func HandleBooking(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	var req BookingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Convert UserID to int64
	var userIDInt int64
	_, err = fmt.Sscanf(req.UserID, "%d", &userIDInt)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	pool := db.GetPool()

	// Fetch available vehicles
	queryAvailable := `SELECT id, vehicle_no, vehicle_type, latitude, longitude, busy FROM vehicles 
	                    WHERE vehicle_type = $1 AND busy = FALSE`
	rows, err := pool.Query(context.Background(), queryAvailable, req.VehicleType)
	if err != nil {
		http.Error(w, "Error fetching vehicles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var availableVehicles []models.Vehicle
	for rows.Next() {
		var vehicle models.Vehicle
		err := rows.Scan(&vehicle.ID, &vehicle.VehicleNo, &vehicle.VehicleType, &vehicle.Latitude, &vehicle.Longitude, &vehicle.Busy)
		if err != nil {
			http.Error(w, "Error scanning vehicle data", http.StatusInternalServerError)
			return
		}
		availableVehicles = append(availableVehicles, vehicle)
	}

	if len(availableVehicles) == 0 {
		http.Error(w, "No available vehicles found", http.StatusNotFound)
		return
	}

	// Calculate distance and find the closest vehicle
	var closestVehicle models.Vehicle
	var shortestDistance float64 = math.Inf(1)

	for _, vehicle := range availableVehicles {
		distance := utils.HaversineDistance(
			req.PickupCoords.Latitude, req.PickupCoords.Longitude,
			vehicle.Latitude, vehicle.Longitude,
		)
		if distance < shortestDistance {
			shortestDistance = distance
			closestVehicle = vehicle
		}
	}

	// Fetch the DriverID from the assignments table
	var driverID int64
	queryDriver := `SELECT uid FROM assignments WHERE vehicle_no = $1`
	err = pool.QueryRow(context.Background(), queryDriver, closestVehicle.VehicleNo).Scan(&driverID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			http.Error(w, "No driver assigned to this vehicle", http.StatusInternalServerError)
			return
		}
		http.Error(w, "Error fetching driver assignment", http.StatusInternalServerError)
		return
	}

	// Create a new booking
	insertBookingQuery := `INSERT INTO bookings 
	                        (user_id, vehicle_no, driver_id, pickup_latitude, pickup_longitude, 
	                         dropoff_latitude, dropoff_longitude, distance, cost, job_status, created_at) 
	                        VALUES 
	                        ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	var bookingID int64
	jobStatus := "in-transit"
	createdAt := time.Now()

	err = pool.QueryRow(context.Background(), insertBookingQuery,
		userIDInt,
		closestVehicle.VehicleNo,
		driverID,
		req.PickupCoords.Latitude,
		req.PickupCoords.Longitude,
		req.DropoffCoords.Latitude,
		req.DropoffCoords.Longitude,
		req.Distance,
		req.Cost,
		jobStatus,
		createdAt,
	).Scan(&bookingID)
	if err != nil {
		http.Error(w, "Error creating booking", http.StatusInternalServerError)
		return
	}

	// Update the assignment with the new booking ID
	updateAssignmentQuery := `UPDATE assignments SET booking_id = $1 WHERE uid = $2`
	_, err = pool.Exec(context.Background(), updateAssignmentQuery, bookingID, driverID)
	if err != nil {
		http.Error(w, "Error updating assignment with booking ID", http.StatusInternalServerError)
		return
	}

	// Mark the vehicle as busy
	updateVehicleQuery := `UPDATE vehicles SET busy = TRUE WHERE vehicle_no = $1`
	_, err = pool.Exec(context.Background(), updateVehicleQuery, closestVehicle.VehicleNo)
	if err != nil {
		http.Error(w, "Error updating vehicle status", http.StatusInternalServerError)
		return
	}

	// Respond with the created booking details
	booking := models.Booking{
		ID:              bookingID,
		UserID:          userIDInt,
		VehicleNo:       closestVehicle.VehicleNo,
		DriverID:        driverID,
		PickupLocation:  req.PickupCoords,
		DropoffLocation: req.DropoffCoords,
		Distance:        req.Distance,
		Cost:            req.Cost,
		JobStatus:       jobStatus,
		CreatedAt:       createdAt,
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"booking": booking,
	})
}

type BookingResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Booking models.Booking `json:"booking,omitempty"`
}

type BookingsResponse struct {
	Success  bool             `json:"success"`
	Message  string           `json:"message"`
	Bookings []models.Booking `json:"bookings,omitempty"`
}

// GetAllBookings retrieves all bookings.
func GetAllBookings(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	pool := db.GetPool()

	// Query to get all bookings
	query := `SELECT id, user_id, vehicle_no, driver_id, pickup_latitude, pickup_longitude, 
	                 dropoff_latitude, dropoff_longitude, distance, cost, job_status, created_at 
	          FROM bookings`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookings []models.Booking

	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.VehicleNo,
			&booking.DriverID,
			&booking.PickupLocation.Latitude,
			&booking.PickupLocation.Longitude,
			&booking.DropoffLocation.Latitude,
			&booking.DropoffLocation.Longitude,
			&booking.Distance,
			&booking.Cost,
			&booking.JobStatus,
			&booking.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Error scanning booking data", http.StatusInternalServerError)
			return
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating bookings", http.StatusInternalServerError)
		return
	}

	// Return all bookings
	json.NewEncoder(w).Encode(bookings)
}

// GetBookingByID retrieves booking details by booking ID.
func GetBookingByID(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	bookingID := params["bookingId"]
	if bookingID == "" {
		http.Error(w, "Missing Booking ID parameter", http.StatusBadRequest)
		return
	}

	// Convert bookingID to int64
	var bookingIDInt int64
	_, err := fmt.Sscanf(bookingID, "%d", &bookingIDInt)
	if err != nil {
		http.Error(w, "Invalid Booking ID format", http.StatusBadRequest)
		return
	}

	pool := db.GetPool()

	// Query to get booking by ID
	query := `SELECT id, user_id, vehicle_no, driver_id, pickup_latitude, pickup_longitude, 
	                 dropoff_latitude, dropoff_longitude, distance, cost, job_status, created_at 
	          FROM bookings WHERE id = $1`
	var booking models.Booking

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.QueryRow(ctx, query, bookingIDInt).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.VehicleNo,
		&booking.DriverID,
		&booking.PickupLocation.Latitude,
		&booking.PickupLocation.Longitude,
		&booking.DropoffLocation.Latitude,
		&booking.DropoffLocation.Longitude,
		&booking.Distance,
		&booking.Cost,
		&booking.JobStatus,
		&booking.CreatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(BookingResponse{Success: false, Message: "Booking not found"})
			return
		}
		http.Error(w, "Failed to fetch booking", http.StatusInternalServerError)
		return
	}

	// Return the booking details
	json.NewEncoder(w).Encode(BookingResponse{
		Success: true,
		Message: "Booking retrieved successfully",
		Booking: booking,
	})
}

// GetBookingsByUID retrieves bookings by user UID.
func GetBookingsByUID(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	uid := params["uid"]
	if uid == "" {
		http.Error(w, "Missing UID parameter", http.StatusBadRequest)
		return
	}

	// Convert UID to int64
	var uidInt int64
	_, err := fmt.Sscanf(uid, "%d", &uidInt)
	if err != nil {
		http.Error(w, "Invalid UID format", http.StatusBadRequest)
		return
	}

	pool := db.GetPool()

	// Query to get bookings by user UID
	query := `SELECT id, user_id, vehicle_no, driver_id, pickup_latitude, pickup_longitude, 
	                 dropoff_latitude, dropoff_longitude, distance, cost, job_status, created_at 
	          FROM bookings WHERE user_id = $1`
	rows, err := pool.Query(context.Background(), query, uidInt)
	if err != nil {
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookings []models.Booking

	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.VehicleNo,
			&booking.DriverID,
			&booking.PickupLocation.Latitude,
			&booking.PickupLocation.Longitude,
			&booking.DropoffLocation.Latitude,
			&booking.DropoffLocation.Longitude,
			&booking.Distance,
			&booking.Cost,
			&booking.JobStatus,
			&booking.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Error scanning booking data", http.StatusInternalServerError)
			return
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating bookings", http.StatusInternalServerError)
		return
	}

	if len(bookings) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(BookingsResponse{Success: false, Message: "No bookings found for this user"})
		return
	}

	// Return the bookings
	json.NewEncoder(w).Encode(BookingsResponse{
		Success:  true,
		Message:  "Bookings retrieved successfully",
		Bookings: bookings,
	})
}

// GetBookingsByDriverID retrieves bookings by driver ID.
func GetBookingsByDriverID(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	driverID := params["driverId"]
	if driverID == "" {
		http.Error(w, "Missing DriverId parameter", http.StatusBadRequest)
		return
	}

	// Convert DriverID to int64
	var driverIDInt int64
	_, err := fmt.Sscanf(driverID, "%d", &driverIDInt)
	if err != nil {
		http.Error(w, "Invalid DriverId format", http.StatusBadRequest)
		return
	}

	pool := db.GetPool()

	// Query to get bookings by driver ID
	query := `SELECT id, user_id, vehicle_no, driver_id, pickup_latitude, pickup_longitude, 
	                 dropoff_latitude, dropoff_longitude, distance, cost, job_status, created_at 
	          FROM bookings WHERE driver_id = $1`
	rows, err := pool.Query(context.Background(), query, driverIDInt)
	if err != nil {
		http.Error(w, "Failed to fetch bookings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookings []models.Booking

	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.VehicleNo,
			&booking.DriverID,
			&booking.PickupLocation.Latitude,
			&booking.PickupLocation.Longitude,
			&booking.DropoffLocation.Latitude,
			&booking.DropoffLocation.Longitude,
			&booking.Distance,
			&booking.Cost,
			&booking.JobStatus,
			&booking.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Error scanning booking data", http.StatusInternalServerError)
			return
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating bookings", http.StatusInternalServerError)
		return
	}

	if len(bookings) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(BookingsResponse{Success: false, Message: "No bookings found for this driver"})
		return
	}

	// Return the bookings
	json.NewEncoder(w).Encode(BookingsResponse{
		Success:  true,
		Message:  "Bookings retrieved successfully",
		Bookings: bookings,
	})
}
