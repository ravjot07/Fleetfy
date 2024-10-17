package authentication

import (
	"context"
	"encoding/json"

	"fleetfy/backend/db"
	"fleetfy/backend/models"
	"net/http"
	"time"


)

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	UserType    string `json:"type"`
	Success     bool   `json:"success"`
	Message     string `json:"message"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	pool := db.GetPool()

	// Query to find user by phone number and password
	query := `SELECT id, name, address, phone_number, user_type FROM users 
	          WHERE phone_number = $1 AND password = $2`

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.QueryRow(ctx, query, loginReq.PhoneNumber, loginReq.Password).
		Scan(&user.ID, &user.Name, &user.Address, &user.PhoneNumber, &user.UserType)
	if err != nil {
		http.Error(w, "Invalid phone number or password", http.StatusUnauthorized)
		return
	}

	// User authenticated successfully
	loginRes := LoginResponse{
		ID:          user.ID,
		Name:        user.Name,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		UserType:    user.UserType,
		Success:     true,
		Message:     "Authentication successful",
	}

	json.NewEncoder(w).Encode(loginRes)
}
