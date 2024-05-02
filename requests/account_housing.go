package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func AccountHousing(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		AccountHousingPost(w, r, db)
	case "GET":
		AccountHousingGet(w, r, db)
	case "PUT":
		AccountHousingPut(w, r, db)
	case "DELETE":
		AccountHousingDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func AccountHousingPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `account`, `housing`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    account_ := tools.BodyValueToString(body, "account")
	housing_ := tools.BodyValueToString(body, "housing")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty() {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, ) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, ) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

    if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(housing_) {
		if !tools.ElementExists(db, "HOUSING", "uuid", housing_) {
			tools.JsonResponse(w, 400, `{"error": "This housing does not exist"}`) 
			return
		}
	}
	
	
	

	

	

	

	// Inserting the AccountHousing in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `ACCOUNT_HOUSING` (`account`, `housing`, `account`, `housing`) VALUES (?, ?)", account_, housing_, account_, housing_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "AccountHousing created"`

	// Adding the return fields of the query
	fields, err := AccountHousingGetAll(db, account_, housing_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func AccountHousingGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `creation_date`, `account`, `housing`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `creation_date`, `account`, `housing` FROM `ACCOUNT_HOUSING`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `ACCOUNT_HOUSING`"
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
	jsonResponse, err := AccountHousingGetAllAssociation(result, true)
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

func AccountHousingPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, ``) || tools.ValuesNotInQuery(query, `account`, `housing`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]
	housing_ := query["housing"]
	
    

	// Checking if the values are empty
	if tools.ValueIsEmpty(account_, housing_) {
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
	if tools.ValueTooShort(4, ) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, ) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

    if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(housing_) {
		if !tools.ElementExists(db, "HOUSING", "uuid", housing_) {
			tools.JsonResponse(w, 400, `{"error": "This housing does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExists(db, "ACCOUNT_HOUSING", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountHousing does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "ACCOUNT_HOUSING", "housing", housing_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountHousing does not exist"}`) 
		return
	}
	

	

    

	request := "UPDATE `ACCOUNT_HOUSING` SET "
	var params []interface{}
	
	for key, value := range body {
		if !tools.ValueInArray(key, `account`, `housing`) {
			if key == "password" {
				value = tools.HashPassword(value.(string))
			}
			tools.AppendUpdate(&request, &params, key, value)
		}
	}

	// Removing the last ","
	request = request[:len(request)-2]

	request += " WHERE account = ?, housing = ?"
	params = append(params, account_, housing_)

	// Updating the account in the database
	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "AccountHousing updated"`
	
	// Adding the return fields of the query
	fields, err := AccountHousingGetAll(db, account_, housing_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func AccountHousingDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `account`, `housing`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]
	housing_ := query["housing"]
	

	if !tools.ElementExists(db, "ACCOUNT_HOUSING", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountHousing does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "ACCOUNT_HOUSING", "housing", housing_) {
		tools.JsonResponse(w, 400, `{"error": "This AccountHousing does not exist"}`) 
		return
	}
	

	// Deleting the AccountHousing in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `ACCOUNT_HOUSING` WHERE account = ?, housing = ?", account_, housing_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "AccountHousing deleted", "account": "` + account_ + `", "housing": "` + housing_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func AccountHousingGetAll(db *sql.DB, account_ string, housing_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `creation_date`, `account`, `housing` FROM `ACCOUNT_HOUSING` WHERE account = ?, housing = ?", account_, housing_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return AccountHousingGetAllAssociation(result, arrayOutput)
}

func AccountHousingGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var creation_date_, account_, housing_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&creation_date_, &account_, &housing_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"creation_date": "` + creation_date_ + `", "account": "` + account_ + `", "housing": "` + housing_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&creation_date_, &account_, &housing_)
			if err != nil {
				return "", err
			}
		}
		return `"creation_date": "` + creation_date_ + `", "account": "` + account_ + `", "housing": "` + housing_ + `"`, nil
	}
}