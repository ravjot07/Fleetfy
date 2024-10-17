package handler

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "fmc/models"
    "log"
)

func RegisterHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username string `json:"username"`
            Password string `json:"password"`
            Role     string `json:"role"`
        }

        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        // Validate input
        if req.Username == "" || req.Password == "" || req.Role == "" {
            http.Error(w, "All fields are required", http.StatusBadRequest)
            return
        }

        // Register the user in the database
        err = models.RegisterUser(db, req.Username, req.Password, req.Role)
        if err != nil {
            // Log and return error if registration failed
            log.Printf("Error registering user: %v", err)
            http.Error(w, "Error registering user", http.StatusInternalServerError)
            return
        }

        // Success
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
    }
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }

        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        user, err := models.AuthenticateUser(db, req.Username, req.Password)
        if err != nil {
            http.Error(w, "Invalid credentials", http.StatusUnauthorized)
            return
        }

        json.NewEncoder(w).Encode(map[string]interface{}{
            "message":  "Login successful",
            "username": user.Username,
            "userID":   user.ID,
            "role":     user.Role,
        })
    }
}
