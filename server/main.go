package main

import (
    "fmc/database"
    "fmc/handler"
    "fmc/middleware"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"
)

func main() {
    // Initialize the database
    database.InitDB()
    db := database.DB

    // Initialize the router
    r := mux.NewRouter()

    // Apply global CORS middleware
    r.Use(middleware.CORS)

    // Authentication routes (User registration and login)
    r.HandleFunc("/register", handler.RegisterHandler(db)).Methods("POST")
    r.HandleFunc("/login", handler.LoginHandler(db)).Methods("POST")

    // Protected routes for Admin
    adminRouter := r.PathPrefix("/admin").Subrouter()
    adminRouter.Use(middleware.RoleMiddleware("admin"))  // Protect with admin role middleware
    adminRouter.HandleFunc("/getVehicles", handler.GetAllVehiclesHandler(db)).Methods("GET")  // Admin gets all vehicles
	adminRouter.HandleFunc("/vehicles", handler.CreateVehicleHandler(db)).Methods("POST")  // Admin creates a vehicle
    adminRouter.HandleFunc("/bookings", handler.GetAllBookingsHandler(db)).Methods("GET")  // Get all bookings
    adminRouter.HandleFunc("/bookings/{id}/complete", handler.CompleteBookingHandler(db)).Methods("PUT")  // Mark a booking as complete
    adminRouter.HandleFunc("/drivers/active-bookings", handler.GetDriverActiveBookingsCount(db)).Methods("GET")  // Get active bookings count per driver

    // adminRouter.HandleFunc("/bookings/pending", handler.GetPendingBookingsHandler(db)).Methods("GET")
    // adminRouter.Use(middleware.RoleMiddleware("admin", "driver"))  // Allow both 'admin' and 'driver'


    // Routes for Users to book vehicles
    userRouter := r.PathPrefix("/user").Subrouter()
    userRouter.Use(middleware.RoleMiddleware("user"))  // Protect with user role middleware
    userRouter.HandleFunc("/bookings", handler.CreateBookingHandler(db)).Methods("POST")  // Create a booking

    // Routes for Drivers to accept bookings
    driverRouter := r.PathPrefix("/driver").Subrouter()
    driverRouter.Use(middleware.RoleMiddleware("driver"))  // Protect with driver role middleware
    driverRouter.HandleFunc("/bookings/pending", handler.GetPendingBookingsHandler(db)).Methods("GET")
    driverRouter.HandleFunc("/bookings/{id}/accept", handler.AcceptBookingHandler(db)).Methods("PUT")  // Driver accepts booking

    // Add CORS support for frontend
    headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "User-ID", "Driver-ID", "Role"})
    methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
    origins := handlers.AllowedOrigins([]string{"http://localhost:5173"}) // Allow Vite dev server requests

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server is running on port %s...", port)
    log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(origins, headers, methods)(r)))
}