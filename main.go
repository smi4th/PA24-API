package main

import (
	"net/http"
	"requests"
	"tools"
)

func main() {

	// Initialize the database connection
	db := tools.InitDatabaseConnection()
	defer tools.CloseDatabaseConnection(db)

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		tools.JsonResponse(w, 200, `{"message": "Hello, API!"}`)
	})

	http.HandleFunc("/api/account", func(w http.ResponseWriter, r *http.Request) {
		requests.Account(w, r)
	})

	tools.InfoLog("Server is running on port 80")
	http.ListenAndServe(":80", nil)

}