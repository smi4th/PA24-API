package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Admin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "GET":
		AdminGet(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func AdminGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query[`account`]

	// Checking if the values are empty
	if tools.ValueIsEmpty(account_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the account exists in the database
	if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
		tools.JsonResponse(w, 404, `{"message": "Account not found"}`)
		return
	}

	// Getting the account_type from the database
	account_type_ := tools.GetElement(db, "ACCOUNT_TYPE", "admin", "uuid", tools.GetElement(db, "ACCOUNT", "account_type", "uuid", account_))

	jsonResponse := ""
	// Creating the response
	if account_type_ == "true" {
		jsonResponse = `{"admin": true}`
	} else {
		jsonResponse = `{"admin": false}`
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}