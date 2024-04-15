package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Message(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		MessagePost(w, r, db)
	case "GET":
		MessageGet(w, r, db)
	case "PUT":
		MessagePut(w, r, db)
	case "DELETE":
		MessageDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func MessagePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `uuid`, `creation_date`, `content`, `note`, `account`, `author`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    uuid_ := tools.BodyValueToString(body, "uuid")
	creation_date_ := tools.BodyValueToString(body, "creation_date")
	content_ := tools.BodyValueToString(body, "content")
	note_ := tools.BodyValueToString(body, "note")
	account_ := tools.BodyValueToString(body, "account")
	author_ := tools.BodyValueToString(body, "author")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(uuid_, creation_date_, content_, note_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, uuid_, creation_date_, content_, note_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, uuid_, creation_date_, content_, note_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

    if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(author_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", author_) {
			tools.JsonResponse(w, 400, `{"error": "This author does not exist"}`) 
			return
		}
	}
	
	
	

	

	

	

	// Inserting the Message in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `MESSAGE` (`account`, `author`, `uuid`, `uuid`, `creation_date`, `content`, `note`, `account`, `author`) VALUES (?, ?, ?, ?, ?, ?)", account_, author_, uuid_, uuid_, creation_date_, content_, note_, account_, author_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Message created"`

	// Adding the return fields of the query
	fields, err := MessageGetAll(db, account_, author_, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func MessageGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`, `creation_date`, `content`, `note`, `account`, `author`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `creation_date`, `content`, `note`, `account`, `author` FROM `MESSAGE`"
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
	jsonResponse, err := MessageGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func MessagePut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `creation_date`, `content`, `note`) || tools.ValuesNotInQuery(query, `account`, `author`, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	account_ := query["account"]
	author_ := query["author"]
	
    creation_date_ := tools.BodyValueToString(body, "creation_date")
	content_ := tools.BodyValueToString(body, "content")
	note_ := tools.BodyValueToString(body, "note")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(account_, author_, uuid_) {
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
	if tools.ValueTooShort(4, creation_date_, content_, note_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, creation_date_, content_, note_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

    if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(author_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", author_) {
			tools.JsonResponse(w, 400, `{"error": "This author does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExists(db, "MESSAGE", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This Message does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "MESSAGE", "author", author_) {
		tools.JsonResponse(w, 400, `{"error": "This Message does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "MESSAGE", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Message does not exist"}`) 
		return
	}
	

	

    

	request := "UPDATE `MESSAGE` SET "
	var params []interface{}
	
	for key, value := range body {
		if !tools.ValueInArray(key, `account`, `author`, `uuid`) {
			if key == "password" {
				value = tools.HashPassword(value.(string))
			}
			tools.AppendUpdate(&request, &params, key, value)
		}
	}

	// Removing the last ","
	request = request[:len(request)-2]

	request += " WHERE account = ?, author = ?, uuid = ?"
	params = append(params, account_, author_, uuid_)

	// Updating the account in the database
	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Message updated"`
	
	// Adding the return fields of the query
	fields, err := MessageGetAll(db, account_, author_, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func MessageDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `account`, `author`, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]
	author_ := query["author"]
	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "MESSAGE", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This Message does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "MESSAGE", "author", author_) {
		tools.JsonResponse(w, 400, `{"error": "This Message does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "MESSAGE", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Message does not exist"}`) 
		return
	}
	

	// Deleting the Message in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `MESSAGE` WHERE account = ?, author = ?, uuid = ?", account_, author_, uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Message deleted", "account": "` + account_ + `", "author": "` + author_ + `", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func MessageGetAll(db *sql.DB, account_ string, author_ string, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `creation_date`, `content`, `note`, `account`, `author` FROM `MESSAGE` WHERE account = ?, author = ?, uuid = ?", account_, author_, uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return MessageGetAllAssociation(result, arrayOutput)
}

func MessageGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, creation_date_, content_, note_, account_, author_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &creation_date_, &content_, &note_, &account_, &author_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "creation_date": "` + creation_date_ + `", "content": "` + content_ + `", "note": "` + note_ + `", "account": "` + account_ + `", "author": "` + author_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &creation_date_, &content_, &note_, &account_, &author_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "creation_date": "` + creation_date_ + `", "content": "` + content_ + `", "note": "` + note_ + `", "account": "` + account_ + `", "author": "` + author_ + `"`, nil
	}
}