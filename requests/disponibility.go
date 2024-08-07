package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Disponibility(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		DisponibilityPost(w, r, db)
	case "GET":
		DisponibilityGet(w, r, db)
	case "PUT":
		if tools.GetUUID(r, db) == tools.GetElement(db, "DISPONIBILITY", "account", "uuid", tools.ReadQuery(r)["uuid"]) || tools.IsAdmin(r, db) {
			DisponibilityPut(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	case "DELETE":
		if tools.GetUUID(r, db) == tools.GetElement(db, "DISPONIBILITY", "account", "uuid", tools.ReadQuery(r)["uuid"]) || tools.IsAdmin(r, db) {
			DisponibilityDelete(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func DisponibilityPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `start_date`, `end_date`, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    start_date_ := tools.BodyValueToString(body, "start_date")
	end_date_ := tools.BodyValueToString(body, "end_date")
	account_ := tools.BodyValueToString(body, "account")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(start_date_, end_date_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
		tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`)
		return
	}
	

	

	uuid_ := tools.GenerateUUID()

	// Inserting the Disponibility in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `DISPONIBILITY` (`uuid`, `start_date`, `end_date`, `account`) VALUES (?, ?, ?, ?)", uuid_, start_date_, end_date_, account_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Disponibility created"`

	// Adding the return fields of the query
	fields, err := DisponibilityGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func DisponibilityGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`, `start_date`, `end_date`, `account`, `all`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `start_date`, `end_date`, `account` FROM `DISPONIBILITY`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `DISPONIBILITY`"
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
	jsonResponse, err := DisponibilityGetAllAssociation(result, true)
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

func DisponibilityPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `start_date`, `end_date`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	if !tools.AtLeastOneValueInBody(body, `uuid`, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Cannot update all fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

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

    

	if !tools.ElementExists(db, "DISPONIBILITY", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Disponibility does not exist"}`) 
		return
	}
	

    

	request := "UPDATE `DISPONIBILITY` SET "
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
	jsonResponse := `{"message": "Disponibility updated"`
	
	// Adding the return fields of the query
	fields, err := DisponibilityGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func DisponibilityDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "DISPONIBILITY", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Disponibility does not exist"}`) 
		return
	}
	

	// Deleting the Disponibility in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `DISPONIBILITY` WHERE uuid = ?", uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Disponibility deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func DisponibilityGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `start_date`, `end_date`, `account` FROM `DISPONIBILITY` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return DisponibilityGetAllAssociation(result, arrayOutput)
}

func DisponibilityGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, start_date_, end_date_, account_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &start_date_, &end_date_, &account_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "start_date": "` + start_date_ + `", "end_date": "` + end_date_ + `", "account": "` + account_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &start_date_, &end_date_, &account_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "start_date": "` + start_date_ + `", "end_date": "` + end_date_ + `", "account": "` + account_ + `"`, nil
	}
}