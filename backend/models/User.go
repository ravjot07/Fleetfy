package models

import "time"

type User struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Address     string    `json:"address"`
    PhoneNumber string    `json:"phone_number"`
    Password    string    `json:"password"`
    UserType    string    `json:"user_type"`
    CreatedAt   time.Time `json:"created_at"`
}
