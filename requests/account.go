package requests

import (
	"net/http"
	"tools"
	"database/sql"
	"strings"
)

func Account(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		AccountPost(w, r, db)
	case "GET":
		AccountGet(w, r, db)
	case "PUT":
		if tools.GetUUID(r, db) == tools.GetElement(db, "ACCOUNT", "uuid", "uuid", tools.ReadQuery(r)["uuid"]) || tools.IsAdmin(r, db) {
			AccountPut(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	case "DELETE":
		if tools.GetUUID(r, db) == tools.GetElement(db, "ACCOUNT", "uuid", "uuid", tools.ReadQuery(r)["uuid"]) || tools.IsAdmin(r, db) {
			AccountDelete(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
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
	imgPath_ := tools.BodyValueToString(body, "imgPath")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(username_, password_, first_name_, last_name_, email_, account_type_) {
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

	if !tools.ElementExists(db, "ACCOUNT_TYPE", "uuid", account_type_) {
		tools.JsonResponse(w, 400, `{"error": "This account_type does not exist"}`) 
		return
	}
	
	
	if tools.PasswordNotStrong(password_) {
		tools.JsonResponse(w, 400, `{"error": "Password is not strong enough"}`) 
		return
	} else {
		password_ = tools.HashPassword(password_)
	}

	if !tools.EmailIsValid(email_) {
		tools.JsonResponse(w, 400, `{"error": "Email is not valid"}`) 
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
	

	uuid_ := tools.GenerateUUID()

	// Inserting the Account in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `ACCOUNT` (`uuid`, `username`, `password`, `first_name`, `last_name`, `email`, `account_type`, `imgPath`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", uuid_, username_, password_, first_name_, last_name_, email_, account_type_, imgPath_)
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
	if tools.AtLeastOneValueInQuery(query, `uuid`, `username`, `first_name`, `last_name`, `email`, `creation_date`, `account_type`, "all", "provider", "imgPath", "token") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `username`, `first_name`, `last_name`, `email`, `creation_date`, `account_type`, `provider`, `imgPath` FROM `ACCOUNT`"
	var params []interface{}

	countRequest := "SELECT COUNT(*) FROM `ACCOUNT`"
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
	jsonResponse, err := AccountGetAllAssociation(result, true)
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

func AccountPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `username`, `password`, `first_name`, `last_name`, `email`, `account_type`, `provider`, "imgPath") || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	if !tools.AtLeastOneValueInBody(body, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Cannot update all fields"}`)
		return
	}

	uuid_ := query["uuid"]
	
    username_ := tools.BodyValueToString(body, "username")
	password_ := tools.BodyValueToString(body, "password")
	first_name_ := tools.BodyValueToString(body, "first_name")
	last_name_ := tools.BodyValueToString(body, "last_name")
	email_ := tools.BodyValueToString(body, "email")
	account_type_ := tools.BodyValueToString(body, "account_type")
	provider_ := tools.BodyValueToString(body, "provider")
	

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

	if !tools.ValueIsEmpty(provider_) {
		if !tools.ElementExists(db, "PROVIDER", "uuid", provider_) {
			tools.JsonResponse(w, 400, `{"error": "This provider does not exist"}`)
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
	
	testsRequests := []string{
		"SELECT count(*) FROM `BASKET_BEDROOM` WHERE `BEDROOM` IN (SELECT `uuid` FROM `BED_ROOM` WHERE `HOUSING` IN (SELECT `uuid` FROM `HOUSING` WHERE `account` = ?))",
		"SELECT count(*) FROM `BASKET_EQUIPMENT` WHERE `EQUIPMENT` IN (SELECT `uuid` FROM `EQUIPMENT` WHERE `HOUSING` IN (SELECT `uuid` FROM `HOUSING` WHERE `account` = ?))",
		"SELECT count(*) FROM `BASKET_SERVICE` WHERE `SERVICE` IN (SELECT `uuid` FROM `SERVICES` WHERE `account` = ?)",
		"SELECT count(*) FROM `BASKET_HOUSING` WHERE `HOUSING` IN (SELECT `uuid` FROM `HOUSING` WHERE `account` = ?)",
	}

	for _, request := range testsRequests {
		result, err := tools.ExecuteQuery(db, request, uuid_)
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

		if count != "0" {
			tools.JsonResponse(w, 400, `{"message": "This account has ` + count + " " + strings.Split(strings.Split(request, "FROM `")[1], "`")[0] + ` in some baskets"}`)
			return
		}
	}

	deleteRequests := []string{
		"DELETE FROM `MESSAGE` WHERE `author` = ?",
		"DELETE FROM `MESSAGE` WHERE `account` = ?",
		"DELETE FROM `REVIEW` WHERE `account` = ?",
		"DELETE FROM `REVIEW` WHERE `HOUSING` IN (SELECT `uuid` FROM `HOUSING` WHERE `account` = ?)",
		"DELETE FROM `REVIEW` WHERE `SERVICE` IN (SELECT `uuid` FROM `SERVICES` WHERE `account` = ?)",
		"DELETE FROM `REVIEW` WHERE `BEDROOM` IN (SELECT `uuid` FROM `BED_ROOM` WHERE `HOUSING` IN (SELECT `uuid` FROM `HOUSING` WHERE `account` = ?))",
		"DELETE FROM `BASKET_SERVICE` WHERE `BASKET` IN (SELECT `uuid` FROM `BASKET` WHERE `account` = ?)",
		"DELETE FROM `BASKET_HOUSING` WHERE `BASKET` IN (SELECT `uuid` FROM `BASKET` WHERE `account` = ?)",
		"DELETE FROM `BASKET_BEDROOM` WHERE `BASKET` IN (SELECT `uuid` FROM `BASKET` WHERE `account` = ?)",
		"DELETE FROM `BASKET_EQUIPMENT` WHERE `BASKET` IN (SELECT `uuid` FROM `BASKET` WHERE `account` = ?)",
		"DELETE FROM `BASKET_SERVICE` WHERE `BASKET` IN (SELECT `uuid` FROM `BASKET` WHERE `account` = ?)",
		"DELETE FROM `BED_ROOM` WHERE `HOUSING` IN (SELECT `uuid` FROM `HOUSING` WHERE `account` = ?)",
		"DELETE FROM `EQUIPMENT` WHERE `HOUSING` IN (SELECT `uuid` FROM `HOUSING` WHERE `account` = ?)",
		"DELETE FROM `HOUSING` WHERE `account` = ?",
		"DELETE FROM `DISPONIBILITY` WHERE `account` = ?",
		"DELETE FROM `SERVICES` WHERE `account` = ?",
		"DELETE FROM `ACCOUNT_SUBSCRIPTION` WHERE `account` = ?",
		"DELETE FROM `ACCOUNT` WHERE `uuid` = ?",
	}

	// Create the transaction
	tx, err := db.Begin()
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	for _, request := range deleteRequests {
		_, err := tx.Exec(request, uuid_)
		if err != nil {
			tx.Rollback()
			tools.ErrorLog(err.Error())
			tools.ErrorLog(request)
			tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
			return
		}
	}

	tx.Commit()
	
	// Creating the response	
	jsonResponse := `{"message": "Account deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func AccountGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `username`, `first_name`, `last_name`, `email`, `creation_date`, `account_type`, `provider`, `imgPath` FROM `ACCOUNT` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return AccountGetAllAssociation(result, arrayOutput)
}

func AccountGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, username_, first_name_, last_name_, email_, creation_date_, account_type_, imgPath_ string
	var provider_ sql.NullString

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &username_, &first_name_, &last_name_, &email_, &creation_date_, &account_type_, &provider_, &imgPath_)
			if err != nil {
				return "", err
			}
			var provider string
			if provider_.Valid {
				provider = provider_.String
			} else {
				provider = ""
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "username": "` + username_ + `", "first_name": "` + first_name_ + `", "last_name": "` + last_name_ + `", "email": "` + email_ + `", "creation_date": "` + creation_date_ + `", "account_type": "` + account_type_ + `", "provider": "` + provider + `", "imgPath": "` + imgPath_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &username_, &first_name_, &last_name_, &email_, &creation_date_, &account_type_, &provider_, &imgPath_)
			if err != nil {
				return "", err
			}
		}
		var provider string
		if provider_.Valid {
			provider = provider_.String
		} else {
			provider = ""
		}
		return `"uuid": "` + uuid_ + `", "username": "` + username_ + `", "first_name": "` + first_name_ + `", "last_name": "` + last_name_ + `", "email": "` + email_ + `", "creation_date": "` + creation_date_ + `", "account_type": "` + account_type_ + `", "provider": "` + provider + `", "imgPath": "` + imgPath_ + `"`, nil
	}
}