package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Review(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		ReviewPost(w, r, db)
	case "GET":
		ReviewGet(w, r, db)
	case "PUT":
		if tools.GetUUID(r, db) == tools.GetElement(db, "REVIEW", "account", "uuid", tools.ReadQuery(r)["uuid"]) || tools.IsAdmin(r, db) {
			ReviewPut(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	case "DELETE":
		if tools.GetUUID(r, db) == tools.GetElement(db, "REVIEW", "account", "uuid", tools.ReadQuery(r)["uuid"]) || tools.IsAdmin(r, db) {
			ReviewDelete(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func ReviewPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `content`, `note`, `account`) || tools.AtLeastOneValueInBody(body, `service`, `housing`, `bedRoom`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    content_ := tools.BodyValueToString(body, "content")
	note_ := tools.BodyValueToString(body, "note")
	account_ := tools.BodyValueToString(body, "account")
	services_ := tools.BodyValueToString(body, "service")
	housing_ := tools.BodyValueToString(body, "housing")
	bedRoom_ := tools.BodyValueToString(body, "bedRoom")

	if tools.GetUUID(r, db) != account_ && !tools.IsAdmin(r, db) {
		tools.JsonResponse(w, 403, `{"error": "Forbidden"}`)
		return
	}
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(content_, note_, account_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, content_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}

	requestValue := ""
	request := ""

	if !tools.ValueIsEmpty(services_) {
		tools.InfoLog("Value not empty")
		if !tools.ElementExists(db, "SERVICES", "uuid", services_) {
			tools.JsonResponse(w, 400, `{"error": "This services does not exist"}`)
			return
		} else {
			tools.InfoLog("Value exists")
			request = "service"
			requestValue = services_
		}
	} else if !tools.ValueIsEmpty(housing_) {
		if !tools.ElementExists(db, "HOUSING", "uuid", housing_) {
			tools.JsonResponse(w, 400, `{"error": "This housing does not exist"}`)
			return
		} else {
			request = "housing"
			requestValue = housing_
		}
	} else if !tools.ValueIsEmpty(bedRoom_) {
		if !tools.ElementExists(db, "BED_ROOM", "uuid", bedRoom_) {
			tools.JsonResponse(w, 400, `{"error": "This bedRoom does not exist"}`)
			return
		} else {
			request = "bedRoom"
			requestValue = bedRoom_
		}
	} else {
		tools.JsonResponse(w, 400, `{"error": "No element found"}`)
		return
	}

	if tools.ElementExists(db, "REVIEW", "account", account_) {
		tools.JsonResponse(w, 400, `{"error": "This account already posted a review on this element"}`)
		return
	}


	tools.InfoLog(request)

	// Generating the UUID
	uuid_ := tools.GenerateUUID()

	// Inserting the Review in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `REVIEW` (`uuid`, `content`, `note`, `account`, `" + request + "`) VALUES (?, ?, ?, ?, ?)", uuid_, content_, note_, account_, requestValue)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Review created"`

	// Adding the return fields of the query
	fields, err := ReviewGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func ReviewGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`, `account`, `services`, `housing`, `bedRoom`, `note`, `content`, `all`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `content`, `note`, `account`, case when `service` is null then 'NULL' else `service` end as `service`, case when `housing` is null then 'NULL' else `housing` end as `housing`, case when `bedRoom` is null then 'NULL' else `bedRoom` end as `bedRoom` FROM `REVIEW`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `REVIEW`"
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
	jsonResponse, err := ReviewGetAllAssociation(result, true)
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

func ReviewPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `content`, `note`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	// check if service, housing or bedRoom is in the body
	if !tools.AtLeastOneValueInBody(body, `service`, `housing`, `bedRoom`, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Cannot update the element"}`)
		return
	}

	uuid_ := query["uuid"]
	
	content_ := tools.BodyValueToString(body, "content")

	// Checking if the values are empty
	if tools.ValueIsEmpty(uuid_) {
		tools.JsonResponse(w, 400, `{"message": "Empty fields"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, content_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
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

	// Creating the request
	request := "UPDATE `REVIEW` SET "
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
	jsonResponse := `{"message": "Review updated"`
	
	// Adding the return fields of the query
	fields, err := ReviewGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func ReviewDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "REVIEW", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Review does not exist"}`) 
		return
	}
	

	// Deleting the Review in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `REVIEW` WHERE uuid = ?", uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Review deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func ReviewGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `content`, `note`, `account`, case when `service` is null then 'NULL' else `service` end as `service`, case when `housing` is null then 'NULL' else `housing` end as `housing`, case when `bedRoom` is null then 'NULL' else `bedRoom` end as `bedRoom` FROM `REVIEW` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return ReviewGetAllAssociation(result, arrayOutput)
}

func ReviewGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, account_, content_, note_, services_, housing_, bedRoom_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &content_, &note_, &account_, &services_, &housing_, &bedRoom_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "account": "` + account_ + `", "content": "` + content_ + `", "note": "` + note_ + `", "services": "` + services_ + `", "housing": "` + housing_ + `", "bedRoom": "` + bedRoom_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &content_, &note_, &account_, &services_, &housing_, &bedRoom_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "account": "` + account_ + `", "content": "` + content_ + `", "note": "` + note_ + `", "services": "` + services_ + `", "housing": "` + housing_ + `", "bedRoom": "` + bedRoom_ + `"`, nil
	}
}