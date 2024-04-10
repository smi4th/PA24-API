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
		ServicesPut(w, r, db)
	case "DELETE":
		ServicesDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func ServicesPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `price`, `service_type`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    price_ := tools.BodyValueToString(body, "price")
	service_type_ := tools.BodyValueToString(body, "service_type")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(price_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, price_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, price_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

    if !tools.ValueIsEmpty(service_type_) {
		if !tools.ElementExists(db, "SERVICES_TYPES", "uuid", service_type_) {
			tools.JsonResponse(w, 400, `{"error": "This service_type does not exist"}`) 
			return
		}
	}
	
	
	

	

	

	uuid_ := tools.GenerateUUID()

	// Inserting the Services in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `SERVICES` (`uuid`, `price`, `service_type`) VALUES (?, ?, ?)", uuid_, price_, service_type_)
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
	if tools.AtLeastOneValueInQuery(query, `uuid`, `price`, `service_type`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `price`, `service_type` FROM `SERVICES`"
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
	jsonResponse, err := ServicesGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ServicesPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `price`, `service_type`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	
    price_ := tools.BodyValueToString(body, "price")
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
	if tools.ValueTooShort(4, price_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, price_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

    if !tools.ValueIsEmpty(service_type_) {
		if !tools.ElementExists(db, "SERVICES_TYPES", "uuid", service_type_) {
			tools.JsonResponse(w, 400, `{"error": "This service_type does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExists(db, "SERVICES", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Services does not exist"}`) 
		return
	}
	

	{{passwordCheck}}

    {{emailCheck}}

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
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `price`, `service_type` FROM `SERVICES` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return ServicesGetAllAssociation(result, arrayOutput)
}

func ServicesGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, price_, service_type_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &price_, &service_type_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "price": "` + price_ + `", "service_type": "` + service_type_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &price_, &service_type_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "price": "` + price_ + `", "service_type": "` + service_type_ + `"`, nil
	}
}