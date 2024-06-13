package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func VerifyPassword(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "GET":
		VerifyPasswordGet(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func VerifyPasswordGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `password`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid:= tools.GetUUID(r, db)

	// Checking if the account exists
	if !tools.ElementExists(db, "ACCOUNT", "uuid", uuid) {
		tools.JsonResponse(w, 404, `{"message": "User not found"}`)
		return
	}

	// Getting the account's password
	password := tools.GetElement(db, "ACCOUNT", "password", "uuid", uuid)

	// Checking if the password is correct
	if tools.ComparePassword(password, query["password"]) {
		tools.JsonResponse(w, 200, `{"correct": true}`)
		return
	}

	tools.JsonResponse(w, 200, `{"correct": false}`)

}