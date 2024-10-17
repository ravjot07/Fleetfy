package user

import (
	"context"
	"encoding/json"
	
	"fleetfy/backend/db"
	"fleetfy/backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type UserResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	User    models.User `json:"user,omitempty"`
}

// GetUserInfoById retrieves user information by UID
func GetUserInfoById(w http.ResponseWriter, r *http.Request) {
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	uid := params["uid"]
	if uid == "" {
		http.Error(w, "Missing UID parameter", http.StatusBadRequest)
		return
	}

	// Convert UID to int64
	uidInt, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		http.Error(w, "Invalid UID format", http.StatusBadRequest)
		return
	}

	pool := db.GetPool()

	// Query to fetch user by ID
	query := `SELECT id, name, address, phone_number, user_type, created_at 
	          FROM users WHERE id = $1`

	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.QueryRow(ctx, query, uidInt).Scan(
		&user.ID,
		&user.Name,
		&user.Address,
		&user.PhoneNumber,
		&user.UserType,
		&user.CreatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(UserResponse{Success: false, Message: "User not found"})
			return
		}
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	// Return the user details
	json.NewEncoder(w).Encode(UserResponse{
		Success: true,
		Message: "User retrieved successfully",
		User:    user,
	})
}
