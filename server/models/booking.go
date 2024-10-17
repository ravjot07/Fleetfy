package models

import (
    "database/sql"
    "errors"
    "time"
	"log"
)

type Booking struct {
    ID             int       `json:"id"`
    UserID         int       `json:"user_id"`
    DriverID       int       `json:"driver_id"`
    VehicleID      int       `json:"vehicle_id"`
    PickupLocation string    `json:"pickup_location"`
    DropoffLocation string   `json:"dropoff_location"`
    VehicleType    string    `json:"vehicle_type"`
    EstimatedCost  float64   `json:"estimated_cost"`
    Status         string    `json:"status"`
    CreatedAt      time.Time `json:"created_at"`
}

// CreateBooking creates a new booking for a user
func CreateBooking(db *sql.DB, userID int, pickupLocation, dropoffLocation, vehicleType string, estimatedCost float64) (int, error) {
    var bookingID int
    query := `
        INSERT INTO bookings (user_id, pickup_location, dropoff_location, vehicle_type, estimated_cost, status)
        VALUES ($1, $2, $3, $4, $5, 'pending') RETURNING id`
    
    err := db.QueryRow(query, userID, pickupLocation, dropoffLocation, vehicleType, estimatedCost).Scan(&bookingID)
    if err != nil {
        log.Printf("Error executing SQL query: %v", err) // Log the query error
        return 0, err
    }

    return bookingID, nil
}

func AcceptBooking(db *sql.DB, driverID, bookingID int) error {
    // Assign booking to driver and update status
    query := `UPDATE bookings SET driver_id=$1, status='accepted' WHERE id=$2 AND status='pending'`
    result, err := db.Exec(query, driverID, bookingID)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        return errors.New("Booking not available or already accepted")
    }

    return nil
}
