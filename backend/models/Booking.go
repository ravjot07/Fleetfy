package models

import "time"

type Booking struct {
    ID              int64       `json:"id"`
    UserID          int64       `json:"user_id"`
    VehicleNo       string      `json:"vehicle_no"`
    DriverID        int64       `json:"driver_id"`
    PickupLocation  Coordinates `json:"pickup_location"`
    DropoffLocation Coordinates `json:"dropoff_location"`
    Distance        float64     `json:"distance"`
    Cost            float64     `json:"cost"`
    JobStatus       string      `json:"job_status"`
    CreatedAt       time.Time   `json:"created_at"`
}
