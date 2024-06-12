package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func TMessage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		TMessagePost(w, r, db)
	case "GET":
		TMessageGet(w, r, db)
	case "PUT":
		if tools.IsAdmin(r, db) {
			TMessagePut(w, r, db)
		} else {
			tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		}
	case "DELETE":
		if tools.IsAdmin(r, db) {
			TMessageDelete(w, r, db)
		} else {
			tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		}
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func TMessagePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `content`, `ticket`, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    content_ := tools.BodyValueToString(body, "content")
	ticket_ := tools.BodyValueToString(body, "ticket")
	account_ := tools.BodyValueToString(body, "account")

	if tools.GetUUID(r, db) != account_ && !tools.IsAdmin(r, db) {
		tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		return
	}

	// Checking if the values are empty
	if tools.ValueIsEmpty(content_, ticket_, account_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, content_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}

	uuid_ := tools.GenerateUUID()

	// Inserting the TMessage in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `TMESSAGE` (`uuid`, `content`, `ticket`, `account`) VALUES (?, ?, ?, ?)", uuid_, content_, ticket_, account_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "TMessage created"`

	// Adding the return fields of the query
	fields, err := TMessageGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func TMessageGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`, `all`, `content`, `creation_date`, `ticket`, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `content`, `creation_date`, `ticket`, `account` FROM `TMESSAGE`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `TMESSAGE`"
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
	jsonResponse, err := TMessageGetAllAssociation(result, true)
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

func TMessagePut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `content`, `creation_date`, `ticket`, `account`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	if !tools.AtLeastOneValueInBody(body, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Empty fields"}`)
		return
	}

	uuid_ := query["uuid"]
	
	ticket_ := tools.BodyValueToString(body, "ticket")
	account_ := tools.BodyValueToString(body, "account")
	content_ := tools.BodyValueToString(body, "content")

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
	if tools.ValueTooShort(4, content_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
    

	if !tools.ElementExists(db, "TMessage", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This TMessage does not exist"}`) 
		return
	}
	if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This Account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(ticket_) {
		if !tools.ElementExists(db, "TICKET", "uuid", ticket_) {
			tools.JsonResponse(w, 400, `{"error": "This Ticket does not exist"}`)
			return
		}
	}

	request := "UPDATE `TMESSAGE` SET "
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
	jsonResponse := `{"message": "TMessage updated"`
	
	// Adding the return fields of the query
	fields, err := TMessageGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func TMessageDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "TMESSAGE", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This TMessage does not exist"}`) 
		return
	}
	

	// Deleting the TMessage in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `TMESSAGE` WHERE uuid = ?", uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "TMessage deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func TMessageGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `content`, `creation_date`, `ticket`, `account` FROM `TMESSAGE` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return TMessageGetAllAssociation(result, arrayOutput)
}

func TMessageGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, content_, creation_date_, ticket_, account_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &content_, &creation_date_, &ticket_, &account_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "content": "` + content_ + `", "creation_date": "` + creation_date_ + `", "ticket": "` + ticket_ + `", "account": "` + account_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &content_, &creation_date_, &ticket_, &account_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "content": "` + content_ + `", "creation_date": "` + creation_date_ + `", "ticket": "` + ticket_ + `", "account": "` + account_ + `"`, nil
	}
}