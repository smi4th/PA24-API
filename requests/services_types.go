package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func ServicesTypes(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		ServicesTypesPost(w, r, db)
	case "GET":
		ServicesTypesGet(w, r, db)
	case "PUT":
		ServicesTypesPut(w, r, db)
	case "DELETE":
		ServicesTypesDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func ServicesTypesPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `type`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    type_ := tools.BodyValueToString(body, "type")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(type_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, type_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, type_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

    
	
	

	

	if tools.ElementExists(db, "SERVICES_TYPES", "type", type_) {
		tools.JsonResponse(w, 400, `{"error": "This type already exists"}`) 
		return
	}
	

	uuid_ := tools.GenerateUUID()

	// Inserting the ServicesTypes in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `SERVICES_TYPES` (`uuid`, `type`) VALUES (?, ?)", uuid_, type_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "ServicesTypes created"`

	// Adding the return fields of the query
	fields, err := ServicesTypesGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func ServicesTypesGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`, `type`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `type` FROM `SERVICES_TYPES`"
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
	jsonResponse, err := ServicesTypesGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ServicesTypesPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `type`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	
    type_ := tools.BodyValueToString(body, "type")
	

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
	if tools.ValueTooShort(4, type_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, type_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

    

	if !tools.ElementExists(db, "SERVICES_TYPES", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This ServicesTypes does not exist"}`) 
		return
	}
	if tools.ElementExists(db, "SERVICES_TYPES", "type", type_) {
		tools.JsonResponse(w, 400, `{"error": "This type already exists"}`) 
		return
	}
	

	{{passwordCheck}}

    {{emailCheck}}

	request := "UPDATE `SERVICES_TYPES` SET "
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
	jsonResponse := `{"message": "ServicesTypes updated"`
	
	// Adding the return fields of the query
	fields, err := ServicesTypesGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func ServicesTypesDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "SERVICES_TYPES", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This ServicesTypes does not exist"}`) 
		return
	}
	

	// Deleting the ServicesTypes in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `SERVICES_TYPES` WHERE uuid = ?", uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "ServicesTypes deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ServicesTypesGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `type` FROM `SERVICES_TYPES` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return ServicesTypesGetAllAssociation(result, arrayOutput)
}

func ServicesTypesGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, type_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &type_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "type": "` + type_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &type_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "type": "` + type_ + `"`, nil
	}
}