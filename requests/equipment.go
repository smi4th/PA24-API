package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Equipment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		EquipmentPost(w, r, db)
	case "GET":
		EquipmentGet(w, r, db)
	case "PUT":
		if tools.GetUUID(r, db) == tools.GetElement(db, "HOUSING", "account", "uuid", tools.GetElement(db, "EQUIPMENT", "housing", "uuid", tools.ReadQuery(r)["uuid"])) {
			EquipmentPut(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	case "DELETE":
		if tools.GetUUID(r, db) == tools.GetElement(db, "HOUSING", "account", "uuid", tools.GetElement(db, "EQUIPMENT", "housing", "uuid", tools.ReadQuery(r)["uuid"])) {
			EquipmentDelete(w, r, db)
		} else {
			tools.JsonResponse(w, 403, `{"message": "Forbidden"}`)
		}
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func EquipmentPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `name`, `description`, `price`, `number`, `equipment_type`, `housing`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    name_ := tools.BodyValueToString(body, "name")
	description_ := tools.BodyValueToString(body, "description")
	price_ := tools.BodyValueToString(body, "price")
	number_ := tools.BodyValueToString(body, "number")
	equipment_type_ := tools.BodyValueToString(body, "equipment_type")
	housing_ := tools.BodyValueToString(body, "housing")

	if tools.GetUUID(r, db) != tools.GetElement(db, "HOUSING", "account", "uuid", housing_) {
		tools.JsonResponse(w, 403, `{"error": "Forbidden"}`)
		return
	}

	// Checking if the values are empty
	if tools.ValueIsEmpty(name_, description_, price_, number_, equipment_type_, housing_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, name_, description_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, name_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

	if !tools.ElementExists(db, "EQUIPMENT_TYPE", "uuid", equipment_type_) {
		tools.JsonResponse(w, 400, `{"error": "This equipment_type does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "HOUSING", "uuid", housing_) {
		tools.JsonResponse(w, 400, `{"error": "This housing does not exist"}`)
		return
	}
	

	uuid_ := tools.GenerateUUID()

	// Inserting the Equipment in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `EQUIPMENT` (`uuid`, `name`, `description`, `price`, `number`, `equipment_type`, `housing`) VALUES (?, ?, ?, ?, ?, ?, ?)", uuid_, name_, description_, price_, number_, equipment_type_, housing_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Equipment created"`

	// Adding the return fields of the query
	fields, err := EquipmentGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func EquipmentGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`, `name`, `description`, `price`, `number`, `equipment_type`, `housing`, `all`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `uuid`, `name`, `description`, `price`, `number`, `equipment_type`, `housing` FROM `EQUIPMENT`"
	var params []interface{}
	countRequest := "SELECT COUNT(*) FROM `EQUIPMENT`"
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
	jsonResponse, err := EquipmentGetAllAssociation(result, true)
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

func EquipmentPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `name`, `description`, `price`, `number`, `housing`, `equipment_type`) || tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	
    name_ := tools.BodyValueToString(body, "name")
	description_ := tools.BodyValueToString(body, "description")
	equipment_type_ := tools.BodyValueToString(body, "equipment_type")
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
	if tools.ValueTooShort(4, name_, description_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, name_, description_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

    if !tools.ValueIsEmpty(equipment_type_) {
		if !tools.ElementExists(db, "EQUIPMENT_TYPE", "uuid", equipment_type_) {
			tools.JsonResponse(w, 400, `{"error": "This equipment_type does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(housing_) {
		if !tools.ElementExists(db, "HOUSING", "uuid", housing_) {
			tools.JsonResponse(w, 400, `{"error": "This housing does not exist"}`)
			return
		}
	}
	

	if !tools.ElementExists(db, "EQUIPMENT", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Equipment does not exist"}`) 
		return
	}

    

	request := "UPDATE `EQUIPMENT` SET "
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
	jsonResponse := `{"message": "Equipment updated"`
	
	// Adding the return fields of the query
	fields, err := EquipmentGetAll(db, uuid_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func EquipmentDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid_ := query["uuid"]
	

	if !tools.ElementExists(db, "EQUIPMENT", "uuid", uuid_) {
		tools.JsonResponse(w, 400, `{"error": "This Equipment does not exist"}`) 
		return
	}
	

	// Deleting the Equipment in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `EQUIPMENT` WHERE uuid = ?", uuid_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Equipment deleted", "uuid": "` + uuid_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func EquipmentGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `uuid`, `name`, `description`, `price`, `number`, `equipment_type`, `housing` FROM `EQUIPMENT` WHERE uuid = ?", uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return EquipmentGetAllAssociation(result, arrayOutput)
}

func EquipmentGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var uuid_, name_, description_, price_, number_, equipment_type_, housing_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&uuid_, &name_, &description_, &price_, &number_, &equipment_type_, &housing_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"uuid": "` + uuid_ + `", "name": "` + name_ + `", "description": "` + description_ + `", "price": "` + price_ + `", "number": "` + number_ + `", "equipment_type": "` + equipment_type_ + `", "housing": "` + housing_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&uuid_, &name_, &description_, &price_, &number_, &equipment_type_, &housing_)
			if err != nil {
				return "", err
			}
		}
		return `"uuid": "` + uuid_ + `", "name": "` + name_ + `", "description": "` + description_ + `", "price": "` + price_ + `", "number": "` + number_ + `", "equipment_type": "` + equipment_type_ + `", "housing": "` + housing_ + `"`, nil
	}
}