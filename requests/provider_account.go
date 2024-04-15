package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func ProviderAccount(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		ProviderAccountPost(w, r, db)
	case "GET":
		ProviderAccountGet(w, r, db)
	case "PUT":
		ProviderAccountPut(w, r, db)
	case "DELETE":
		ProviderAccountDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func ProviderAccountPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `administration_level`, `provider`, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    administration_level_ := tools.BodyValueToString(body, "administration_level")
	provider_ := tools.BodyValueToString(body, "provider")
	account_ := tools.BodyValueToString(body, "account")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(administration_level_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, administration_level_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, administration_level_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

    if !tools.ValueIsEmpty(provider_) {
		if !tools.ElementExists(db, "PROVIDER", "uuid", provider_) {
			tools.JsonResponse(w, 400, `{"error": "This provider does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	
	
	

	

	

	

	// Inserting the ProviderAccount in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `PROVIDER_ACCOUNT` (`provider`, `account`, `administration_level`, `provider`, `account`) VALUES (?, ?, ?)", provider_, account_, administration_level_, provider_, account_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "ProviderAccount created"`

	// Adding the return fields of the query
	fields, err := ProviderAccountGetAll(db, provider_, account_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func ProviderAccountGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `administration_level`, `provider`, `account`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `administration_level`, `provider`, `account` FROM `PROVIDER_ACCOUNT`"
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
	jsonResponse, err := ProviderAccountGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ProviderAccountPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `administration_level`) || tools.ValuesNotInQuery(query, `provider`, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	provider_ := query["provider"]
	account_ := query["account"]
	
    administration_level_ := tools.BodyValueToString(body, "administration_level")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(provider_, account_) {
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
	if tools.ValueTooShort(4, administration_level_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, administration_level_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

    if !tools.ValueIsEmpty(provider_) {
		if !tools.ElementExists(db, "PROVIDER", "uuid", provider_) {
			tools.JsonResponse(w, 400, `{"error": "This provider does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(account_) {
		if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
			tools.JsonResponse(w, 400, `{"error": "This account does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExists(db, "PROVIDER_ACCOUNT", "provider", provider_) {
		tools.JsonResponse(w, 400, `{"error": "This ProviderAccount does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "PROVIDER_ACCOUNT", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This ProviderAccount does not exist"}`) 
		return
	}
	

	

    

	request := "UPDATE `PROVIDER_ACCOUNT` SET "
	var params []interface{}
	
	for key, value := range body {
		if !tools.ValueInArray(key, `provider`, `account`) {
			if key == "password" {
				value = tools.HashPassword(value.(string))
			}
			tools.AppendUpdate(&request, &params, key, value)
		}
	}

	// Removing the last ","
	request = request[:len(request)-2]

	request += " WHERE provider = ?, account = ?"
	params = append(params, provider_, account_)

	// Updating the account in the database
	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "ProviderAccount updated"`
	
	// Adding the return fields of the query
	fields, err := ProviderAccountGetAll(db, provider_, account_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func ProviderAccountDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `provider`, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	provider_ := query["provider"]
	account_ := query["account"]
	

	if !tools.ElementExists(db, "PROVIDER_ACCOUNT", "provider", provider_) {
		tools.JsonResponse(w, 400, `{"error": "This ProviderAccount does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "PROVIDER_ACCOUNT", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This ProviderAccount does not exist"}`) 
		return
	}
	

	// Deleting the ProviderAccount in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `PROVIDER_ACCOUNT` WHERE provider = ?, account = ?", provider_, account_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "ProviderAccount deleted", "provider": "` + provider_ + `", "account": "` + account_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ProviderAccountGetAll(db *sql.DB, provider_ string, account_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `administration_level`, `provider`, `account` FROM `PROVIDER_ACCOUNT` WHERE provider = ?, account = ?", provider_, account_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return ProviderAccountGetAllAssociation(result, arrayOutput)
}

func ProviderAccountGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var administration_level_, provider_, account_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&administration_level_, &provider_, &account_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"administration_level": "` + administration_level_ + `", "provider": "` + provider_ + `", "account": "` + account_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&administration_level_, &provider_, &account_)
			if err != nil {
				return "", err
			}
		}
		return `"administration_level": "` + administration_level_ + `", "provider": "` + provider_ + `", "account": "` + account_ + `"`, nil
	}
}