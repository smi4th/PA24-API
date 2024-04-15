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
	if tools.ValuesNotInBody(body, `username`, `password`, `first_name`, `last_name`, `email`, `account_type`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    username_ := tools.BodyValueToString(body, "username")
	password_ := tools.BodyValueToString(body, "password")
	first_name_ := tools.BodyValueToString(body, "first_name")
	last_name_ := tools.BodyValueToString(body, "last_name")
	email_ := tools.BodyValueToString(body, "email")
	account_type_ := tools.BodyValueToString(body, "account_type")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(username_, password_, first_name_, last_name_, email_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, username_, password_, first_name_, last_name_, email_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, username_, password_, first_name_, last_name_, email_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

    if !tools.ValueIsEmpty(account_type_) {
		if !tools.ElementExists(db, "ACCOUNT_TYPE", "uuid", account_type_) {
			tools.JsonResponse(w, 400, `{"error": "This account_type does not exist"}`) 
			return
		}
	}
	
	
	if !tools.ValueIsEmpty(password_) {
		if tools.PasswordNotStrong(password_) {
			tools.JsonResponse(w, 400, `{"error": "Password is not strong enough"}`) 
			return
		} else {
			password_ = tools.HashPassword(password_)
		}
	}

	if !tools.ValueIsEmpty(email_) {
		if !tools.EmailIsValid(email_) {
			tools.JsonResponse(w, 400, `{"error": "Email is not valid"}`) 
			return
		}
	}

	if tools.ElementExists(db, "ACCOUNT", "username", username_) {
		tools.JsonResponse(w, 400, `{"error": "This username already exists"}`) 
		return
	}
	if tools.ElementExists(db, "ACCOUNT", "email", email_) {
		tools.JsonResponse(w, 400, `{"error": "This email already exists"}`) 
		return
	}
	

	uuid_ := tools.GenerateUUID()

	// Inserting the Account in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `ACCOUNT` (`uuid`, `username`, `password`, `first_name`, `last_name`, `email`, `account_type`) VALUES (?, ?, ?, ?, ?, ?, ?)", uuid_, username_, password_, first_name_, last_name_, email_, account_type_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Account created"`

	// Adding the return fields of the query
	fields, err := AccountGetAll(db, uuid_, false)
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
	if tools.AtLeastOneValueInQuery(query, `uuid`, `username`, `first_name`, `last_name`, `email`, `creation_date`, `account_type`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `username`, `first_name`, `last_name`, `email`, `creation_date`, `account_type` FROM `ACCOUNT`"
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
	if tools.AtLeastOneValueInBody(body, `username`, `password`, `first_name`, `last_name`, `email`, `account_type`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	
    username_ := tools.BodyValueToString(body, "username")
	password_ := tools.BodyValueToString(body, "password")
	first_name_ := tools.BodyValueToString(body, "first_name")
	last_name_ := tools.BodyValueToString(body, "last_name")
	email_ := tools.BodyValueToString(body, "email")
	account_type_ := tools.BodyValueToString(body, "account_type")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(uuid_) {
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
	if tools.ValueTooShort(4, username_, password_, first_name_, last_name_, email_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, username_, password_, first_name_, last_name_, email_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

    if !tools.ValueIsEmpty(account_type_) {
		if !tools.ElementExists(db, "ACCOUNT_TYPE", "uuid", account_type_) {
			tools.JsonResponse(w, 400, `{"error": "This account_type does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExists(db, "ACCOUNT", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Account does not exist"}`) 
		return
	}
	if tools.ElementExists(db, "ACCOUNT", "username", username_) {
		tools.JsonResponse(w, 400, `{"error": "This username already exists"}`) 
		return
	}
	if tools.ElementExists(db, "ACCOUNT", "email", email_) {
		tools.JsonResponse(w, 400, `{"error": "This email already exists"}`) 
		return
	}
	

	if !tools.ValueIsEmpty(password_) {
		if tools.PasswordNotStrong(password_) {
			tools.JsonResponse(w, 400, `{"error": "Password is not strong enough"}`) 
			return
		} else {
			password_ = tools.HashPassword(password_)
		}
	}

    if !tools.ValueIsEmpty(email_) {
		if !tools.EmailIsValid(email_) {
			tools.JsonResponse(w, 400, `{"error": "Email is not valid"}`) 
			return
		}
	}

	request := "UPDATE `ACCOUNT` SET "
	var params []interface{}
	
	for key, value := range body {
		if !tools.ValueInArray(key, `uuid`) {
			if key == "password" {
				value = tools.HashPassword(value.(string))
			}
			tools.AppendUpdate(&request, &params, key, value)
		}
	}

	// Removing the last ","
	request = request[:len(request)-2]

	request += " WHERE uuid = ?"
	params = append(params, uuid_)

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
	fields, err := AccountGetAll(db, uuid_, false)
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
	if tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "ACCOUNT", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Account does not exist"}`) 
		return
	}
	

	// Deleting the Account in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `ACCOUNT` WHERE uuid = ?", uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Account deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func AccountGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `username`, `first_name`, `last_name`, `email`, `creation_date`, `account_type` FROM `ACCOUNT` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return AccountGetAllAssociation(result, arrayOutput)
}

func AccountGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, username_, first_name_, last_name_, email_, creation_date_, account_type_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &username_, &first_name_, &last_name_, &email_, &creation_date_, &account_type_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "username": "` + username_ + `", "first_name": "` + first_name_ + `", "last_name": "` + last_name_ + `", "email": "` + email_ + `", "creation_date": "` + creation_date_ + `", "account_type": "` + account_type_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &username_, &first_name_, &last_name_, &email_, &creation_date_, &account_type_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "username": "` + username_ + `", "first_name": "` + first_name_ + `", "last_name": "` + last_name_ + `", "email": "` + email_ + `", "creation_date": "` + creation_date_ + `", "account_type": "` + account_type_ + `"`, nil
	}
}