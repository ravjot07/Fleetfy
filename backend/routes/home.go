package routes

import "net/http"

func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<H1>Welcome to fleetfy</H1>"))
}
