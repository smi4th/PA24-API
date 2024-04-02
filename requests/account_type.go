package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func AccountType(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		AccountTypePost(w, r, db)
	case "GET":
		AccountTypeGet(w, r, db)
	case "PUT":
		AccountTypePut(w, r, db)
	case "DELETE":
		AccountTypeDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func AccountTypePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, "type") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	type_ := tools.BodyValueToString(body, "type")

	// Checking if the values are empty
	if tools.ValueIsEmpty(type_) {
		tools.JsonResponse(w, 400, `{"message": "Type cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, type_) {
		tools.JsonResponse(w, 400, `{"message": "Type too short"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooLong(32, type_) {
		tools.JsonResponse(w, 400, `{"message": "Type too long"}`)
		return
	}

	// Checking if the account type is valid
	if tools.ElementExists(db, "ACCOUNT_TYPE", "type", type_) {
		tools.JsonResponse(w, 400, `{"message": "Account type already exists"}`)
		return
	}

	uuid := tools.GenerateUUID()

	// Inserting the account in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `ACCOUNT_TYPE` (`id`, `type`) VALUES (?, ?)", uuid, type_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Account type created"`

	// Adding the return fields of the query
	fields, err := AccountTypeGetAll(db, uuid, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	jsonResponse += "," + fields

	tools.InfoLog(tools.RowsToJson(result))

	jsonResponse += "}"

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse)

}

func AccountTypeGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, "id", "type") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `id`, `type` FROM `ACCOUNT_TYPE` WHERE "
	var params []interface{}
	strictSearch := query["strictSearch"] == "true"

	for key, value := range query {
		tools.AppendCondition(&request, &params, key, value, strictSearch)
	}

	// Removing the last "AND"
	request = request[:len(request)-3]

	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse, err := AccountTypeGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func AccountTypePut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, "type") || tools.ValuesNotInQuery(query, "id") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	id := query["id"]
	type_ := tools.BodyValueToString(body, "type")

	// Checking if the values are empty
	if tools.ValueIsEmpty(id, type_) {
		tools.JsonResponse(w, 400, `{"message": "ID and type cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, type_) {
		tools.JsonResponse(w, 400, `{"message": "type too short"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooLong(32, type_) {
		tools.JsonResponse(w, 400, `{"message": "type too long"}`)
		return
	}

	// Checking if the account exists
	if !tools.ElementExists(db, "ACCOUNT_TYPE", "id", id) {
		tools.JsonResponse(w, 400, `{"message": "Account type does not exist"}`)
		return
	}

	// Checking if the account type is valid
	if tools.ElementExists(db, "ACCOUNT_TYPE", "type", type_) {
		tools.JsonResponse(w, 400, `{"message": "Account type already exists"}`)
		return
	}

	// Updating the account in the database
	result, err := tools.ExecuteQuery(db, "UPDATE `ACCOUNT_TYPE` SET `type` = ? WHERE `id` = ?", type_, id)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Account type updated"`

	// Adding the return fields of the query
	fields, err := AccountTypeGetAll(db, id, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	jsonResponse += "," + fields + "}"

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func AccountTypeDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, "id") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	id := query["id"]

	// Checking if the account exists
	if !tools.ElementExists(db, "ACCOUNT_TYPE", "id", id) {
		tools.JsonResponse(w, 400, `{"message": "Account type does not exist"}`)
		return
	}

	// Deleting the account in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `ACCOUNT_TYPE` WHERE `id` = ?", id)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Account type deleted", "id": "` + id + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func AccountTypeGetAll(db *sql.DB, uuid string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `id`, `type` FROM `ACCOUNT_TYPE` WHERE `id` = ?", uuid)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return AccountTypeGetAllAssociation(result, arrayOutput)
}

func AccountTypeGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var id string
	var type_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&id, &type_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"id": "` + id + `", "type": "` + type_ + `"},`
		}
		jsonResponse = jsonResponse[:len(jsonResponse)-1]
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&id, &type_)
			if err != nil {
				return "", err
			}
		}
		return `"id": "` + id + `", "type": "` + type_ + `"`, nil
	}
}