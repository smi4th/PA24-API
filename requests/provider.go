package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Provider(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		ProviderPost(w, r, db)
	case "GET":
		ProviderGet(w, r, db)
	case "PUT":
		ProviderPut(w, r, db)
	case "DELETE":
		ProviderDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func ProviderPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, "name", "email") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	name := tools.BodyValueToString(body, "name")
	email := tools.BodyValueToString(body, "email")

	// Checking if the values are empty
	if tools.ValueIsEmpty(name, email) {
		tools.JsonResponse(w, 400, `{"message": "Empty fields"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(8, name) {
		tools.JsonResponse(w, 400, `{"message": "Value too short"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooLong(64, name) {
		tools.JsonResponse(w, 400, `{"message": "Value too long"}`)
		return
	}

	if tools.EmailIsValid(email) == false {
		tools.JsonResponse(w, 400, `{"message": "Invalid email"}`)
		return
	}
	
	// Checking if the name is already taken
	if tools.ElementExists(db, "PROVIDER", "name", name) {
		tools.JsonResponse(w, 400, `{"message": "Provider already exists"}`)
		return
	}

	if tools.ElementExists(db, "PROVIDER", "email", email) {
		tools.JsonResponse(w, 400, `{"message": "Email already exists"}`)
		return
	}

	uuid := tools.GenerateUUID()

	// Inserting the provider in the database
	_, err := tools.ExecuteQuery(db, "INSERT INTO `PROVIDER` (`id`, `name`, `email`) VALUES (?, ?, ?)", uuid, name, email)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Creating the response
	jsonResponse := `{"message": "Provider created"`

	// Adding the return fields of the query
	fields, err := ProviderGetAll(db, uuid, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func ProviderGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, "name", "id", "email", "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `id`, `name`, `email` FROM `PROVIDER`"
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
	jsonResponse, err := ProviderGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ProviderPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, "name", "email") || tools.ValuesNotInQuery(query, "id") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	id := query["id"]
	name := tools.BodyValueToString(body, "name")
	email := tools.BodyValueToString(body, "email")

	// Checking if the values are empty
	if tools.ValueIsEmpty(id) {
		tools.JsonResponse(w, 400, `{"message": "Empty fields"}`)
		return
	}

	for key, _ := range body {
		if tools.ValueIsEmpty(tools.BodyValueToString(body, key)) {
			tools.JsonResponse(w, 400, `{"message": "Empty fields"}`)
			return
		}
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(8, name) {
		tools.JsonResponse(w, 400, `{"message": "Values too short"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooLong(64, name) {
		tools.JsonResponse(w, 400, `{"message": "Values too long"}`)
		return
	}

	if !tools.ValueIsEmpty(email) {
		if tools.EmailIsValid(email) == false {
			tools.JsonResponse(w, 400, `{"message": "Invalid email"}`)
			return
		}
	}

	if !tools.ElementExists(db, "PROVIDER", "id", id) {
		tools.JsonResponse(w, 400, `{"message": "Provider does not exist"}`)
		return
	}

	if tools.ElementExists(db, "PROVIDER", "name", name) {
		tools.JsonResponse(w, 400, `{"message": "Provider already exists"}`)
		return
	}

	if tools.ElementExists(db, "PROVIDER", "email", email) {
		tools.JsonResponse(w, 400, `{"message": "Email already exists"}`)
		return
	}

	request := "UPDATE `PROVIDER` SET "
	var params []interface{}
	
	for key, value := range body {
		if key != "id" {
			tools.AppendUpdate(&request, &params, key, value)
		}
	}

	// Removing the last ","
	request = request[:len(request)-2]

	request += " WHERE `id` = ?"
	params = append(params, id)

	// Updating the provider in the database
	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Provider updated"`
	
	// Adding the return fields of the query
	fields, err := ProviderGetAll(db, id, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func ProviderDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, "id") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	id := query["id"]

	// Checking if the provider exists
	if !tools.ElementExists(db, "PROVIDER", "id", id) {
		tools.JsonResponse(w, 400, `{"message": "Provider does not exist"}`)
		return
	}

	// Deleting the provider in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `PROVIDER` WHERE `id` = ?", id)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Provider deleted", "id": "` + id + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ProviderGetAll(db *sql.DB, uuid string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `id`, `name`, `email` FROM `PROVIDER` WHERE `id` = ?", uuid)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return ProviderGetAllAssociation(result, arrayOutput)
}

func ProviderGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var id, name, email string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&id, &name, &email)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"id": "` + id + `", "name": "` + name + `", "email": "` + email + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&id, &name, &email)
			if err != nil {
				return "", err
			}
		}
		return `"id": "` + id + `", "name": "` + name + `", "email": "` + email + `"`, nil
	}
}