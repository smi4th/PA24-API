package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func EquipmentType(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		if tools.IsAdmin(r, db) {
			EquipmentTypePost(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	case "GET":
		EquipmentTypeGet(w, r, db)
	case "PUT":
		if tools.IsAdmin(r, db) {
			EquipmentTypePut(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	case "DELETE":
		if tools.IsAdmin(r, db) {
			EquipmentTypeDelete(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func EquipmentTypePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `name`, `imgPath`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    name_ := tools.BodyValueToString(body, "name")
	imgPath_ := tools.BodyValueToString(body, "imgPath")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(name_, imgPath_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, name_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, name_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

    
	
	

	

	if tools.ElementExists(db, "EQUIPMENT_TYPE", "name", name_) {
		tools.JsonResponse(w, 400, `{"error": "This name already exists"}`) 
		return
	}
	

	uuid_ := tools.GenerateUUID()

	// Inserting the EquipmentType in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `EQUIPMENT_TYPE` (`uuid`, `name`, `imgPath`) VALUES (?, ?, ?)", uuid_, name_, imgPath_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "EquipmentType created"`

	// Adding the return fields of the query
	fields, err := EquipmentTypeGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func EquipmentTypeGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`, `name`, "all", "imgPath") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `name`, `imgPath` FROM `EQUIPMENT_TYPE`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `EQUIPMENT_TYPE`"
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
	jsonResponse, err := EquipmentTypeGetAllAssociation(result, true)
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

func EquipmentTypePut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `name`, `imgPath`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	if !tools.AtLeastOneValueInBody(body, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Cannot update all fields"}`)
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

    

	if !tools.ElementExists(db, "EQUIPMENT_TYPE", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This EquipmentType does not exist"}`) 
		return
	}
	if tools.ElementExists(db, "EQUIPMENT_TYPE", "name", name_) {
		tools.JsonResponse(w, 400, `{"error": "This name already exists"}`) 
		return
	}
	

	

    

	request := "UPDATE `EQUIPMENT_TYPE` SET "
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
	jsonResponse := `{"message": "EquipmentType updated"`
	
	// Adding the return fields of the query
	fields, err := EquipmentTypeGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func EquipmentTypeDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "EQUIPMENT_TYPE", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This EquipmentType does not exist"}`) 
		return
	}
	

	// Deleting the EquipmentType in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `EQUIPMENT_TYPE` WHERE uuid = ?", uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "EquipmentType deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func EquipmentTypeGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `name`, `imgPath` FROM `EQUIPMENT_TYPE` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return EquipmentTypeGetAllAssociation(result, arrayOutput)
}

func EquipmentTypeGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, name_, imgPath_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &name_, &imgPath_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "name": "` + name_ + `", "imgPath": "` + imgPath_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &name_, &imgPath_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "name": "` + name_ + `", "imgPath": "` + imgPath_ + `"`, nil
	}
}