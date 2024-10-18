package handler

import (
    "database/sql"
    "encoding/json"
    "net/http"
	"log"
    "github.com/gorilla/mux"
    "time"
    
    "fmc/models"
)

func GetAllVehiclesHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vehicles, err := models.FetchAllVehicles(db)
        if err != nil {
            http.Error(w, "Could not fetch vehicles", http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(vehicles)
    }
}

func CreateVehicleHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Type         string `json:"type"`
            Availability bool   `json:"availability"`
        }

        // Parse the JSON request body
        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        // Create the vehicle in the database
        vehicleID, err := models.CreateVehicle(db, req.Type, req.Availability)
        if err != nil {
            log.Printf("Error creating vehicle: %v", err)
            http.Error(w, "Could not create vehicle", http.StatusInternalServerError)
            return
        }

        // Return a success response
        json.NewEncoder(w).Encode(map[string]interface{}{
            "message":    "Vehicle created",
            "vehicle_id": vehicleID,
        })
    }
}

// GetAllBookingsHandler fetches all bookings for admin
func GetAllBookingsHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, err := db.Query(`SELECT id, user_id, driver_id, pickup_location, dropoff_location, vehicle_type, estimated_cost, status FROM bookings`)
        if err != nil {
            log.Printf("Error fetching bookings: %v", err)
            http.Error(w, "Error fetching bookings", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var bookings []struct {
            ID             int     `json:"id"`
            UserID         int     `json:"user_id"`
            DriverID       *int    `json:"driver_id"` // Pointer to handle NULL values
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
                DriverID       *int    `json:"driver_id"`
                PickupLocation string  `json:"pickup_location"`
                DropoffLocation string `json:"dropoff_location"`
                VehicleType    string  `json:"vehicle_type"`
                EstimatedCost  float64 `json:"estimated_cost"`
                Status         string  `json:"status"`
            }
            err := rows.Scan(&booking.ID, &booking.UserID, &booking.DriverID, &booking.PickupLocation, &booking.DropoffLocation, &booking.VehicleType, &booking.EstimatedCost, &booking.Status)
            if err != nil {
                log.Printf("Error scanning booking row: %v", err)
                http.Error(w, "Error processing bookings", http.StatusInternalServerError)
                return
            }
            bookings = append(bookings, booking)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(bookings)
    }
}

// CompleteBookingHandler marks a booking as complete
func CompleteBookingHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        bookingID := mux.Vars(r)["id"]

        query := `UPDATE bookings SET status = 'completed' WHERE id = $1 AND status = 'accepted'`
        result, err := db.Exec(query, bookingID)
        if err != nil {
            log.Printf("Error marking booking as complete: %v", err)
            http.Error(w, "Error completing booking", http.StatusInternalServerError)
            return
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil || rowsAffected == 0 {
            http.Error(w, "No booking found or booking is not in an accepted state", http.StatusBadRequest)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Booking marked as complete"})
    }
}

// GetDriverActiveBookingsCount fetches the count of active bookings for each driver
func GetDriverActiveBookingsCount(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, err := db.Query(`
            SELECT driver_id, COUNT(*) as active_bookings
            FROM bookings
            WHERE status = 'accepted'
            GROUP BY driver_id
        `)
        if err != nil {
            log.Printf("Error fetching driver active bookings: %v", err)
            http.Error(w, "Error fetching driver active bookings", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var driverBookings []struct {
            DriverID      int `json:"driver_id"`
            ActiveBookings int `json:"active_bookings"`
        }

        for rows.Next() {
            var db struct {
                DriverID      int `json:"driver_id"`
                ActiveBookings int `json:"active_bookings"`
            }
            err := rows.Scan(&db.DriverID, &db.ActiveBookings)
            if err != nil {
                log.Printf("Error scanning row: %v", err)
                http.Error(w, "Error processing data", http.StatusInternalServerError)
                return
            }
            driverBookings = append(driverBookings, db)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(driverBookings)
    }
}

type VehicleStatus struct {
    Active int `json:"active"`
    Idle   int `json:"idle"`
}

func GetVehicleStatus(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var status VehicleStatus

        // Count active vehicles (in use)
        err := db.QueryRow(`SELECT COUNT(*) FROM vehicles WHERE availability = FALSE`).Scan(&status.Active)
        if err != nil {
            http.Error(w, "Error fetching active vehicles", http.StatusInternalServerError)
            return
        }

        // Count idle vehicles (available)
        err = db.QueryRow(`SELECT COUNT(*) FROM vehicles WHERE availability = TRUE`).Scan(&status.Idle)
        if err != nil {
            http.Error(w, "Error fetching idle vehicles", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(status)
    }
}
type DriverPerformance struct {
    DriverName   string `json:"driver_name"`
    Deliveries   int    `json:"deliveries"`
}

func GetDriverPerformance(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, err := db.Query(`
            SELECT drivers.name, COUNT(bookings.id) 
            FROM drivers 
            JOIN bookings ON drivers.id = bookings.driver_id
            WHERE bookings.status = 'completed'
            GROUP BY drivers.name
        `)
        if err != nil {
            http.Error(w, "Error fetching driver performance", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var data []DriverPerformance
        for rows.Next() {
            var performance DriverPerformance
            if err := rows.Scan(&performance.DriverName, &performance.Deliveries); err != nil {
                http.Error(w, "Error processing data", http.StatusInternalServerError)
                return
            }
            data = append(data, performance)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(data)
    }
}
type RevenueData struct {
    Date    time.Time `json:"date"`
    Revenue float64   `json:"revenue"`
}

func GetRevenueOverTime(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, err := db.Query(`
            SELECT date_trunc('day', created_at) AS day, SUM(estimated_cost) 
            FROM bookings 
            WHERE status = 'completed' 
            GROUP BY day
            ORDER BY day
        `)
        if err != nil {
            http.Error(w, "Error fetching revenue data", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var data []RevenueData
        for rows.Next() {
            var revenue RevenueData
            if err := rows.Scan(&revenue.Date, &revenue.Revenue); err != nil {
                http.Error(w, "Error processing data", http.StatusInternalServerError)
                return
            }
            data = append(data, revenue)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(data)
    }
}
type BookingStatus struct {
    Status string `json:"status"`
    Count  int    `json:"count"`
}

func GetBookingStatusDistribution(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, err := db.Query(`
            SELECT status, COUNT(*) 
            FROM bookings 
            GROUP BY status
        `)
        if err != nil {
            http.Error(w, "Error fetching booking statuses", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var data []BookingStatus
        for rows.Next() {
            var status BookingStatus
            if err := rows.Scan(&status.Status, &status.Count); err != nil {
                http.Error(w, "Error processing data", http.StatusInternalServerError)
                return
            }
            data = append(data, status)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(data)
    }
}