package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"fleetfy/backend/db"
	"fleetfy/backend/models"
	"fleetfy/backend/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func CompleteJobHandler(w http.ResponseWriter, r *http.Request) {
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

	// Start a transaction
	tx, err := pool.Begin(context.Background())
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	// Fetch booking details using booking ID
	var booking models.Booking
	queryBooking := `SELECT id, user_id, vehicle_no, driver_id, pickup_latitude, pickup_longitude, 
	                        dropoff_latitude, dropoff_longitude, distance, cost, job_status, created_at 
	                 FROM bookings WHERE id = $1`
	err = tx.QueryRow(context.Background(), queryBooking, bookingIDInt).Scan(
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
			http.Error(w, `{"message": "Booking not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch booking", http.StatusInternalServerError)
		return
	}

	// Fetch vehicle details using vehicle number
	var vehicle models.Vehicle
	queryVehicle := `SELECT id, vehicle_no, vehicle_type, latitude, longitude, busy FROM vehicles WHERE vehicle_no = $1`
	err = tx.QueryRow(context.Background(), queryVehicle, booking.VehicleNo).Scan(
		&vehicle.ID,
		&vehicle.VehicleNo,
		&vehicle.VehicleType,
		&vehicle.Latitude,
		&vehicle.Longitude,
		&vehicle.Busy,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			http.Error(w, `{"message": "Vehicle not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch vehicle", http.StatusInternalServerError)
		return
	}

	// Calculate distance using the HaversineDistance function
	distance := utils.HaversineDistance(
		vehicle.Latitude, vehicle.Longitude,
		booking.DropoffLocation.Latitude, booking.DropoffLocation.Longitude,
	)

	if distance <= 5.0 {
		// Update job status to completed
		updateJobStatusQuery := `UPDATE bookings SET job_status = $1 WHERE id = $2`
		_, err = tx.Exec(context.Background(), updateJobStatusQuery, "completed", booking.ID)
		if err != nil {
			http.Error(w, `{"message": "Failed to update job status"}`, http.StatusInternalServerError)
			return
		}

		// Mark vehicle as not busy
		updateVehicleStatusQuery := `UPDATE vehicles SET busy = FALSE WHERE id = $1`
		_, err = tx.Exec(context.Background(), updateVehicleStatusQuery, vehicle.ID)
		if err != nil {
			http.Error(w, `{"message": "Failed to update vehicle status"}`, http.StatusInternalServerError)
			return
		}

		// Clear booking ID in assignments
		updateAssignmentQuery := `UPDATE assignments SET booking_id = NULL WHERE vehicle_no = $1`
		_, err = tx.Exec(context.Background(), updateAssignmentQuery, booking.VehicleNo)
		if err != nil {
			http.Error(w, "Failed to update assignment", http.StatusInternalServerError)
			return
		}

		// Commit the transaction
		if err := tx.Commit(context.Background()); err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		// Respond with success message
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Job marked as completed",
		})
	} else {
		// Vehicle is out of range
		http.Error(w, `{"message": "You are not in range of the delivery location"}`, http.StatusBadRequest)
		return
	}
}
