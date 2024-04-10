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
		DisponibilityPut(w, r, db)
	case "DELETE":
		DisponibilityDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func DisponibilityPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `start_date`, `end_date`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    start_date_ := tools.BodyValueToString(body, "start_date")
	end_date_ := tools.BodyValueToString(body, "end_date")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(start_date_, end_date_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, start_date_, end_date_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, start_date_, end_date_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

    
	
	

	

	

	uuid_ := tools.GenerateUUID()

	// Inserting the Disponibility in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `DISPONIBILITY` (`uuid`, `start_date`, `end_date`) VALUES (?, ?, ?)", uuid_, start_date_, end_date_)
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
	if tools.AtLeastOneValueInQuery(query, `uuid`, `start_date`, `end_date`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `start_date`, `end_date` FROM `DISPONIBILITY`"
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
	jsonResponse, err := DisponibilityGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

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

	uuid_ := query["uuid"]
	
    start_date_ := tools.BodyValueToString(body, "start_date")
	end_date_ := tools.BodyValueToString(body, "end_date")
	

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
	if tools.ValueTooShort(4, start_date_, end_date_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, start_date_, end_date_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

    

	if !tools.ElementExists(db, "DISPONIBILITY", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Disponibility does not exist"}`) 
		return
	}
	

	{{passwordCheck}}

    {{emailCheck}}

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
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `start_date`, `end_date` FROM `DISPONIBILITY` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return DisponibilityGetAllAssociation(result, arrayOutput)
}

func DisponibilityGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, start_date_, end_date_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &start_date_, &end_date_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "start_date": "` + start_date_ + `", "end_date": "` + end_date_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &start_date_, &end_date_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "start_date": "` + start_date_ + `", "end_date": "` + end_date_ + `"`, nil
	}
}