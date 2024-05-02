package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Housing(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		HousingPost(w, r, db)
	case "GET":
		HousingGet(w, r, db)
	case "PUT":
		HousingPut(w, r, db)
	case "DELETE":
		HousingDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func HousingPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `surface`, `price`, `validated`, `street_nb`, `city`, `zip_code`, `street`, `description`, `house_type`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    surface_ := tools.BodyValueToString(body, "surface")
	price_ := tools.BodyValueToString(body, "price")
	validated_ := tools.BodyValueToString(body, "validated")
	street_nb_ := tools.BodyValueToString(body, "street_nb")
	city_ := tools.BodyValueToString(body, "city")
	zip_code_ := tools.BodyValueToString(body, "zip_code")
	street_ := tools.BodyValueToString(body, "street")
	description_ := tools.BodyValueToString(body, "description")
	house_type_ := tools.BodyValueToString(body, "house_type")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(surface_, price_, validated_, street_nb_, city_, zip_code_, street_, description_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, surface_, price_, validated_, street_nb_, city_, zip_code_, street_, description_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, surface_, price_, validated_, street_nb_, city_, zip_code_, street_, description_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

    if !tools.ValueIsEmpty(house_type_) {
		if !tools.ElementExists(db, "HOUSE_TYPE", "uuid", house_type_) {
			tools.JsonResponse(w, 400, `{"error": "This house_type does not exist"}`) 
			return
		}
	}
	
	
	

	

	

	uuid_ := tools.GenerateUUID()

	// Inserting the Housing in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `HOUSING` (`uuid`, `surface`, `price`, `validated`, `street_nb`, `city`, `zip_code`, `street`, `description`, `house_type`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", uuid_, surface_, price_, validated_, street_nb_, city_, zip_code_, street_, description_, house_type_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Housing created"`

	// Adding the return fields of the query
	fields, err := HousingGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func HousingGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`, `surface`, `price`, `validated`, `street_nb`, `city`, `zip_code`, `street`, `description`, `house_type`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `surface`, `price`, `validated`, `street_nb`, `city`, `zip_code`, `street`, `description`, `house_type` FROM `HOUSING`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `HOUSING`"
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
	jsonResponse, err := HousingGetAllAssociation(result, true)
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

func HousingPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `surface`, `price`, `validated`, `street_nb`, `city`, `zip_code`, `street`, `description`, `house_type`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	
    surface_ := tools.BodyValueToString(body, "surface")
	price_ := tools.BodyValueToString(body, "price")
	validated_ := tools.BodyValueToString(body, "validated")
	street_nb_ := tools.BodyValueToString(body, "street_nb")
	city_ := tools.BodyValueToString(body, "city")
	zip_code_ := tools.BodyValueToString(body, "zip_code")
	street_ := tools.BodyValueToString(body, "street")
	description_ := tools.BodyValueToString(body, "description")
	house_type_ := tools.BodyValueToString(body, "house_type")
	

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
	if tools.ValueTooShort(4, surface_, price_, validated_, street_nb_, city_, zip_code_, street_, description_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, surface_, price_, validated_, street_nb_, city_, zip_code_, street_, description_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

    if !tools.ValueIsEmpty(house_type_) {
		if !tools.ElementExists(db, "HOUSE_TYPE", "uuid", house_type_) {
			tools.JsonResponse(w, 400, `{"error": "This house_type does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExists(db, "HOUSING", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Housing does not exist"}`) 
		return
	}
	

	

    

	request := "UPDATE `HOUSING` SET "
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
	jsonResponse := `{"message": "Housing updated"`
	
	// Adding the return fields of the query
	fields, err := HousingGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func HousingDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "HOUSING", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Housing does not exist"}`) 
		return
	}
	

	// Deleting the Housing in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `HOUSING` WHERE uuid = ?", uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Housing deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func HousingGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `surface`, `price`, `validated`, `street_nb`, `city`, `zip_code`, `street`, `description`, `house_type` FROM `HOUSING` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return HousingGetAllAssociation(result, arrayOutput)
}

func HousingGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, surface_, price_, validated_, street_nb_, city_, zip_code_, street_, description_, house_type_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &surface_, &price_, &validated_, &street_nb_, &city_, &zip_code_, &street_, &description_, &house_type_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "surface": "` + surface_ + `", "price": "` + price_ + `", "validated": "` + validated_ + `", "street_nb": "` + street_nb_ + `", "city": "` + city_ + `", "zip_code": "` + zip_code_ + `", "street": "` + street_ + `", "description": "` + description_ + `", "house_type": "` + house_type_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &surface_, &price_, &validated_, &street_nb_, &city_, &zip_code_, &street_, &description_, &house_type_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "surface": "` + surface_ + `", "price": "` + price_ + `", "validated": "` + validated_ + `", "street_nb": "` + street_nb_ + `", "city": "` + city_ + `", "zip_code": "` + zip_code_ + `", "street": "` + street_ + `", "description": "` + description_ + `", "house_type": "` + house_type_ + `"`, nil
	}
}