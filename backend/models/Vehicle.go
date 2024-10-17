package models

import "time"

type Vehicle struct {
    ID          int64       `json:"id"`
    VehicleNo   string      `json:"vehicle_no"`
    VehicleType string      `json:"vehicle_type"`
    Coordinates Coordinates `json:"coordinates"`
    Busy        bool        `json:"busy"`
    CreatedAt   time.Time   `json:"created_at"`
}
