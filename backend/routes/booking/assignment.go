package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"fleetfy/backend/db"
	"fleetfy/backend/models"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AssignmentResponse struct {
	Success    bool              `json:"success"`
	Message    string            `json:"message"`
	Assignment models.Assignment `json:"assignment,omitempty"`
}

// GetAssignmentByUid retrieves the assignment details of a driver based on their UID.
func GetAssignmentByUid(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Get the 'uid' from the URL params
	params := mux.Vars(r)
	uid := params["uid"]
	if uid == "" {
		http.Error(w, "Missing UID parameter", http.StatusBadRequest)
		return
	}

	// Convert uid to int64
	var uidInt int64
	_, err := fmt.Sscanf(uid, "%d", &uidInt)
	if err != nil {
		http.Error(w, "Invalid UID format", http.StatusBadRequest)
		return
	}

	pool := db.GetPool()

	// Query for the assignment based on the UID
	query := `SELECT id, uid, vehicle_no, booking_id FROM assignments WHERE uid = $1`
	var assignment models.Assignment

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.QueryRow(ctx, query, uidInt).Scan(&assignment.ID, &assignment.UID, &assignment.VehicleNo, &assignment.BookingID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(AssignmentResponse{Success: false, Message: "Assignment not found"})
			return
		}
		http.Error(w, "Failed to fetch assignment", http.StatusInternalServerError)
		return
	}

	// Return the assignment details
	json.NewEncoder(w).Encode(AssignmentResponse{
		Success:    true,
		Message:    "Assignment retrieved successfully",
		Assignment: assignment,
	})
}

// GetAssignmentByVehicleNo retrieves the assignment details of a vehicle based on its vehicle number.
func GetAssignmentByVehicleNo(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Get the 'vehicle_no' from the URL params
	params := mux.Vars(r)
	vehicleNo := params["vehicle_no"]
	if vehicleNo == "" {
		http.Error(w, "Missing vehicle_no parameter", http.StatusBadRequest)
		return
	}

	pool := db.GetPool()

	// Query for the assignment based on the vehicle number
	query := `SELECT id, uid, vehicle_no, booking_id FROM assignments WHERE vehicle_no = $1`
	var assignment models.Assignment

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := pool.QueryRow(ctx, query, vehicleNo).Scan(&assignment.ID, &assignment.UID, &assignment.VehicleNo, &assignment.BookingID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(AssignmentResponse{Success: false, Message: "Assignment not found"})
			return
		}
		http.Error(w, "Failed to fetch assignment", http.StatusInternalServerError)
		return
	}

	// Return the assignment details
	json.NewEncoder(w).Encode(AssignmentResponse{
		Success:    true,
		Message:    "Assignment retrieved successfully",
		Assignment: assignment,
	})
}

// GetAllAssignments retrieves all assignments.
func GetAllAssignments(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	pool := db.GetPool()

	// Query to get all assignments
	query := `SELECT id, uid, vehicle_no, booking_id FROM assignments`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Failed to fetch assignments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var assignments []models.Assignment

	for rows.Next() {
		var assignment models.Assignment
		err := rows.Scan(&assignment.ID, &assignment.UID, &assignment.VehicleNo, &assignment.BookingID)
		if err != nil {
			http.Error(w, "Failed to decode assignment", http.StatusInternalServerError)
			return
		}
		assignments = append(assignments, assignment)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating assignments", http.StatusInternalServerError)
		return
	}

	// Return all assignments
	json.NewEncoder(w).Encode(assignments)
}

// AssignVehicle assigns a vehicle to a driver based on UID.
func AssignVehicle(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Get the 'uid' from the URL params
	params := mux.Vars(r)
	uid := params["uid"]
	if uid == "" {
		http.Error(w, "Missing UID parameter", http.StatusBadRequest)
		return
	}

	var assignmentUpdate struct {
		VehicleNo string `json:"vehicle_no"`
	}

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&assignmentUpdate); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Convert uid to int64
	var uidInt int64
	_, err := fmt.Sscanf(uid, "%d", &uidInt)
	if err != nil {
		http.Error(w, "Invalid UID format", http.StatusBadRequest)
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

	// Check if the vehicle exists and is not busy
	var isBusy bool
	checkVehicleQuery := `SELECT busy FROM vehicles WHERE vehicle_no = $1`
	err = tx.QueryRow(context.Background(), checkVehicleQuery, assignmentUpdate.VehicleNo).Scan(&isBusy)
	if err != nil {
		if err.Error() == "no rows in result set" {
			http.Error(w, "Vehicle not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to check vehicle status", http.StatusInternalServerError)
		return
	}

	if isBusy {
		http.Error(w, "Vehicle is currently busy", http.StatusConflict)
		return
	}

	// Update the assignment with the new vehicle number
	updateAssignmentQuery := `UPDATE assignments SET vehicle_no = $1 WHERE uid = $2`
	result, err := tx.Exec(context.Background(), updateAssignmentQuery, assignmentUpdate.VehicleNo, uidInt)
	if err != nil {
		http.Error(w, "Failed to update assignment", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "No assignment found for the given UID", http.StatusNotFound)
		return
	}

	// Mark the vehicle as busy
	updateVehicleQuery := `UPDATE vehicles SET busy = TRUE WHERE vehicle_no = $1`
	_, err = tx.Exec(context.Background(), updateVehicleQuery, assignmentUpdate.VehicleNo)
	if err != nil {
		http.Error(w, "Failed to update vehicle status", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Vehicle assigned successfully",
	})
}
