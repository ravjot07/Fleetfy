package main

import (
	"fmt"
	"log"
	"fleetfy/backend/db"
	"fleetfy/backend/routes"
	analytics "fleetfy/backend/routes/Analytics"
	authentication "fleetfy/backend/routes/Authentication"
	booking "fleetfy/backend/routes/Booking"
	user "fleetfy/backend/routes/User"
	vehicles "fleetfy/backend/routes/Vehicles"
	"net/http"
	"os"
	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {
	// PostgreSQL connection string
	// It's recommended to use environment variables for sensitive information
	// Example format: "postgres://username:password@localhost:5432/dbname?sslmode=disable"
	err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found. Proceeding with environment variables.")
    }

    // Retrieve DATABASE_URL from environment variables
    dbURL := getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/logi_craft?sslmode=disable")


	// Connect to PostgreSQL
	err := db.ConnectDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.CloseDB()

	// Initialize the router
	router := mux.NewRouter()
	router.HandleFunc("/", routes.GetHome).Methods("GET")

	// =====================
	// === Authentication ===
	// =====================
	router.HandleFunc("/login", authentication.LoginHandler).Methods("POST")
	router.HandleFunc("/signup", authentication.SignupHandler).Methods("POST")

	// =====================
	// ===== Bookings ======
	// =====================
	router.HandleFunc("/book", booking.HandleBooking).Methods("POST")
	router.HandleFunc("/booking/{bookingId}", booking.GetBookingByID).Methods("GET")
	router.HandleFunc("/bookings/id/user/{uid}", booking.GetBookingsByUID).Methods("GET")
	router.HandleFunc("/bookings/id/driver/{driverId}", booking.GetBookingsByDriverID).Methods("GET")
	router.HandleFunc("/bookings", booking.GetAllBookings).Methods("GET")
	router.HandleFunc("/complete-job/{bookingId}", booking.CompleteJobHandler).Methods("GET")

	// =====================
	// ===== Analytics =====
	// =====================
	router.HandleFunc("/analysis/bookings", analytics.GetBookingAnalysis).Methods("GET")
	router.HandleFunc("/analysis/vehicles", analytics.GetVehicleAnalysis).Methods("GET")
	router.HandleFunc("/analysis/drivers", analytics.GetDriverAnalysis).Methods("GET")

	// =====================
	// ======== Users =======
	// =====================
	router.HandleFunc("/users/{uid}", user.GetUserInfoById).Methods("GET")

	// =====================
	// ===== Assignments ====
	// =====================
	router.HandleFunc("/assignment-user/{uid}", booking.GetAssignmentByUid).Methods("GET")
	router.HandleFunc("/assignment-vehicle/{vehicle_no}", booking.GetAssignmentByVehicleNo).Methods("GET")
	router.HandleFunc("/assignments", booking.GetAllAssignments).Methods("GET")
	router.HandleFunc("/assignments/{uid}/assign_vehicle", booking.AssignVehicle).Methods("PUT")

	// =====================
	// ======= Vehicles =====
	// =====================
	router.HandleFunc("/vehicle/{vehicleNo}", vehicles.GetVehicleInfoById).Methods("GET")
	router.HandleFunc("/vehicles", vehicles.GetAllVehicles).Methods("GET")
	router.HandleFunc("/add/vehicle", vehicles.AddVehicleHandler).Methods("POST")
	router.HandleFunc("/vehicle/update-location/{vehicle_no}", vehicles.UpdateVehicleLocation).Methods("PUT")
	router.HandleFunc("/vehicle-coords/{vehicle_no}", vehicles.GetVehicleCoordsHandler).Methods("GET")

	// Determine the port to listen on
	port := "4001" // Default port
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	fmt.Printf("Starting server on port %s\n", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Unable to start server on port %s: %v", port, err)
	}
}

// getEnv retrieves environment variables or returns a default value if not set.
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
