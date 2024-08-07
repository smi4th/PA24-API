package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Subscription(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		if tools.IsAdmin(r, db) {
			SubscriptionPost(w, r, db)
		} else {
			tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		}
	case "GET":
		SubscriptionGet(w, r, db)
	case "PUT":
		if tools.IsAdmin(r, db) {
			SubscriptionPut(w, r, db)
		} else {
			tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		}
	case "DELETE":
		if tools.IsAdmin(r, db) {
			SubscriptionDelete(w, r, db)
		} else {
			tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		}
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func SubscriptionPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `name`, `price`, `ads`, `VIP`, `description`, `duration`, `imgPath`, `taxes`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    name_ := tools.BodyValueToString(body, "name")
	price_ := tools.BodyValueToString(body, "price")
	ads_ := tools.BodyValueToString(body, "ads")
	VIP_ := tools.BodyValueToString(body, "VIP")
	description_ := tools.BodyValueToString(body, "description")
	duration_ := tools.BodyValueToString(body, "duration")
	imgPath_ := tools.BodyValueToString(body, "imgPath")
	taxes_ := tools.BodyValueToString(body, "taxes")	

	// Checking if the values are empty
	if tools.ValueIsEmpty(name_, price_, ads_, VIP_, description_, duration_, imgPath_, taxes_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, name_, description_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, name_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}
	

	

	if tools.ElementExists(db, "SUBSCRIPTION", "name", name_) {
		tools.JsonResponse(w, 400, `{"error": "This name already exists"}`) 
		return
	}
	
	if tools.ElementExists(db, "TAXES", "uuid", taxes_) {
		tools.JsonResponse(w, 400, `{"error": "This taxes already exists"}`)
		return
	}

	uuid_ := tools.GenerateUUID()

	// Inserting the Subscription in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `SUBSCRIPTION` (`uuid`, `name`, `price`, `ads`, `VIP`, `description`, `duration`, `imgPath`, `taxes`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", uuid_, name_, price_, ads_, VIP_, description_, duration_, imgPath_, taxes_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Subscription created"`

	// Adding the return fields of the query
	fields, err := SubscriptionGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func SubscriptionGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`, `name`, `price`, `ads`, `VIP`, `description`, `duration`, `all`, `imgPath`, `taxes`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `name`, `price`, `ads`, `VIP`, `description`, `duration`, `imgPath` FROM `SUBSCRIPTION`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `SUBSCRIPTION`"
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
	jsonResponse, err := SubscriptionGetAllAssociation(result, true)
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

func SubscriptionPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `name`, `price`, `ads`, `VIP`, `description`, `duration`, `imgPath`, `taxes`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	if !tools.AtLeastOneValueInBody(body, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Empty fields"}`)
		return
	}

	uuid_ := query["uuid"]
	
    name_ := tools.BodyValueToString(body, "name")
	

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
	if tools.ValueTooShort(4, name_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, name_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}
    

	if !tools.ElementExists(db, "SUBSCRIPTION", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Subscription does not exist"}`) 
		return
	}
	if tools.ElementExists(db, "SUBSCRIPTION", "name", name_) {
		tools.JsonResponse(w, 400, `{"error": "This name already exists"}`) 
		return
	}
	

	

    

	request := "UPDATE `SUBSCRIPTION` SET "
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
	jsonResponse := `{"message": "Subscription updated"`
	
	// Adding the return fields of the query
	fields, err := SubscriptionGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func SubscriptionDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "SUBSCRIPTION", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Subscription does not exist"}`) 
		return
	}
	

	// Deleting the Subscription in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `SUBSCRIPTION` WHERE uuid = ?", uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Subscription deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func SubscriptionGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `name`, `price`, `ads`, `VIP`, `description`, `duration`, `imgPath`, `taxes` FROM `SUBSCRIPTION` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return SubscriptionGetAllAssociation(result, arrayOutput)
}

func SubscriptionGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, name_, price_, ads_, VIP_, description_, duration_, imgPath_, taxes_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &name_, &price_, &ads_, &VIP_, &description_, &duration_, &imgPath_, &taxes_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "name": "` + name_ + `", "price": "` + price_ + `", "ads": "` + ads_ + `", "VIP": "` + VIP_ + `", "description": "` + description_ + `", "duration": "` + duration_ + `", "imgPath": "` + imgPath_ + `", "taxes": "` + taxes_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &name_, &price_, &ads_, &VIP_, &description_, &duration_, &imgPath_, &taxes_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "name": "` + name_ + `", "price": "` + price_ + `", "ads": "` + ads_ + `", "VIP": "` + VIP_ + `", "description": "` + description_ + `", "duration": "` + duration_ + `", "imgPath": "` + imgPath_ + `", "taxes": "` + taxes_ + `"`, nil
	}
}