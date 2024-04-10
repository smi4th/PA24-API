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
	if tools.ValuesNotInBody(body, "provider", "account") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	// administration_level := tools.BodyValueToString(body, "administration_level")
	provider := tools.BodyValueToString(body, "provider")
	account := tools.BodyValueToString(body, "account")

	// Checking if the values are empty
	if tools.ValueIsEmpty(provider, account) {
		tools.JsonResponse(w, 400, `{"message": "Empty fields"}`)
		return
	}

	// Checking if the values are in theyr respective tables
	// if !tools.ElementExists(db, "ADMINISTRATION_LEVEL", "id", administration_level) {
	// 	tools.JsonResponse(w, 400, `{"message": "Administration level does not exist"}`)
	// 	return
	// }

	if !tools.ElementExists(db, "PROVIDER", "id", provider) {
		tools.JsonResponse(w, 400, `{"message": "Provider does not exist"}`)
		return
	}

	if !tools.ElementExists(db, "ACCOUNT", "id", account) {
		tools.JsonResponse(w, 400, `{"message": "Account does not exist"}`)
		return
	}

	// Checking if the provider is already associated with the account
	if tools.ElementExistsInLinkTable(db, "PROVIDER_ACCOUNT", "provider", provider, "account", account) {
		tools.JsonResponse(w, 400, `{"message": "Provider already associated with the account"}`)
		return
	}

	// Inserting the provider in the database
	_, err := tools.ExecuteQuery(db, "INSERT INTO `PROVIDER_ACCOUNT` (`provider`, `account`) VALUES (?, ?)", provider, account)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Creating the response
	jsonResponse := `{"message": "ProviderAccount created"`

	// Adding the return fields of the query
	fields, err := ProviderAccountGetAll(db, provider, account, false)
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
	if tools.AtLeastOneValueInQuery(query, "provider", "account", "all") {
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

func ProviderAccountDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, "provider", "account") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	provider := query["provider"]
	account := query["account"]

	if !tools.ElementExists(db, "PROVIDER", "id", provider) {
		tools.JsonResponse(w, 400, `{"message": "Provider does not exist"}`)
		return
	}

	if !tools.ElementExists(db, "ACCOUNT", "id", account) {
		tools.JsonResponse(w, 400, `{"message": "Account does not exist"}`)
		return
	}

	// Checking if the provider exists
	if !tools.ElementExistsInLinkTable(db, "PROVIDER_ACCOUNT", "provider", provider, "account", account) {
		tools.JsonResponse(w, 400, `{"message": "Provider not associated with the account"}`)
		return
	}

	// Deleting the provider in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `PROVIDER_ACCOUNT` WHERE `provider` = ? AND `account` = ?", provider, account)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "ProviderAccount deleted", "provider": "` + provider + `", "account": "` + account + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ProviderAccountGetAll(db *sql.DB, provider, account string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `administration_level`, `provider`, `account` FROM `PROVIDER_ACCOUNT` WHERE `provider` = ? AND `account` = ?", provider, account)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return ProviderAccountGetAllAssociation(result, arrayOutput)
}

func ProviderAccountGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var provider, account, administration_level string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&administration_level, &provider, &account)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"administration_level" : "` + administration_level + `", "provider": "` + provider + `", "account": "` + account + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&administration_level, &provider, &account)
			if err != nil {
				return "", err
			}
		}
		return `"administration_level" : "` + administration_level + `", "provider": "` + provider + `", "account": "` + account + `"`, nil
	}
}