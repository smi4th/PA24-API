package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Subscription(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		SubscriptionPost(w, r, db)
	case "GET":
		SubscriptionGet(w, r, db)
	case "PUT":
		SubscriptionPut(w, r, db)
	case "DELETE":
		SubscriptionDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func SubscriptionPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, "name") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	name := tools.BodyValueToString(body, "name")

	// Checking if the values are empty
	if tools.ValueIsEmpty(name) {
		tools.JsonResponse(w, 400, `{"message": "Name cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, name) {
		tools.JsonResponse(w, 400, `{"message": "Name too short"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooLong(32, name) {
		tools.JsonResponse(w, 400, `{"message": "Name too long"}`)
		return
	}

	// Checking if the subscription type is valid
	if tools.ElementExists(db, "SUBSCRIPTION", "name", name) {
		tools.JsonResponse(w, 400, `{"message": "Subscription already exists"}`)
		return
	}

	uuid := tools.GenerateUUID()

	// Inserting the account in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `SUBSCRIPTION` (`id`, `name`) VALUES (?, ?)", uuid, name)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Subscription created"`

	// Adding the return fields of the query
	fields, err := SubscriptionGetAll(db, uuid, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	jsonResponse += "," + fields

	tools.InfoLog(tools.RowsToJson(result))

	jsonResponse += "}"

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse)

}

func SubscriptionGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, "id", "name", "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `id`, `name` FROM `SUBSCRIPTION`"
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
	jsonResponse, err := SubscriptionGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func SubscriptionPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, "name") || tools.ValuesNotInQuery(query, "id") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	id := query["id"]
	name := tools.BodyValueToString(body, "name")

	// Checking if the values are empty
	if tools.ValueIsEmpty(id, name) {
		tools.JsonResponse(w, 400, `{"message": "ID and name cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, name) {
		tools.JsonResponse(w, 400, `{"message": "Name too short"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooLong(32, name) {
		tools.JsonResponse(w, 400, `{"message": "Name too long"}`)
		return
	}

	// Checking if the subscription exists
	if !tools.ElementExists(db, "SUBSCRIPTION", "id", id) {
		tools.JsonResponse(w, 400, `{"message": "Subscription does not exist"}`)
		return
	}

	// Checking if the subscription type is valid
	if tools.ElementExists(db, "SUBSCRIPTION", "name", name) {
		tools.JsonResponse(w, 400, `{"message": "Subscription already exists"}`)
		return
	}

	// Updating the subscription in the database
	result, err := tools.ExecuteQuery(db, "UPDATE `SUBSCRIPTION` SET `name` = ? WHERE `id` = ?", name, id)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Subscription updated"`

	// Adding the return fields of the query
	fields, err := SubscriptionGetAll(db, id, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	jsonResponse += "," + fields + "}"

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func SubscriptionDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, "id") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	id := query["id"]

	// Checking if the subscription exists
	if !tools.ElementExists(db, "SUBSCRIPTION", "id", id) {
		tools.JsonResponse(w, 400, `{"message": "Subscription does not exist"}`)
		return
	}

	// Deleting the subscription in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `SUBSCRIPTION` WHERE `id` = ?", id)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Subscription deleted", "id": "` + id + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func SubscriptionGetAll(db *sql.DB, uuid string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `id`, `name` FROM `SUBSCRIPTION` WHERE `id` = ?", uuid)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return SubscriptionGetAllAssociation(result, arrayOutput)
}

func SubscriptionGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var id, name string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&id, &name)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"id": "` + id + `", "name": "` + name + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&id, &name)
			if err != nil {
				return "", err
			}
		}
		return `"id": "` + id + `", "name": "` + name + `"`, nil
	}
}