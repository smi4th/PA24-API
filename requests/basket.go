package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Basket(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "GET":
		BasketGet(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func BasketGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {

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

	// Getting the services from the database
	services_, err := tools.ExecuteQuery(db, "SELECT `services` FROM `CONSUME` WHERE `account` = ?", account_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer services_.Close()

	reservation_bedroom_, err := tools.ExecuteQuery(db, "SELECT `bed_room` FROM `reservation_bedroom` WHERE `account` = ?", account_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer reservation_bedroom_.Close()

	reservation_housing_, err := tools.ExecuteQuery(db, "SELECT `housing` FROM `reservation_housing` WHERE `account` = ?", account_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer reservation_housing_.Close()

	// Creating the response
	jsonResponse := `{"services": [ `
	for services_.Next() {
		var services string
		services_.Scan(&services)
		jsonResponse += services + `,`
	}

	jsonResponse = jsonResponse[:len(jsonResponse)-1] + `], "reservation_bedroom": [ `
	for reservation_bedroom_.Next() {
		var bedroom string
		reservation_bedroom_.Scan(&bedroom)
		jsonResponse += bedroom + `,`
	}

	jsonResponse = jsonResponse[:len(jsonResponse)-1] + `], "reservation_housing": [ `
	for reservation_housing_.Next() {
		var housing string
		reservation_housing_.Scan(&housing)
		jsonResponse += housing + `,`
	}

	jsonResponse = jsonResponse[:len(jsonResponse)-1] + `]}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}