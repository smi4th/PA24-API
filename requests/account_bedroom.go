package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func AccountBedroom(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		AccountBedroomPost(w, r, db)
	case "GET":
		AccountBedroomGet(w, r, db)
	case "PUT":
		AccountBedroomPut(w, r, db)
	case "DELETE":
		AccountBedroomDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func AccountBedroomPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `account`, `bedroom`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    account_ := tools.BodyValueToString(body, "account")
	bedroom_ := tools.BodyValueToString(body, "bedroom")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty() {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, ) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, ) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

    if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(bedroom_) {
		if !tools.ElementExists(db, "BED_ROOM", "uuid", bedroom_) {
			tools.JsonResponse(w, 400, `{"error": "This bedroom does not exist"}`) 
			return
		}
	}
	
	
	

	

	

	

	// Inserting the AccountBedroom in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `ACCOUNT_BEDROOM` (`account`, `bedroom`, `account`, `bedroom`) VALUES (?, ?)", account_, bedroom_, account_, bedroom_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "AccountBedroom created"`

	// Adding the return fields of the query
	fields, err := AccountBedroomGetAll(db, account_, bedroom_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func AccountBedroomGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `creation_date`, `account`, `bedroom`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `creation_date`, `account`, `bedroom` FROM `ACCOUNT_BEDROOM`"
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
	jsonResponse, err := AccountBedroomGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func AccountBedroomPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, ``) || tools.ValuesNotInQuery(query, `account`, `bedroom`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]
	bedroom_ := query["bedroom"]
	
    

	// Checking if the values are empty
	if tools.ValueIsEmpty(account_, bedroom_) {
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
	if tools.ValueTooShort(4, ) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, ) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

    if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(bedroom_) {
		if !tools.ElementExists(db, "BED_ROOM", "uuid", bedroom_) {
			tools.JsonResponse(w, 400, `{"error": "This bedroom does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExists(db, "ACCOUNT_BEDROOM", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountBedroom does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "ACCOUNT_BEDROOM", "bedroom", bedroom_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountBedroom does not exist"}`) 
		return
	}
	

	

    

	request := "UPDATE `ACCOUNT_BEDROOM` SET "
	var params []interface{}
	
	for key, value := range body {
		if !tools.ValueInArray(key, `account`, `bedroom`) {
			if key == "password" {
				value = tools.HashPassword(value.(string))
			}
			tools.AppendUpdate(&request, &params, key, value)
		}
	}

	// Removing the last ","
	request = request[:len(request)-2]

	request += " WHERE account = ?, bedroom = ?"
	params = append(params, account_, bedroom_)

	// Updating the account in the database
	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "AccountBedroom updated"`
	
	// Adding the return fields of the query
	fields, err := AccountBedroomGetAll(db, account_, bedroom_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func AccountBedroomDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `account`, `bedroom`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]
	bedroom_ := query["bedroom"]
	

	if !tools.ElementExists(db, "ACCOUNT_BEDROOM", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountBedroom does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "ACCOUNT_BEDROOM", "bedroom", bedroom_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountBedroom does not exist"}`) 
		return
	}
	

	// Deleting the AccountBedroom in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `ACCOUNT_BEDROOM` WHERE account = ?, bedroom = ?", account_, bedroom_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "AccountBedroom deleted", "account": "` + account_ + `", "bedroom": "` + bedroom_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func AccountBedroomGetAll(db *sql.DB, account_ string, bedroom_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `creation_date`, `account`, `bedroom` FROM `ACCOUNT_BEDROOM` WHERE account = ?, bedroom = ?", account_, bedroom_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return AccountBedroomGetAllAssociation(result, arrayOutput)
}

func AccountBedroomGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var creation_date_, account_, bedroom_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&creation_date_, &account_, &bedroom_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"creation_date": "` + creation_date_ + `", "account": "` + account_ + `", "bedroom": "` + bedroom_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&creation_date_, &account_, &bedroom_)
			if err != nil {
				return "", err
			}
		}
		return `"creation_date": "` + creation_date_ + `", "account": "` + account_ + `", "bedroom": "` + bedroom_ + `"`, nil
	}
}