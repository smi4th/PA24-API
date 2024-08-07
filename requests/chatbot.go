package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Chatbot(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		ChatbotPost(w, r, db)
	case "GET":
		ChatbotGet(w, r, db)
	case "PUT":
		ChatbotPut(w, r, db)
	case "DELETE":
		ChatbotDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func ChatbotPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `keyword`, `text`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    keyword_ := tools.BodyValueToString(body, "keyword")
	text_ := tools.BodyValueToString(body, "text")

	// Checking if the values are empty
	if tools.ValueIsEmpty(keyword_, text_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}
	
	if tools.ValueTooLong(255, keyword_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}
	

	

	if tools.ElementExists(db, "CHATBOT", "keyword", keyword_) {
		tools.JsonResponse(w, 400, `{"error": "This keyword already exists"}`)
		return
	}

	uuid_ := tools.GenerateUUID()

	// Inserting the Chatbot in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `CHATBOT` (`uuid`, `keyword`, `text`) VALUES (?, ?, ?)", uuid_, keyword_, text_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Chatbot created"`

	// Adding the return fields of the query
	fields, err := ChatbotGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func ChatbotGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`, `keyword`, `text`, `all`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `keyword`, `text` FROM `CHATBOT`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `CHATBOT`"
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
	jsonResponse, err := ChatbotGetAllAssociation(result, true)
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

func ChatbotPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `keyword`, `text`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	if !tools.AtLeastOneValueInBody(body, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Empty fields"}`)
		return
	}

	uuid_ := query["uuid"]
	
	keyword_ := tools.BodyValueToString(body, "keyword")
	

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

	if tools.ValueTooLong(255, keyword_) {
		tools.JsonResponse(w, 400, `{"message": "keyword too long"}`)
		return
	}
    

	if !tools.ElementExists(db, "CHATBOT", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Chatbot does not exist"}`) 
		return
	}
	if tools.ElementExists(db, "CHATBOT", "keyword", keyword_) {
		tools.JsonResponse(w, 400, `{"error": "This keyword already exists"}`)
		return
	}
	

	

    

	request := "UPDATE `CHATBOT` SET "
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
	jsonResponse := `{"message": "Chatbot updated"`
	
	// Adding the return fields of the query
	fields, err := ChatbotGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func ChatbotDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "CHATBOT", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Chatbot does not exist"}`) 
		return
	}
	

	// Deleting the Chatbot in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `CHATBOT` WHERE uuid = ?", uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Chatbot deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ChatbotGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `keyword`, `text` FROM `CHATBOT` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return ChatbotGetAllAssociation(result, arrayOutput)
}

func ChatbotGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, keyword_, text_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &keyword_, &text_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "keyword": "` + keyword_ + `", "text": "` + text_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &keyword_, &text_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "keyword": "` + keyword_ + `", "text": "` + text_ + `"`, nil
	}
}