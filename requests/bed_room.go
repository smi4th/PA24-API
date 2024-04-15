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
		BedRoomPut(w, r, db)
	case "DELETE":
		BedRoomDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func BedRoomPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `nbPlaces`, `price`, `description`, `validated`, `housing`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    nbPlaces_ := tools.BodyValueToString(body, "nbPlaces")
	price_ := tools.BodyValueToString(body, "price")
	description_ := tools.BodyValueToString(body, "description")
	validated_ := tools.BodyValueToString(body, "validated")
	housing_ := tools.BodyValueToString(body, "housing")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(nbPlaces_, price_, description_, validated_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, nbPlaces_, price_, description_, validated_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, nbPlaces_, price_, description_, validated_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
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
	result, err := tools.ExecuteQuery(db, "INSERT INTO `BED_ROOM` (`uuid`, `nbPlaces`, `price`, `description`, `validated`, `housing`) VALUES (?, ?, ?, ?, ?, ?)", uuid_, nbPlaces_, price_, description_, validated_, housing_)
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
	jsonResponse, err := BedRoomGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

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
	
    nbPlaces_ := tools.BodyValueToString(body, "nbPlaces")
	price_ := tools.BodyValueToString(body, "price")
	description_ := tools.BodyValueToString(body, "description")
	validated_ := tools.BodyValueToString(body, "validated")
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

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, nbPlaces_, price_, description_, validated_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, nbPlaces_, price_, description_, validated_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
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