package analytics

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"fleetfy/backend/db"

	// "github.com/jackc/pgx/v4/pgxpool"
)

type VehicleStatus struct {
	TotalVehicles int `json:"total_vehicles"`
	Free          int `json:"free"`
	Busy          int `json:"busy"`
}

func GetVehicleAnalysis(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool := db.GetPool()

	// Get total vehicles
	var totalCount int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM vehicles").Scan(&totalCount)
	if err != nil {
		http.Error(w, "Failed to fetch total vehicles", http.StatusInternalServerError)
		return
	}

	// Get free vehicles (busy: false)
	var freeCount int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM vehicles WHERE busy = $1", false).Scan(&freeCount)
	if err != nil {
		http.Error(w, "Failed to fetch free vehicles", http.StatusInternalServerError)
		return
	}

	// Calculate busy vehicles (total - free)
	busyCount := totalCount - freeCount

	// Create response
	response := VehicleStatus{
		TotalVehicles: totalCount,
		Free:          freeCount,
		Busy:          busyCount,
	}

	// Encode response to JSON
	json.NewEncoder(w).Encode(response)
}
