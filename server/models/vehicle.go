package models

import (
    "database/sql"
	"log"
)

// Vehicle represents a vehicle in the fleet
type Vehicle struct {
    ID          int    `json:"id"`
    Type        string `json:"type"`
    Availability bool   `json:"availability"`
    DriverID    *int    `json:"driver_id"`
}

// FetchAllVehicles fetches all vehicles from the database
func FetchAllVehicles(db *sql.DB) ([]Vehicle, error) {
    log.Println("Fetching all vehicles...")

    rows, err := db.Query(`SELECT id, type, availability, driver_id FROM vehicles`)
    if err != nil {
        log.Printf("Error executing query: %v", err)
        return nil, err
    }
    defer rows.Close()

    vehicles := []Vehicle{}
    for rows.Next() {
        var v Vehicle
        err := rows.Scan(&v.ID, &v.Type, &v.Availability, &v.DriverID)
        if err != nil {
            log.Printf("Error scanning row: %v", err)
            return nil, err
        }
        log.Printf("Vehicle fetched: %+v", v) // Logging each vehicle fetched
        vehicles = append(vehicles, v)
    }

    // Check for errors after iterating over rows
    if err = rows.Err(); err != nil {
        log.Printf("Error during row iteration: %v", err)
        return nil, err
    }

    log.Printf("Total vehicles fetched: %d", len(vehicles))
    return vehicles, nil
}
func CreateVehicle(db *sql.DB, vehicleType string, availability bool) (int, error) {
    var vehicleID int
    query := `INSERT INTO vehicles (type, availability) VALUES ($1, $2) RETURNING id`
    err := db.QueryRow(query, vehicleType, availability).Scan(&vehicleID)
    if err != nil {
        return 0, err
    }
    return vehicleID, nil
}