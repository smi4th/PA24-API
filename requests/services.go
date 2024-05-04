package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Services(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		ServicesPost(w, r, db)
	case "GET":
		ServicesGet(w, r, db)
	case "PUT":
		if tools.GetUUID(r, db) == tools.GetElement(db, "SERVICES", "account", "uuid", tools.ReadQuery(r)["uuid"]) || tools.IsAdmin(r, db) {
			ServicesPut(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"error": "Forbidden"}`)
		}
	case "DELETE":
		if tools.GetUUID(r, db) == tools.GetElement(db, "SERVICES", "account", "uuid", tools.ReadQuery(r)["uuid"]) || tools.IsAdmin(r, db) {
			ServicesDelete(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"error": "Forbidden"}`)
		}
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func ServicesPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `price`, `description`, `account`, `service_type`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    price_ := tools.BodyValueToString(body, "price")
	description_ := tools.BodyValueToString(body, "description")
	account_ := tools.BodyValueToString(body, "account")
	service_type_ := tools.BodyValueToString(body, "service_type")

	if tools.GetUUID(r, db) != tools.GetElement(db, "ACCOUNT", "uuid", "uuid", account_) && !tools.IsAdmin(r, db) {
		tools.JsonResponse(w, 403, `{"error": "Forbidden"}`)
		return
	}
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(price_, description_, account_, service_type_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, description_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}

	if !tools.ElementExists(db, "SERVICES_TYPES", "uuid", service_type_) {
		tools.JsonResponse(w, 400, `{"error": "This service_type does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
		tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`)
		return
	}
	

	

	uuid_ := tools.GenerateUUID()

	// Inserting the Services in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `SERVICES` (`uuid`, `price`, `description`, `account`, `service_type`) VALUES (?, ?, ?, ?, ?)", uuid_, price_, description_, account_, service_type_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Services created"`

	// Adding the return fields of the query
	fields, err := ServicesGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func ServicesGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`, `price`, `description`, `account`, `service_type`, `all`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `price`, `description`, `account`, `service_type` FROM `SERVICES`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `SERVICES`"
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
	jsonResponse, err := ServicesGetAllAssociation(result, true)
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

func ServicesPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `price`, `description`, `account`, `service_type`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	
	description_ := tools.BodyValueToString(body, "description")
	account_ := tools.BodyValueToString(body, "account")
	service_type_ := tools.BodyValueToString(body, "service_type")
	

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
	if tools.ValueTooShort(4, description_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}

    if !tools.ValueIsEmpty(service_type_) {
		if !tools.ElementExists(db, "SERVICES_TYPES", "uuid", service_type_) {
			tools.JsonResponse(w, 400, `{"error": "This service_type does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`)
			return
		}
	}

	

	if !tools.ElementExists(db, "SERVICES", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Services does not exist"}`) 
		return
	}
	

	

    

	request := "UPDATE `SERVICES` SET "
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
	jsonResponse := `{"message": "Services updated"`
	
	// Adding the return fields of the query
	fields, err := ServicesGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func ServicesDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "SERVICES", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Services does not exist"}`) 
		return
	}
	

	// Deleting the Services in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `SERVICES` WHERE uuid = ?", uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Services deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ServicesGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `price`, `description`, `account`, `service_type` FROM `SERVICES` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return ServicesGetAllAssociation(result, arrayOutput)
}

func ServicesGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, price_, description_, account_, service_type_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &price_, &description_, &account_, &service_type_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "price": "` + price_ + `", "description": "` + description_ + `", "account": "` + account_ + `", "service_type": "` + service_type_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &price_, &description_, &account_, &service_type_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "price": "` + price_ + `", "description": "` + description_ + `", "account": "` + account_ + `", "service_type": "` + service_type_ + `"`, nil
	}
}