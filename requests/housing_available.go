package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func HousingAvailable(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "GET":
		HousingAvailableGet(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func HousingAvailableGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `start_time`, `end_time`, `housing`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	start_time_ := query[`start_time`]
	end_time_ := query[`end_time`]
	housing_ := query[`housing`]

	// Checking if the values are empty
	if tools.ValueIsEmpty(start_time_, end_time_, housing_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the house exists in the database
	if !tools.ElementExists(db, "HOUSING", "uuid", housing_) {
		tools.JsonResponse(w, 404, `{"message": "House not found"}`)
		return
	}


	result := tools.PeriodeOverlap(db, "RESERVATION_HOUSING", "start_time", "end_time", "housing", housing_, start_time_, end_time_)

	// Checking if the house is available
	jsonResponse := `{"available": true}`
	if result {
		jsonResponse = `{"available": false}`
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}