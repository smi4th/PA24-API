package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func AccountServices(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		AccountServicesPost(w, r, db)
	case "GET":
		AccountServicesGet(w, r, db)
	case "PUT":
		AccountServicesPut(w, r, db)
	case "DELETE":
		AccountServicesDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func AccountServicesPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `account`, `services`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    account_ := tools.BodyValueToString(body, "account")
	services_ := tools.BodyValueToString(body, "services")
	

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
	if !tools.ValueIsEmpty(services_) {
		if !tools.ElementExists(db, "SERVICES", "uuid", services_) {
			tools.JsonResponse(w, 400, `{"error": "This services does not exist"}`) 
			return
		}
	}
	
	
	

	

	

	

	// Inserting the AccountServices in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `ACCOUNT_SERVICES` (`account`, `services`, `account`, `services`) VALUES (?, ?)", account_, services_, account_, services_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "AccountServices created"`

	// Adding the return fields of the query
	fields, err := AccountServicesGetAll(db, account_, services_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func AccountServicesGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `account`, `services`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `account`, `services` FROM `ACCOUNT_SERVICES`"
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
	jsonResponse, err := AccountServicesGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func AccountServicesPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, ``) || tools.ValuesNotInQuery(query, `account`, `services`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]
	services_ := query["services"]
	
    

	// Checking if the values are empty
	if tools.ValueIsEmpty(account_, services_) {
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
	if !tools.ValueIsEmpty(services_) {
		if !tools.ElementExists(db, "SERVICES", "uuid", services_) {
			tools.JsonResponse(w, 400, `{"error": "This services does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExists(db, "ACCOUNT_SERVICES", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountServices does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "ACCOUNT_SERVICES", "services", services_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountServices does not exist"}`) 
		return
	}
	

	

    

	request := "UPDATE `ACCOUNT_SERVICES` SET "
	var params []interface{}
	
	for key, value := range body {
		if !tools.ValueInArray(key, `account`, `services`) {
			if key == "password" {
				value = tools.HashPassword(value.(string))
			}
			tools.AppendUpdate(&request, &params, key, value)
		}
	}

	// Removing the last ","
	request = request[:len(request)-2]

	request += " WHERE account = ?, services = ?"
	params = append(params, account_, services_)

	// Updating the account in the database
	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "AccountServices updated"`
	
	// Adding the return fields of the query
	fields, err := AccountServicesGetAll(db, account_, services_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func AccountServicesDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `account`, `services`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]
	services_ := query["services"]
	

	if !tools.ElementExists(db, "ACCOUNT_SERVICES", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountServices does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "ACCOUNT_SERVICES", "services", services_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountServices does not exist"}`) 
		return
	}
	

	// Deleting the AccountServices in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `ACCOUNT_SERVICES` WHERE account = ?, services = ?", account_, services_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "AccountServices deleted", "account": "` + account_ + `", "services": "` + services_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func AccountServicesGetAll(db *sql.DB, account_ string, services_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `account`, `services` FROM `ACCOUNT_SERVICES` WHERE account = ?, services = ?", account_, services_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return AccountServicesGetAllAssociation(result, arrayOutput)
}

func AccountServicesGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var account_, services_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&account_, &services_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"account": "` + account_ + `", "services": "` + services_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&account_, &services_)
			if err != nil {
				return "", err
			}
		}
		return `"account": "` + account_ + `", "services": "` + services_ + `"`, nil
	}
}