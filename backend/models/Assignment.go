package models

import "time"

type Assignment struct {
    ID          int64     `json:"id"`
    UID         int64     `json:"uid"`
    VehicleNo   string    `json:"vehicle_no"`
    BookingID   int64     `json:"booking_id"`
    CreatedAt   time.Time `json:"created_at"`
}
