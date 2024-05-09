package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func HousingReservation(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "GET":
		HousingReservationGet(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func HousingReservationGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `housing`, `all`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}
	housing_ := query[`housing`]

	// Checking if the house exists in the database
	if !tools.ValueIsEmpty(housing_) {
		if !tools.ElementExists(db, "HOUSING", "uuid", housing_) {
			tools.JsonResponse(w, 404, `{"message": "housing not found"}`)
			return
		}
	}

	request := `SELECT start_time, end_time FROM BASKET_HOUSING WHERE `

	if !tools.ValueIsEmpty(housing_) {
		request += `housing = ?`
	} else {
		request += `1 = 1 OR 1 = ?`
	}

	rows, err := tools.ExecuteQuery(db, request, housing_)
	if err != nil {
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	var response string
	response = `[`
	for rows.Next() {
		var start_time, end_time string
		rows.Scan(&start_time, &end_time)
		response += `{"start_time": "` + start_time + `", "end_time": "` + end_time + `"},`
	}
	response = response[:len(response)-1]
	response += `]`
	tools.JsonResponse(w, 200, response)
}