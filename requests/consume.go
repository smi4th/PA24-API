package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Consume(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		ConsumePost(w, r, db)
	case "GET":
		ConsumeGet(w, r, db)
	case "PUT":
		ConsumePut(w, r, db)
	case "DELETE":
		ConsumeDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func ConsumePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `report`, `notice`, `price`, `note`, `services`, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    report_ := tools.BodyValueToString(body, "report")
	notice_ := tools.BodyValueToString(body, "notice")
	price_ := tools.BodyValueToString(body, "price")
	note_ := tools.BodyValueToString(body, "note")
	services_ := tools.BodyValueToString(body, "services")
	account_ := tools.BodyValueToString(body, "account")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(report_, notice_, price_, note_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, report_, notice_, price_, note_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, report_, notice_, price_, note_) {
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
	
	
	

	

	

	

	// Inserting the Consume in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `CONSUME` (`account`, `services`, `report`, `notice`, `price`, `note`, `services`, `account`) VALUES (?, ?, ?, ?, ?, ?)", account_, services_, report_, notice_, price_, note_, services_, account_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Consume created"`

	// Adding the return fields of the query
	fields, err := ConsumeGetAll(db, account_, services_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func ConsumeGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `report`, `notice`, `price`, `note`, `services`, `account`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `report`, `notice`, `price`, `note`, `services`, `account` FROM `CONSUME`"
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
	jsonResponse, err := ConsumeGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ConsumePut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `report`, `notice`, `price`, `note`) || tools.ValuesNotInQuery(query, `account`, `services`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	services_ := query["services"]
	account_ := query["account"]
	
    report_ := tools.BodyValueToString(body, "report")
	notice_ := tools.BodyValueToString(body, "notice")
	price_ := tools.BodyValueToString(body, "price")
	note_ := tools.BodyValueToString(body, "note")
	

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
	if tools.ValueTooShort(4, report_, notice_, price_, note_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, report_, notice_, price_, note_) {
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
	

	if !tools.ElementExists(db, "CONSUME", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This Consume does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "CONSUME", "services", services_) {
		tools.JsonResponse(w, 400, `{"error": "This Consume does not exist"}`) 
		return
	}
	

	{{passwordCheck}}

    {{emailCheck}}

	request := "UPDATE `CONSUME` SET "
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
	jsonResponse := `{"message": "Consume updated"`
	
	// Adding the return fields of the query
	fields, err := ConsumeGetAll(db, account_, services_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func ConsumeDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
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
	

	if !tools.ElementExists(db, "CONSUME", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This Consume does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "CONSUME", "services", services_) {
		tools.JsonResponse(w, 400, `{"error": "This Consume does not exist"}`) 
		return
	}
	

	// Deleting the Consume in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `CONSUME` WHERE account = ?, services = ?", account_, services_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Consume deleted", "account": "` + account_ + `", "services": "` + services_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ConsumeGetAll(db *sql.DB, account_ string, services_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `report`, `notice`, `price`, `note`, `services`, `account` FROM `CONSUME` WHERE account = ?, services = ?", account_, services_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return ConsumeGetAllAssociation(result, arrayOutput)
}

func ConsumeGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var report_, notice_, price_, note_, services_, account_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&report_, &notice_, &price_, &note_, &services_, &account_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"report": "` + report_ + `", "notice": "` + notice_ + `", "price": "` + price_ + `", "note": "` + note_ + `", "services": "` + services_ + `", "account": "` + account_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&report_, &notice_, &price_, &note_, &services_, &account_)
			if err != nil {
				return "", err
			}
		}
		return `"report": "` + report_ + `", "notice": "` + notice_ + `", "price": "` + price_ + `", "note": "` + note_ + `", "services": "` + services_ + `", "account": "` + account_ + `"`, nil
	}
}