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

	// Handle the requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/account":
			requests.Account(w, r, db)
		case "/api/account_type":
			requests.AccountType(w, r, db)
		case "/api/provider":
			requests.Provider(w, r, db)
		case "/api/provider_account":
			requests.ProviderAccount(w, r, db)
		default:
			tools.JsonResponse(w, 404, `{"message": "Not found"}`)
		}
	})

	tools.InfoLog("Server is running on port 80")
	http.ListenAndServe(":80", nil)

}