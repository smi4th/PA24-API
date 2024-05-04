package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func ReservationBedroom(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		ReservationBedroomPost(w, r, db)
	case "GET":
		ReservationBedroomGet(w, r, db)
	case "PUT":
		if tools.GetUUID(r, db) == tools.GetElementFromLinkTable(db, "RESERVATION_BEDROOM", "account", "bed_room", tools.ReadQuery(r)["bed_room"], "account", tools.ReadQuery(r)["account"]) {
			ReservationBedroomPut(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	case "DELETE":
		if tools.GetUUID(r, db) == tools.GetElementFromLinkTable(db, "RESERVATION_BEDROOM", "account", "bed_room", tools.ReadQuery(r)["bed_room"], "account", tools.ReadQuery(r)["account"]) {
			ReservationBedroomDelete(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func ReservationBedroomPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `start_time`, `end_time`, `account`, `bed_room`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    start_time_ := tools.BodyValueToString(body, "start_time")
	end_time_ := tools.BodyValueToString(body, "end_time")
	review_ := "None"
	review_note_ := "0"
	account_ := tools.BodyValueToString(body, "account")
	bed_room_ := tools.BodyValueToString(body, "bed_room")

	if tools.GetUUID(r, db) != account_ {
		tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		return
	}

	if tools.GetElement(db, "HOUSING", "account", "uuid", tools.GetElement(db, "BED_ROOM", "housing", "uuid", bed_room_)) == tools.GetUUID(r, db) {
		tools.JsonResponse(w, 404, `{"message": "Cannot reserve your own bed_room"}`)
		return
	}
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(start_time_, end_time_, review_, review_note_, account_, bed_room_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, review_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}

	if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
		tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "BED_ROOM", "uuid", bed_room_) {
		tools.JsonResponse(w, 400, `{"error": "This bed_room does not exist"}`) 
		return
	}

	if tools.ElementExistsInLinkTable(db, "RESERVATION_BEDROOM", "account", account_, "bed_room", bed_room_) {
		tools.JsonResponse(w, 400, `{"error": "This ReservationBedroom already exists"}`)
		return
	}
	
	

	

	

	

	// Inserting the ReservationBedroom in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `RESERVATION_BEDROOM` (`account`, `bed_room`, `start_time`, `end_time`, `review`, `review_note`) VALUES (?, ?, ?, ?, ?, ?)", account_, bed_room_, start_time_, end_time_, review_, review_note_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "ReservationBedroom created"`

	// Adding the return fields of the query
	fields, err := ReservationBedroomGetAll(db, account_, bed_room_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func ReservationBedroomGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `start_time`, `end_time`, `review`, `review_note`, `account`, `bed_room`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `start_time`, `end_time`, `review`, `review_note`, `account`, `bed_room` FROM `RESERVATION_BEDROOM`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `RESERVATION_BEDROOM`"
	var countParams []interface{}

	if query["all"] != "true" {
		request += " WHERE "
		countRequest += " WHERE "
		strictSearch := query["strictSearch"] == "true"

		for key, value := range query {
			tools.AppendCondition(&request, &params, key, value, strictSearch)
			tools.AppendCondition(&countRequest, &countParams, key, value, strictSearch)
		}

		// Removing the last "AND"
		request = request[:len(request)-3]
		countRequest = countRequest[:len(countRequest)-3]
	}

	if query["limit"] != "" {
		request += " LIMIT " + query["limit"]

		if query["offset"] != "" {
			request += " OFFSET " + query["offset"]
		}
	}

	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse, err := ReservationBedroomGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	
	result, err = tools.ExecuteQuery(db, countRequest, countParams...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	var count string
	for result.Next() {
		err := result.Scan(&count)
		if err != nil {
			tools.ErrorLog(err.Error())
			tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
			return
		}
	}

	// Sending the response
	tools.JsonResponse(w, 200, `{"total": ` + count + `, "data": ` + jsonResponse + `}`)

}

func ReservationBedroomPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `start_time`, `end_time`, `review`, `review_note`) || tools.ValuesNotInQuery(query, `account`, `bed_room`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]
	bed_room_ := query["bed_room"]
	
	review_ := tools.BodyValueToString(body, "review")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(account_, bed_room_) {
		tools.JsonResponse(w, 400, `{"message": "Empty fields"}`)
		return
	}

	// for each key in the body, if the key is not in the query, return an error
	for key, _ := range body {
		// if the key is empty
		if tools.ValueIsEmpty(tools.BodyValueToString(body, key)) {
			tools.JsonResponse(w, 400, `{"message": "Empty fields"}`)
			return
		}
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, review_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}

    if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(bed_room_) {
		if !tools.ElementExists(db, "BED_ROOM", "uuid", bed_room_) {
			tools.JsonResponse(w, 400, `{"error": "This bed_room does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExistsInLinkTable(db, "RESERVATION_BEDROOM", "account", account_, "bed_room", bed_room_) {
		tools.JsonResponse(w, 400, `{"error": "This ReservationBedroom does not exist"}`) 
		return
	}
	

	

    

	request := "UPDATE `RESERVATION_BEDROOM` SET "
	var params []interface{}
	
	for key, value := range body {
		if !tools.ValueInArray(key, `account`, `bed_room`) {
			if key == "password" {
				value = tools.HashPassword(value.(string))
			}
			tools.AppendUpdate(&request, &params, key, value)
		}
	}

	// Removing the last ","
	request = request[:len(request)-2]

	request += " WHERE account = ?, bed_room = ?"
	params = append(params, account_, bed_room_)

	// Updating the account in the database
	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "ReservationBedroom updated"`
	
	// Adding the return fields of the query
	fields, err := ReservationBedroomGetAll(db, account_, bed_room_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func ReservationBedroomDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `account`, `bed_room`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]
	bed_room_ := query["bed_room"]
	

	if !tools.ElementExists(db, "RESERVATION_BEDROOM", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`)
		return
	}

	if !tools.ElementExists(db, "RESERVATION_BEDROOM", "bed_room", bed_room_) {
		tools.JsonResponse(w, 400, `{"error": "This bed_room does not exist"}`)
		return
	}

	if !tools.ElementExistsInLinkTable(db, "RESERVATION_BEDROOM", "account", account_, "bed_room", bed_room_) {
		tools.JsonResponse(w, 400, `{"error": "This ReservationBedroom does not exist"}`)
		return
	}
	

	// Deleting the ReservationBedroom in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `RESERVATION_BEDROOM` WHERE account = ? AND bed_room = ?", account_, bed_room_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "ReservationBedroom deleted", "account": "` + account_ + `", "bed_room": "` + bed_room_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ReservationBedroomGetAll(db *sql.DB, account_ string, bed_room_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `start_time`, `end_time`, `review`, `review_note`, `account`, `bed_room` FROM `RESERVATION_BEDROOM` WHERE account = ? AND bed_room = ?", account_, bed_room_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return ReservationBedroomGetAllAssociation(result, arrayOutput)
}

func ReservationBedroomGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var start_time_, end_time_, review_, review_note_, account_, bed_room_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&start_time_, &end_time_, &review_, &review_note_, &account_, &bed_room_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"start_time": "` + start_time_ + `", "end_time": "` + end_time_ + `", "price": "` + review_ + `", "review_note": "` + review_note_ + `", "account": "` + account_ + `", "bed_room": "` + bed_room_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&start_time_, &end_time_, &review_, &review_note_, &account_, &bed_room_)
			if err != nil {
				return "", err
			}
		}
		return `"start_time": "` + start_time_ + `", "end_time": "` + end_time_ + `", "price": "` + review_ + `", "review_note": "` + review_note_ + `", "account": "` + account_ + `", "bed_room": "` + bed_room_ + `"`, nil
	}
}