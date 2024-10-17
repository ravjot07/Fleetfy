package analytics

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"fleetfy/backend/db"

	// "github.com/jackc/pgx/v4/pgxpool"
)

type DriverStatus struct {
	TotalDrivers int `json:"total_drivers"`
	NotVerified  int `json:"not_verified"`
	Free         int `json:"free"`
	Busy         int `json:"busy"`
}

func GetDriverAnalysis(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool := db.GetPool()

	// Get total drivers
	var totalCount int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM assignments").Scan(&totalCount)
	if err != nil {
		http.Error(w, "Failed to fetch total drivers", http.StatusInternalServerError)
		return
	}

	// Get not verified drivers (vehicle_no is empty)
	var notVerifiedCount int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM assignments WHERE vehicle_no = $1", "").Scan(&notVerifiedCount)
	if err != nil {
		http.Error(w, "Failed to fetch not verified drivers", http.StatusInternalServerError)
		return
	}

	// Get free drivers (vehicle_no is not empty and booking_id is empty)
	var freeCount int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM assignments WHERE vehicle_no != $1 AND booking_id = $2", "", "").Scan(&freeCount)
	if err != nil {
		http.Error(w, "Failed to fetch free drivers", http.StatusInternalServerError)
		return
	}

	// Calculate busy drivers (total - not_verified - free)
	busyCount := totalCount - notVerifiedCount - freeCount

	// Create response
	response := DriverStatus{
		TotalDrivers: totalCount,
		NotVerified:  notVerifiedCount,
		Free:         freeCount,
		Busy:         busyCount,
	}

	// Encode response to JSON
	json.NewEncoder(w).Encode(response)
}
