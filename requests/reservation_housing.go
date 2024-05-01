package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func ReservationHousing(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		ReservationHousingPost(w, r, db)
	case "GET":
		ReservationHousingGet(w, r, db)
	case "PUT":
		ReservationHousingPut(w, r, db)
	case "DELETE":
		ReservationHousingDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func ReservationHousingPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `start_time`, `end_time`, `price`, `review`, `review_note`, `account`, `housing`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    start_time_ := tools.BodyValueToString(body, "start_time")
	end_time_ := tools.BodyValueToString(body, "end_time")
	price_ := tools.BodyValueToString(body, "price")
	review_ := tools.BodyValueToString(body, "review")
	review_note_ := tools.BodyValueToString(body, "review_note")
	account_ := tools.BodyValueToString(body, "account")
	housing_ := tools.BodyValueToString(body, "housing")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(start_time_, end_time_, price_, review_, review_note_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, start_time_, end_time_, price_, review_, review_note_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, start_time_, end_time_, price_, review_, review_note_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

    if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(housing_) {
		if !tools.ElementExists(db, "HOUSING", "uuid", housing_) {
			tools.JsonResponse(w, 400, `{"error": "This housing does not exist"}`) 
			return
		}
	}
	
	
	

	

	

	

	// Inserting the ReservationHousing in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `RESERVATION_HOUSING` (`account`, `housing`, `start_time`, `end_time`, `price`, `review`, `review_note`, `account`, `housing`) VALUES (?, ?, ?, ?, ?, ?, ?)", account_, housing_, start_time_, end_time_, price_, review_, review_note_, account_, housing_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "ReservationHousing created"`

	// Adding the return fields of the query
	fields, err := ReservationHousingGetAll(db, account_, housing_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func ReservationHousingGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `start_time`, `end_time`, `price`, `review`, `review_note`, `account`, `housing`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `start_time`, `end_time`, `price`, `review`, `review_note`, `account`, `housing` FROM `RESERVATION_HOUSING`"
	var params []interface{}

	if query["all"] != "true" {
		request += " WHERE "
		strictSearch := query["strictSearch"] == "true"

		for key, value := range query {
			tools.AppendCondition(&request, &params, key, value, strictSearch)
		}

		// Removing the last "AND"
		request = request[:len(request)-3]
	}

	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse, err := ReservationHousingGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ReservationHousingPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `start_time`, `end_time`, `price`, `review`, `review_note`) || tools.ValuesNotInQuery(query, `account`, `housing`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]
	housing_ := query["housing"]
	
    start_time_ := tools.BodyValueToString(body, "start_time")
	end_time_ := tools.BodyValueToString(body, "end_time")
	price_ := tools.BodyValueToString(body, "price")
	review_ := tools.BodyValueToString(body, "review")
	review_note_ := tools.BodyValueToString(body, "review_note")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(account_, housing_) {
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
	if tools.ValueTooShort(4, start_time_, end_time_, price_, review_, review_note_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, start_time_, end_time_, price_, review_, review_note_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

    if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(housing_) {
		if !tools.ElementExists(db, "HOUSING", "uuid", housing_) {
			tools.JsonResponse(w, 400, `{"error": "This housing does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExists(db, "RESERVATION_HOUSING", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This ReservationHousing does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "RESERVATION_HOUSING", "housing", housing_) {
		tools.JsonResponse(w, 400, `{"error": "This ReservationHousing does not exist"}`) 
		return
	}
	

	

    

	request := "UPDATE `RESERVATION_HOUSING` SET "
	var params []interface{}
	
	for key, value := range body {
		if !tools.ValueInArray(key, `account`, `housing`) {
			if key == "password" {
				value = tools.HashPassword(value.(string))
			}
			tools.AppendUpdate(&request, &params, key, value)
		}
	}

	// Removing the last ","
	request = request[:len(request)-2]

	request += " WHERE account = ?, housing = ?"
	params = append(params, account_, housing_)

	// Updating the account in the database
	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "ReservationHousing updated"`
	
	// Adding the return fields of the query
	fields, err := ReservationHousingGetAll(db, account_, housing_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func ReservationHousingDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `account`, `housing`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]
	housing_ := query["housing"]
	

	if !tools.ElementExists(db, "RESERVATION_HOUSING", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This ReservationHousing does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "RESERVATION_HOUSING", "housing", housing_) {
		tools.JsonResponse(w, 400, `{"error": "This ReservationHousing does not exist"}`) 
		return
	}
	

	// Deleting the ReservationHousing in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `RESERVATION_HOUSING` WHERE account = ?, housing = ?", account_, housing_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "ReservationHousing deleted", "account": "` + account_ + `", "housing": "` + housing_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ReservationHousingGetAll(db *sql.DB, account_ string, housing_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `start_time`, `end_time`, `price`, `review`, `review_note`, `account`, `housing` FROM `RESERVATION_HOUSING` WHERE account = ?, housing = ?", account_, housing_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return ReservationHousingGetAllAssociation(result, arrayOutput)
}

func ReservationHousingGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var start_time_, end_time_, price_, review_, review_note_, account_, housing_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&start_time_, &end_time_, &price_, &review_, &review_note_, &account_, &housing_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"start_time": "` + start_time_ + `", "end_time": "` + end_time_ + `", "price": "` + price_ + `", "review": "` + review_ + `", "review_note": "` + review_note_ + `", "account": "` + account_ + `", "housing": "` + housing_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&start_time_, &end_time_, &price_, &review_, &review_note_, &account_, &housing_)
			if err != nil {
				return "", err
			}
		}
		return `"start_time": "` + start_time_ + `", "end_time": "` + end_time_ + `", "price": "` + price_ + `", "review": "` + review_ + `", "review_note": "` + review_note_ + `", "account": "` + account_ + `", "housing": "` + housing_ + `"`, nil
	}
}