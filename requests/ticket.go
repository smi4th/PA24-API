package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Ticket(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		TicketPost(w, r, db)
	case "GET":
		TicketGet(w, r, db)
	case "PUT":
		if tools.IsAdmin(r, db) {
			TicketPut(w, r, db)
		} else {
			tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		}
	case "DELETE":
		if tools.IsAdmin(r, db) {
			TicketDelete(w, r, db)
		} else {
			tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		}
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func TicketPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `title`, `description`, `status`, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    title_ := tools.BodyValueToString(body, "title")
	description_ := tools.BodyValueToString(body, "description")
	status_ := tools.BodyValueToString(body, "status")
	account_ := tools.BodyValueToString(body, "account")

	if tools.GetUUID(r, db) != account_ && !tools.IsAdmin(r, db) {
		tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		return
	}

	// Checking if the values are empty
	if tools.ValueIsEmpty(title_, description_, status_, account_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, title_, description_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, title_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

	uuid_ := tools.GenerateUUID()

	// Inserting the Ticket in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `TICKET` (`uuid`, `title`, `description`, `status`, `account`, `support`) VALUES (?, ?, ?, ?, ?, NULL)", uuid_, title_, description_, status_, account_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Ticket created"`

	// Adding the return fields of the query
	fields, err := TicketGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func TicketGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`, `all`, `title`, `description`, `status`, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `title`, `description`, `creation_date`, `status`, `account`, CASE WHEN `support` IS NULL THEN 'NULL' ELSE `support` END AS support FROM `TICKET`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `TICKET`"
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
	jsonResponse, err := TicketGetAllAssociation(result, true)
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

func TicketPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `title`, `description`, `creation_date`, `status`, `account`, `support`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	if !tools.AtLeastOneValueInBody(body, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Empty fields"}`)
		return
	}

	uuid_ := query["uuid"]
	
    title_ := tools.BodyValueToString(body, "title")
	description_ := tools.BodyValueToString(body, "description")
	status_ := tools.BodyValueToString(body, "status")
	account_ := tools.BodyValueToString(body, "account")
	support_ := tools.BodyValueToString(body, "support")

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
	if tools.ValueTooShort(4, title_, description_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, title_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}
    

	if !tools.ElementExists(db, "Ticket", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Ticket does not exist"}`) 
		return
	}
	if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This Account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(support_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", support_) {
			tools.JsonResponse(w, 400, `{"error": "This Account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(status_) {
		if !tools.ElementExists(db, "STATUS", "uuid", status_) {
			tools.JsonResponse(w, 400, `{"error": "This Status does not exist"}`) 
			return
		}
	}

	request := "UPDATE `TICKET` SET "
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
	jsonResponse := `{"message": "Ticket updated"`
	
	// Adding the return fields of the query
	fields, err := TicketGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func TicketDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "TICKET", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Ticket does not exist"}`) 
		return
	}
	

	// Deleting the Ticket in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `TICKET` WHERE uuid = ?", uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Ticket deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func TicketGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `title`, `description`, `creation_date`, `status`, `account`, CASE WHEN `support` IS NULL THEN 'NULL' ELSE `support` END AS support FROM `TICKET` WHERE UUID = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return TicketGetAllAssociation(result, arrayOutput)
}

func TicketGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, title_, description_, creation_date_, status_, account_, support_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &title_, &description_, &creation_date_, &status_, &account_, &support_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "title": "` + title_ + `", "description": "` + description_ + `", "creation_date": "` + creation_date_ + `", "status": "` + status_ + `", "account": "` + account_ + `", "support": "` + support_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &title_, &description_, &creation_date_, &status_, &account_, &support_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "title": "` + title_ + `", "description": "` + description_ + `", "creation_date": "` + creation_date_ + `", "status": "` + status_ + `", "account": "` + account_ + `", "support": "` + support_ + `"`, nil
	}
}