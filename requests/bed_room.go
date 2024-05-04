package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func BedRoom(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		BedRoomPost(w, r, db)
	case "GET":
		BedRoomGet(w, r, db)
	case "PUT":
		if tools.GetUUID(r, db) == tools.GetElement(db, "HOUSING", "account", "uuid", tools.GetElement(db, "BED_ROOM", "housing", "uuid", tools.ReadQuery(r)["uuid"])) || tools.IsAdmin(r, db) {
			BedRoomPut(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	case "DELETE":
		if tools.GetUUID(r, db) == tools.GetElement(db, "HOUSING", "account", "uuid", tools.GetElement(db, "BED_ROOM", "housing", "uuid", tools.ReadQuery(r)["uuid"])) || tools.IsAdmin(r, db) {
			BedRoomDelete(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func BedRoomPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `nbPlaces`, `price`, `description`, `housing`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    nbPlaces_ := tools.BodyValueToString(body, "nbPlaces")
	price_ := tools.BodyValueToString(body, "price")
	description_ := tools.BodyValueToString(body, "description")
	housing_ := tools.BodyValueToString(body, "housing")

	if tools.GetUUID(r, db) != tools.GetElement(db, "HOUSING", "account", "uuid", housing_) && !tools.IsAdmin(r, db) && !tools.IsAdmin(r, db) {
		tools.JsonResponse(w, 403, `{"error": "Forbidden"}`)
		return
	}
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(nbPlaces_, price_, description_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, description_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}

    if !tools.ValueIsEmpty(housing_) {
		if !tools.ElementExists(db, "HOUSING", "uuid", housing_) {
			tools.JsonResponse(w, 400, `{"error": "This housing does not exist"}`) 
			return
		}
	}
	
	
	

	

	

	uuid_ := tools.GenerateUUID()

	// Inserting the BedRoom in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `BED_ROOM` (`uuid`, `nbPlaces`, `price`, `description`, `validated`, `housing`) VALUES (?, ?, ?, ?, false, ?)", uuid_, nbPlaces_, price_, description_, housing_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "BedRoom created"`

	// Adding the return fields of the query
	fields, err := BedRoomGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func BedRoomGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`, `nbPlaces`, `price`, `description`, `validated`, `housing`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `nbPlaces`, `price`, `description`, `validated`, `housing` FROM `BED_ROOM`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `BED_ROOM`"
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
	jsonResponse, err := BedRoomGetAllAssociation(result, true)
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

func BedRoomPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `nbPlaces`, `price`, `description`, `validated`, `housing`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	
	housing_ := tools.BodyValueToString(body, "housing")
	

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

    if !tools.ValueIsEmpty(housing_) {
		if !tools.ElementExists(db, "HOUSING", "uuid", housing_) {
			tools.JsonResponse(w, 400, `{"error": "This housing does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExists(db, "BED_ROOM", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This BedRoom does not exist"}`) 
		return
	}
	

	

    

	request := "UPDATE `BED_ROOM` SET "
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
	jsonResponse := `{"message": "BedRoom updated"`
	
	// Adding the return fields of the query
	fields, err := BedRoomGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func BedRoomDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "BED_ROOM", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This BedRoom does not exist"}`) 
		return
	}
	

	// Deleting the BedRoom in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `BED_ROOM` WHERE uuid = ?", uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "BedRoom deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func BedRoomGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `nbPlaces`, `price`, `description`, `validated`, `housing` FROM `BED_ROOM` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return BedRoomGetAllAssociation(result, arrayOutput)
}

func BedRoomGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, nbPlaces_, price_, description_, validated_, housing_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &nbPlaces_, &price_, &description_, &validated_, &housing_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "nbPlaces": "` + nbPlaces_ + `", "price": "` + price_ + `", "description": "` + description_ + `", "validated": "` + validated_ + `", "housing": "` + housing_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &nbPlaces_, &price_, &description_, &validated_, &housing_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "nbPlaces": "` + nbPlaces_ + `", "price": "` + price_ + `", "description": "` + description_ + `", "validated": "` + validated_ + `", "housing": "` + housing_ + `"`, nil
	}
}