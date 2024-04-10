package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func HousingEquipment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		HousingEquipmentPost(w, r, db)
	case "GET":
		HousingEquipmentGet(w, r, db)
	case "PUT":
		HousingEquipmentPut(w, r, db)
	case "DELETE":
		HousingEquipmentDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func HousingEquipmentPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `number`, `housing`, `equipment`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    number_ := tools.BodyValueToString(body, "number")
	housing_ := tools.BodyValueToString(body, "housing")
	equipment_ := tools.BodyValueToString(body, "equipment")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(number_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(4, number_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too short"}`)
		return
	}
	if tools.ValueTooLong(32, number_) {
		tools.JsonResponse(w, 400, `{"message": "Fields too long"}`)
		return
	}

    if !tools.ValueIsEmpty(housing_) {
		if !tools.ElementExists(db, "HOUSING", "uuid", housing_) {
			tools.JsonResponse(w, 400, `{"error": "This housing does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(equipment_) {
		if !tools.ElementExists(db, "EQUIPMENT", "uuid", equipment_) {
			tools.JsonResponse(w, 400, `{"error": "This equipment does not exist"}`) 
			return
		}
	}
	
	
	

	

	

	

	// Inserting the HousingEquipment in the database
	result, err := tools.ExecuteQuery(db, "INSERT INTO `HOUSING_EQUIPMENT` (`housing`, `equipment`, `number`, `housing`, `equipment`) VALUES (?, ?, ?)", housing_, equipment_, number_, housing_, equipment_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "HousingEquipment created"`

	// Adding the return fields of the query
	fields, err := HousingEquipmentGetAll(db, housing_, equipment_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse + "," + fields + "}")

}

func HousingEquipmentGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `number`, `housing`, `equipment`, "all") {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	request := "SELECT `number`, `housing`, `equipment` FROM `HOUSING_EQUIPMENT`"
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
	jsonResponse, err := HousingEquipmentGetAllAssociation(result, true)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func HousingEquipmentPut(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the body of the request
	body := tools.ReadBody(r)
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `number`) || tools.ValuesNotInQuery(query, `housing`, `equipment`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	housing_ := query["housing"]
	equipment_ := query["equipment"]
	
    number_ := tools.BodyValueToString(body, "number")
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(housing_, equipment_) {
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
	if tools.ValueTooShort(4, number_) {
		tools.JsonResponse(w, 400, `{"message": "values too short"}`)
		return
	}
	if tools.ValueTooLong(32, number_) {
		tools.JsonResponse(w, 400, `{"message": "values too long"}`)
		return
	}

    if !tools.ValueIsEmpty(housing_) {
		if !tools.ElementExists(db, "HOUSING", "uuid", housing_) {
			tools.JsonResponse(w, 400, `{"error": "This housing does not exist"}`) 
			return
		}
	}
	if !tools.ValueIsEmpty(equipment_) {
		if !tools.ElementExists(db, "EQUIPMENT", "uuid", equipment_) {
			tools.JsonResponse(w, 400, `{"error": "This equipment does not exist"}`) 
			return
		}
	}
	

	if !tools.ElementExists(db, "HOUSING_EQUIPMENT", "housing", housing_) {
		tools.JsonResponse(w, 400, `{"error": "This HousingEquipment does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "HOUSING_EQUIPMENT", "equipment", equipment_) {
		tools.JsonResponse(w, 400, `{"error": "This HousingEquipment does not exist"}`) 
		return
	}
	

	{{passwordCheck}}

    {{emailCheck}}

	request := "UPDATE `HOUSING_EQUIPMENT` SET "
	var params []interface{}
	
	for key, value := range body {
		if !tools.ValueInArray(key, `housing`, `equipment`) {
			if key == "password" {
				value = tools.HashPassword(value.(string))
			}
			tools.AppendUpdate(&request, &params, key, value)
		}
	}

	// Removing the last ","
	request = request[:len(request)-2]

	request += " WHERE housing = ?, equipment = ?"
	params = append(params, housing_, equipment_)

	// Updating the account in the database
	result, err := tools.ExecuteQuery(db, request, params...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "HousingEquipment updated"`
	
	// Adding the return fields of the query
	fields, err := HousingEquipmentGetAll(db, housing_, equipment_, false)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse + "," + fields + "}")

}

func HousingEquipmentDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	
	// Getting the query parameters
	query := tools.ReadQuery(r)
	
	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `housing`, `equipment`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	housing_ := query["housing"]
	equipment_ := query["equipment"]
	

	if !tools.ElementExists(db, "HOUSING_EQUIPMENT", "housing", housing_) {
		tools.JsonResponse(w, 400, `{"error": "This HousingEquipment does not exist"}`) 
		return
	}
	if !tools.ElementExists(db, "HOUSING_EQUIPMENT", "equipment", equipment_) {
		tools.JsonResponse(w, 400, `{"error": "This HousingEquipment does not exist"}`) 
		return
	}
	

	// Deleting the HousingEquipment in the database
	result, err := tools.ExecuteQuery(db, "DELETE FROM `HOUSING_EQUIPMENT` WHERE housing = ?, equipment = ?", housing_, equipment_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "HousingEquipment deleted", "housing": "` + housing_ + `", "equipment": "` + equipment_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 200, jsonResponse)

}

func HousingEquipmentGetAll(db *sql.DB, housing_ string, equipment_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, "SELECT `number`, `housing`, `equipment` FROM `HOUSING_EQUIPMENT` WHERE housing = ?, equipment = ?", housing_, equipment_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return HousingEquipmentGetAllAssociation(result, arrayOutput)
}

func HousingEquipmentGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var number_, housing_, equipment_ string

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&number_, &housing_, &equipment_)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"number": "` + number_ + `", "housing": "` + housing_ + `", "equipment": "` + equipment_ + `"},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&number_, &housing_, &equipment_)
			if err != nil {
				return "", err
			}
		}
		return `"number": "` + number_ + `", "housing": "` + housing_ + `", "equipment": "` + equipment_ + `"`, nil
	}
}