package handler

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"
    "fmc/models"
	"github.com/gorilla/mux"
	
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
        driverID, _ := strconv.Atoi(r.Header.Get("Driver-ID"))  // Assuming Driver ID from header

        // Get the booking ID from the URL path
        vars := mux.Vars(r)
        bookingID, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid booking ID", http.StatusBadRequest)
            return
        }

        err = models.AcceptBooking(db, driverID, bookingID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(map[string]string{
            "message": "Booking accepted",
        })
    }
}
