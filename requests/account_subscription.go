package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func AccountSubscription(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		AccountSubscriptionPost(w, r, db)
	case "GET":
		AccountSubscriptionGet(w, r, db)
	case "PUT":
		AccountSubscriptionPut(w, r, db)
	case "DELETE":
		AccountSubscriptionDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func AccountSubscriptionPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `start_date`, `account`, `subscription`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    start_date_ := tools.BodyValueToString(body, "start_date")
	account_ := tools.BodyValueToString(body, "account")
	subscription_ := tools.BodyValueToString(body, "subscription")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(start_date_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, start_date_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, start_date_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

    if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(subscription_) {
		if !tools.ElementExists(db, "SUBSCRIPTION", "uuid", subscription_) {
			tools.JsonResponse(w, 400, `{"error": "This subscription does not exist"}`) 
			return
		}
	}
	
	
	

	

	

	

	// Inserting the AccountSubscription in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `ACCOUNT_SUBSCRIPTION` (`account`, `subscription`, `start_date`, `account`, `subscription`) VALUES (?, ?, ?)", account_, subscription_, start_date_, account_, subscription_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "AccountSubscription created"`

	// Adding the return fields of the query
	fields, err := AccountSubscriptionGetAll(db, account_, subscription_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func AccountSubscriptionGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `start_date`, `account`, `subscription`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `start_date`, `account`, `subscription` FROM `ACCOUNT_SUBSCRIPTION`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `ACCOUNT_SUBSCRIPTION`"
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
	jsonResponse, err := AccountSubscriptionGetAllAssociation(result, true)
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

func AccountSubscriptionPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `start_date`) || tools.ValuesNotInQuery(query, `account`, `subscription`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]
	subscription_ := query["subscription"]
	
    start_date_ := tools.BodyValueToString(body, "start_date")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(account_, subscription_) {
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
	if tools.ValueTooShort(4, start_date_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, start_date_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

    if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(subscription_) {
		if !tools.ElementExists(db, "SUBSCRIPTION", "uuid", subscription_) {
			tools.JsonResponse(w, 400, `{"error": "This subscription does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExists(db, "ACCOUNT_SUBSCRIPTION", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountSubscription does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "ACCOUNT_SUBSCRIPTION", "subscription", subscription_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountSubscription does not exist"}`) 
		return
	}
	

	

    

	request := "UPDATE `ACCOUNT_SUBSCRIPTION` SET "
	var params []interface{}
	
	for key, value := range body {
		if !tools.ValueInArray(key, `account`, `subscription`) {
			if key == "password" {
				value = tools.HashPassword(value.(string))
			}
			tools.AppendUpdate(&request, &params, key, value)
		}
	}

	// Removing the last ","
	request = request[:len(request)-2]

	request += " WHERE account = ?, subscription = ?"
	params = append(params, account_, subscription_)

	// Updating the account in the database
	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "AccountSubscription updated"`
	
	// Adding the return fields of the query
	fields, err := AccountSubscriptionGetAll(db, account_, subscription_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func AccountSubscriptionDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `account`, `subscription`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]
	subscription_ := query["subscription"]
	

	if !tools.ElementExists(db, "ACCOUNT_SUBSCRIPTION", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountSubscription does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "ACCOUNT_SUBSCRIPTION", "subscription", subscription_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountSubscription does not exist"}`) 
		return
	}
	

	// Deleting the AccountSubscription in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `ACCOUNT_SUBSCRIPTION` WHERE account = ?, subscription = ?", account_, subscription_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "AccountSubscription deleted", "account": "` + account_ + `", "subscription": "` + subscription_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func AccountSubscriptionGetAll(db *sql.DB, account_ string, subscription_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `start_date`, `account`, `subscription` FROM `ACCOUNT_SUBSCRIPTION` WHERE account = ?, subscription = ?", account_, subscription_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return AccountSubscriptionGetAllAssociation(result, arrayOutput)
}

func AccountSubscriptionGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var start_date_, account_, subscription_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&start_date_, &account_, &subscription_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"start_date": "` + start_date_ + `", "account": "` + account_ + `", "subscription": "` + subscription_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&start_date_, &account_, &subscription_)
			if err != nil {
				return "", err
			}
		}
		return `"start_date": "` + start_date_ + `", "account": "` + account_ + `", "subscription": "` + subscription_ + `"`, nil
	}
}