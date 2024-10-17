package analytics

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"fleetfy/backend/db"

	// "github.com/jackc/pgx/v4/pgxpool"
)

type BookingStatus struct {
	TotalBookings int `json:"total_bookings"`
	Completed     int `json:"completed"`
	Pending       int `json:"pending"`
}

func GetBookingAnalysis(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool := db.GetPool()

	// Get total bookings
	var totalCount int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM bookings").Scan(&totalCount)
	if err != nil {
		http.Error(w, "Failed to fetch total bookings", http.StatusInternalServerError)
		return
	}

	// Get completed bookings
	var completedCount int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM bookings WHERE job_status = $1", "completed").Scan(&completedCount)
	if err != nil {
		http.Error(w, "Failed to fetch completed bookings", http.StatusInternalServerError)
		return
	}

	// Calculate pending bookings
	pendingCount := totalCount - completedCount

	// Create response
	response := BookingStatus{
		TotalBookings: totalCount,
		Completed:     completedCount,
		Pending:       pendingCount,
	}

	// Encode response to JSON
	json.NewEncoder(w).Encode(response)
}
