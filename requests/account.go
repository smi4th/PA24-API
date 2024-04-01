package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Account(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		AccountPost(w, r, db)
	case "GET":
		AccountGet(w, r, db)
	case "PUT":
		AccountPut(w, r, db)
	case "DELETE":
		AccountDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func AccountPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, "username", "password", "first_name", "last_name", "email", "account_type") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	username := tools.BodyValueToString(body, "username")
	password := tools.BodyValueToString(body, "password")
	firstName := tools.BodyValueToString(body, "first_name")
	lastName := tools.BodyValueToString(body, "last_name")
	email := tools.BodyValueToString(body, "email")
	accountType := tools.BodyValueToString(body, "account_type")

	// Checking if the values are empty
	if tools.ValueIsEmpty(username, password, firstName, lastName, email, accountType) {
		tools.JsonResponse(w, 400, `{"message": "Empty fields"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(8, username, password) {
		tools.JsonResponse(w, 400, `{"message": "Username or password too short"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooLong(32, username, password, firstName, lastName) {
		tools.JsonResponse(w, 400, `{"message": "Username, password, first name or last name too long"}`)
		return
	}

	if tools.ValueTooLong(64, email) {
		tools.JsonResponse(w, 400, `{"message": "Email too long"}`)
		return
	}

	// Checking if the email is valid
	if !tools.EmailIsValid(email) {
		tools.JsonResponse(w, 400, `{"message": "Invalid email"}`)
		return
	}

	// Checking if the password is strong enough
	if tools.PasswordNotStrong(password) {
		tools.JsonResponse(w, 400, `{"message": "Password not strong enough"}`)
		return
	}
	
	// Checking if the username is already taken
	if tools.ElementExists(db, "ACCOUNT", "username", username) {
		tools.JsonResponse(w, 400, `{"message": "Username already taken"}`)
		return
	}

	// Checking if the email is already taken
	if tools.ElementExists(db, "ACCOUNT", "email", email) {
		tools.JsonResponse(w, 400, `{"message": "Email already taken"}`)
		return
	}

	// Checking if the account type is valid
	if !tools.ElementExists(db, "ACCOUNT_TYPE", "id", accountType) {
		tools.JsonResponse(w, 400, `{"message": "Invalid account type"}`)
		return
	}

	// Hashing the password
	hashedPassword := tools.HashPassword(password)

	uuid := tools.GenerateUUID()

	// Inserting the account in the database
	_, err := tools.ExecuteQuery(db, "INSERT INTO `ACCOUNT` (`id`, `username`, `password`, `first_name`, `last_name`, `email`, `account_type`) VALUES (?, ?, ?, ?, ?, ?, ?)", uuid, username, hashedPassword, firstName, lastName, email, accountType)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Creating the response
	jsonResponse := `{"message": "Account created"`

	// Adding the return fields of the query
	fields, err := AccountGetAll(db, uuid, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func AccountGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, "id", "username", "first_name", "last_name", "email", "account_type", "creation_date") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `id`, `username`, `first_name`, `last_name`, `email`, `account_type`, `creation_date` FROM `ACCOUNT` WHERE "
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
	jsonResponse, err := AccountGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func AccountPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, "username", "password", "first_name", "last_name", "email", "account_type", "creation_date") || tools.ValuesNotInQuery(query, "id") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	id := query["id"]
	username := tools.BodyValueToString(body, "username")
	password := tools.BodyValueToString(body, "password")
	firstName := tools.BodyValueToString(body, "first_name")
	lastName := tools.BodyValueToString(body, "last_name")
	email := tools.BodyValueToString(body, "email")
	accountType := tools.BodyValueToString(body, "account_type")
	creationDate := tools.BodyValueToString(body, "creation_date")

	// Checking if the values are empty
	if tools.AtLeastOneValueNotEmpty(username, password, firstName, lastName, email, accountType, creationDate) || tools.ValueIsEmpty(id) {
		tools.JsonResponse(w, 400, `{"message": "Empty fields or incorrect values length"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, username, password, firstName, lastName, creationDate) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooLong(32, username, password, firstName, lastName, creationDate) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

	if tools.ValueTooLong(64, email) {
		tools.JsonResponse(w, 400, `{"message": "Email too long"}`)
		return
	}

	// Checking if the account exists
	if !tools.ElementExists(db, "ACCOUNT", "id", id) {
		tools.JsonResponse(w, 400, `{"message": "Account does not exist"}`)
		return
	}

	request := "UPDATE `ACCOUNT` SET "
	var params []interface{}
	
	for key, value := range body {
		if key != "id" {
			if key == "password" {
				if tools.PasswordNotStrong(value.(string)) {
					tools.JsonResponse(w, 400, `{"message": "Password not strong enough"}`)
					return
				}
				value = tools.HashPassword(value.(string))
			}
			tools.AppendUpdate(&request, &params, key, value)
		}
	}

	// Removing the last ","
	request = request[:len(request)-2]

	request += " WHERE `id` = ?"
	params = append(params, id)

	// Updating the account in the database
	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Account updated"`
	
	// Adding the return fields of the query
	fields, err := AccountGetAll(db, id, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func AccountDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
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
	if !tools.ElementExists(db, "ACCOUNT", "id", id) {
		tools.JsonResponse(w, 400, `{"message": "Account does not exist"}`)
		return
	}

	// Deleting the account in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `ACCOUNT` WHERE `id` = ?", id)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Account deleted", "id": "` + id + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func AccountGetAll(db *sql.DB, uuid string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `id`, `username`, `first_name`, `last_name`, `email`, `account_type`, `creation_date` FROM `ACCOUNT` WHERE `id` = ?", uuid)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return AccountGetAllAssociation(result, arrayOutput)
}

func AccountGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var id, username, firstName, lastName, email, accountType, creation_date string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&id, &username, &firstName, &lastName, &email, &accountType, &creation_date)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"id": "` + id + `", "username": "` + username + `", "first_name": "` + firstName + `", "last_name": "` + lastName + `", "email": "` + email + `", "account_type": "` + accountType + `", "creation_date": "` + creation_date + `"},`
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&id, &username, &firstName, &lastName, &email, &accountType, &creation_date)
			if err != nil {
				return "", err
			}
		}
		return `"id": "` + id + `", "username": "` + username + `", "first_name": "` + firstName + `", "last_name": "` + lastName + `", "email": "` + email + `", "account_type": "` + accountType + `", "creation_date": "` + creation_date + `"`, nil
	}
}