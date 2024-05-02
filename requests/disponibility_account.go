package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func DisponibilityAccount(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		DisponibilityAccountPost(w, r, db)
	case "GET":
		DisponibilityAccountGet(w, r, db)
	case "PUT":
		DisponibilityAccountPut(w, r, db)
	case "DELETE":
		DisponibilityAccountDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func DisponibilityAccountPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `disponibility`, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    disponibility_ := tools.BodyValueToString(body, "disponibility")
	account_ := tools.BodyValueToString(body, "account")
	

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

    if !tools.ValueIsEmpty(disponibility_) {
		if !tools.ElementExists(db, "DISPONIBILITY", "uuid", disponibility_) {
			tools.JsonResponse(w, 400, `{"error": "This disponibility does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	
	
	

	

	

	

	// Inserting the DisponibilityAccount in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `DISPONIBILITY_ACCOUNT` (`disponibility`, `account`, `disponibility`, `account`) VALUES (?, ?)", disponibility_, account_, disponibility_, account_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "DisponibilityAccount created"`

	// Adding the return fields of the query
	fields, err := DisponibilityAccountGetAll(db, disponibility_, account_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func DisponibilityAccountGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `disponibility`, `account`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `disponibility`, `account` FROM `DISPONIBILITY_ACCOUNT`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `DISPONIBILITY_ACCOUNT`"
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
	jsonResponse, err := DisponibilityAccountGetAllAssociation(result, true)
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

func DisponibilityAccountPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, ``) || tools.ValuesNotInQuery(query, `disponibility`, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	disponibility_ := query["disponibility"]
	account_ := query["account"]
	
    

	// Checking if the values are empty
	if tools.ValueIsEmpty(disponibility_, account_) {
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

    if !tools.ValueIsEmpty(disponibility_) {
		if !tools.ElementExists(db, "DISPONIBILITY", "uuid", disponibility_) {
			tools.JsonResponse(w, 400, `{"error": "This disponibility does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExists(db, "DISPONIBILITY_ACCOUNT", "disponibility", disponibility_) {
		tools.JsonResponse(w, 400, `{"error": "This DisponibilityAccount does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "DISPONIBILITY_ACCOUNT", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This DisponibilityAccount does not exist"}`) 
		return
	}
	

	

    

	request := "UPDATE `DISPONIBILITY_ACCOUNT` SET "
	var params []interface{}
	
	for key, value := range body {
		if !tools.ValueInArray(key, `disponibility`, `account`) {
			if key == "password" {
				value = tools.HashPassword(value.(string))
			}
			tools.AppendUpdate(&request, &params, key, value)
		}
	}

	// Removing the last ","
	request = request[:len(request)-2]

	request += " WHERE disponibility = ?, account = ?"
	params = append(params, disponibility_, account_)

	// Updating the account in the database
	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "DisponibilityAccount updated"`
	
	// Adding the return fields of the query
	fields, err := DisponibilityAccountGetAll(db, disponibility_, account_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func DisponibilityAccountDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `disponibility`, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	disponibility_ := query["disponibility"]
	account_ := query["account"]
	

	if !tools.ElementExists(db, "DISPONIBILITY_ACCOUNT", "disponibility", disponibility_) {
		tools.JsonResponse(w, 400, `{"error": "This DisponibilityAccount does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "DISPONIBILITY_ACCOUNT", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This DisponibilityAccount does not exist"}`) 
		return
	}
	

	// Deleting the DisponibilityAccount in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `DISPONIBILITY_ACCOUNT` WHERE disponibility = ?, account = ?", disponibility_, account_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "DisponibilityAccount deleted", "disponibility": "` + disponibility_ + `", "account": "` + account_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func DisponibilityAccountGetAll(db *sql.DB, disponibility_ string, account_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `disponibility`, `account` FROM `DISPONIBILITY_ACCOUNT` WHERE disponibility = ?, account = ?", disponibility_, account_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return DisponibilityAccountGetAllAssociation(result, arrayOutput)
}

func DisponibilityAccountGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var disponibility_, account_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&disponibility_, &account_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"disponibility": "` + disponibility_ + `", "account": "` + account_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&disponibility_, &account_)
			if err != nil {
				return "", err
			}
		}
		return `"disponibility": "` + disponibility_ + `", "account": "` + account_ + `"`, nil
	}
}