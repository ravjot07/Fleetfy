package models

import (
    "database/sql"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID       int
    Username string
    Password string
    Role     string
}

// Register a new user
func RegisterUser(db *sql.DB, username, password, role string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    query := `INSERT INTO users (username, password, role) VALUES ($1, $2, $3)`
    _, err = db.Exec(query, username, hashedPassword, role)
    if err != nil {
        return err
    }

    return nil
}

// Authenticate a user
func AuthenticateUser(db *sql.DB, username, password string) (*User, error) {
    user := &User{}
    query := `SELECT id, username, password, role FROM users WHERE username=$1`
    err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
    if err != nil {
        return nil, err
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, err
    }

    return user, nil
}
