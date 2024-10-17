package handler

import (
    "database/sql"
    "encoding/json"
    "net/http"
	"log"
    
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