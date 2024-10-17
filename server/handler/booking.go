package handler

import (
    "database/sql"
    "encoding/json"
   
    "strconv"
    "fmc/models"
	"github.com/gorilla/mux"
    "net/http"
    "log"
	
)

func CreateBookingHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            PickupLocation  string  `json:"pickup_location"`
            DropoffLocation string  `json:"dropoff_location"`
            VehicleType     string  `json:"vehicle_type"`
            EstimatedCost   float64 `json:"estimated_cost"`
        }

        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        // Get User-ID from the header
        userIDStr := r.Header.Get("User-ID")
        if userIDStr == "" {
            http.Error(w, "User ID is required", http.StatusBadRequest)
            return
        }

        userID, err := strconv.Atoi(userIDStr)
        if err != nil {
            http.Error(w, "Invalid User ID", http.StatusBadRequest)
            return
        }

        // Create the booking
        bookingID, err := models.CreateBooking(db, userID, req.PickupLocation, req.DropoffLocation, req.VehicleType, req.EstimatedCost)
        if err != nil {
            http.Error(w, "Could not create booking", http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]interface{}{
            "message": "Booking created",
            "booking_id": bookingID,
        })
    }
}


func AcceptBookingHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        driverIDStr := r.Header.Get("Driver-ID")
        driverID, err := strconv.Atoi(driverIDStr)
        if err != nil {
            log.Printf("Invalid Driver ID: %s", driverIDStr)
            http.Error(w, "Invalid Driver ID", http.StatusBadRequest)
            return
        }

        // Get the booking ID from the URL path
        vars := mux.Vars(r)
        bookingID, err := strconv.Atoi(vars["id"])
        if err != nil {
            log.Printf("Invalid booking ID: %v", err)
            http.Error(w, "Invalid booking ID", http.StatusBadRequest)
            return
        }

        log.Printf("Driver ID: %d is attempting to accept Booking ID: %d", driverID, bookingID)

        // Try to accept the booking in the database
        err = models.AcceptBooking(db, driverID, bookingID)
        if err != nil {
            log.Printf("Error accepting booking: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]string{
            "message": "Booking accepted",
        })
    }
}
// GetPendingBookingsHandler fetches unassigned bookings for drivers
func GetPendingBookingsHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        role := r.Header.Get("Role")  // Get role from headers
        driverID := r.Header.Get("Driver-ID")  // Get driver ID if applicable

        // Ensure the role is "driver"
        if role != "driver" || driverID == "" {
            http.Error(w, "Unauthorized access", http.StatusUnauthorized)
            return
        }

        // Fetch only unassigned pending bookings
        rows, err := db.Query(`SELECT id, user_id, pickup_location, dropoff_location, vehicle_type, estimated_cost, status FROM bookings WHERE status = 'pending' AND driver_id IS NULL`)
        if err != nil {
            log.Printf("Error fetching pending bookings: %v", err)
            http.Error(w, "Error fetching bookings", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var bookings []struct {
            ID             int     `json:"id"`
            UserID         int     `json:"user_id"`
            PickupLocation string  `json:"pickup_location"`
            DropoffLocation string `json:"dropoff_location"`
            VehicleType    string  `json:"vehicle_type"`
            EstimatedCost  float64 `json:"estimated_cost"`
            Status         string  `json:"status"`
        }

        for rows.Next() {
            var booking struct {
                ID             int     `json:"id"`
                UserID         int     `json:"user_id"`
                PickupLocation string  `json:"pickup_location"`
                DropoffLocation string `json:"dropoff_location"`
                VehicleType    string  `json:"vehicle_type"`
                EstimatedCost  float64 `json:"estimated_cost"`
                Status         string  `json:"status"`
            }
            err := rows.Scan(&booking.ID, &booking.UserID, &booking.PickupLocation, &booking.DropoffLocation, &booking.VehicleType, &booking.EstimatedCost, &booking.Status)
            if err != nil {
                log.Printf("Error scanning booking row: %v", err)
                http.Error(w, "Error processing bookings", http.StatusInternalServerError)
                return
            }
            bookings = append(bookings, booking)
        }

        if err := rows.Err(); err != nil {
            log.Printf("Error after rows.Next(): %v", err)
            http.Error(w, "Error processing bookings", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(bookings)
    }
}
