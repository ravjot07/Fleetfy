package authentication

import (
	"context"
	"encoding/json"
	
	"fleetfy/backend/db"

	"net/http"
	"time"

)

type SignupRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	UserType    string `json:"user_type"`
}

type SignupResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var signupReq SignupRequest
	err := json.NewDecoder(r.Body).Decode(&signupReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	pool := db.GetPool()

	// Check if the phone number already exists
	checkQuery := `SELECT COUNT(*) FROM users WHERE phone_number = $1`
	var count int
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.QueryRow(ctx, checkQuery, signupReq.PhoneNumber).Scan(&count)
	if err != nil {
		http.Error(w, "Failed to check existing phone number", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		// Phone number already exists
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(SignupResponse{Success: false, Message: "Phone number already registered"})
		return
	}

	// Insert the new user
	insertUserQuery := `INSERT INTO users (name, address, phone_number, password, user_type, created_at) 
	                    VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var userID int64
	createdAt := time.Now()

	err = pool.QueryRow(ctx, insertUserQuery, signupReq.Name, signupReq.Address, signupReq.PhoneNumber, signupReq.Password, signupReq.UserType, createdAt).
		Scan(&userID)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// If the user type is "driver", create an entry in the assignments table
	if signupReq.UserType == "driver" {
		insertAssignmentQuery := `INSERT INTO assignments (uid, vehicle_no, booking_id, created_at) 
		                          VALUES ($1, $2, $3, $4)`
		_, err := pool.Exec(ctx, insertAssignmentQuery, userID, "", "", createdAt)
		if err != nil {
			http.Error(w, "Failed to create driver assignment", http.StatusInternalServerError)
			return
		}
	}

	// Respond with success
	json.NewEncoder(w).Encode(SignupResponse{Success: true, Message: "User created successfully"})
}
