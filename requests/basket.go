package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Basket(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		BasketPost(w, r, db)
	case "GET":
		BasketGet(w, r, db)
	case "DELETE":
		BasketDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func BasketPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	body := tools.ReadBody(r)

	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.AtLeastOneValueInBody(body, `account`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := tools.BodyValueToString(body, "account")

	if tools.GetUUID(r, db) != account_ {
		tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		return
	}

	// Checking if the value is empty
	if tools.ValueIsEmpty(account_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the account exists
	if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
		tools.JsonResponse(w, 404, `{"message": "Account not found"}`)
		return
	}

	// Checking if the account has a basket not paid
	rows, err := db.Query("SELECT * FROM BASKET WHERE account = ? AND paid = false", account_)
	if err != nil {
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer rows.Close()

	if rows.Next() {
		tools.JsonResponse(w, 400, `{"message": "A basket is already opened"}`)
		return
	}

	uuid := tools.GenerateUUID()

	// Inserting the basket in the database
	_, err = db.Exec("INSERT INTO BASKET (uuid, account, paid) VALUES (?, ?, false)", uuid, account_)
	if err != nil {
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	// Sending the response
	jsonResponse := `{"message": "Basket created", `

	// Adding the return fields of the query
	fields, err := BasketGetAll(db, uuid, false)
	if err != nil {
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	jsonResponse += fields + `}`
	tools.JsonResponse(w, 201, jsonResponse)

}

func BasketGet(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	query := tools.ReadQuery(r)

	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `account`, `all`, `paid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	account_ := query["account"]

	if tools.GetUUID(r, db) != account_ && !tools.IsAdmin(r, db) {
		tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		return
	}

	// Checking if the account exists
	if !tools.ElementExists(db, "ACCOUNT", "uuid", account_) {
		tools.JsonResponse(w, 404, `{"message": "Account not found"}`)
		return
	}

	mainRequest := `SELECT
			B.account AS ACCOUNT,
			B.paid AS PAID
		FROM BASKET AS B`
	mainParams := []interface{}{}
	housingRequest := `SELECT
		case when H.start_time is null then 'null' else H.start_time end AS HOUSING_startTime,
		case when H.end_time is null then 'null' else H.end_time end AS HOUSING_endTime,
		case when HO.uuid is null then 'null' else HO.uuid end AS HOUSING_uuid,
		case when HO.surface is null then 'null' else HO.surface end AS HOUSING_surface,
		case when HO.price is null then 'null' else HO.price end AS HOUSING_price,
		case when HO.validated is null then 'null' else HO.validated end AS HOUSING_validated,
		case when HO.street_nb is null then 'null' else HO.street_nb end AS HOUSING_streetNb,
		case when HO.city is null then 'null' else HO.city end AS HOUSING_city,
		case when HO.zip_code is null then 'null' else HO.zip_code end AS HOUSING_zipCode,
		case when HO.street is null then 'null' else HO.street end AS HOUSING_street,
		case when HO.description is null then 'null' else HO.description end AS HOUSING_description,
		case when HO.imgPath is null then 'null' else HO.imgPath end AS HOUSING_imgPath,
		case when HO.house_type is null then 'null' else HO.house_type end AS HOUSING_houseType,
		case when HO.account is null then 'null' else HO.account end AS HOUSING_account
	FROM BASKET AS B
		LEFT JOIN BASKET_HOUSING AS H ON H.basket = B.uuid
		LEFT JOIN housing AS HO ON HO.uuid = H.housing`
	housingParams := []interface{}{}
	bedroomRequest := `SELECT
		case when BE.start_time is null then 'null' else BE.start_time end AS BEDROOM_startTime,
		case when BE.end_time is null then 'null' else BE.end_time end AS BEDROOM_endTime,
		case when BED.uuid is null then 'null' else BED.uuid end AS BEDROOM_uuid,
		case when BED.nbPlaces is null then 'null' else BED.nbPlaces end AS BEDROOM_nbPlaces,
		case when BED.price is null then 'null' else BED.price end AS BEDROOM_price,
		case when BED.description is null then 'null' else BED.description end AS BEDROOM_description,
		case when BED.validated is null then 'null' else BED.validated end AS BEDROOM_validated,
		case when BED.imgPath is null then 'null' else BED.imgPath end AS BEDROOM_imgPath,
		case when BED.housing is null then 'null' else BED.housing end AS BEDROOM_housing
	FROM BASKET AS B	
		LEFT JOIN BASKET_BEDROOM AS BE ON BE.basket = B.uuid
		LEFT JOIN bed_room AS BED ON BED.uuid = BE.bedroom`
	bedroomParams := []interface{}{}
	serviceRequest := `SELECT	
		case when S.start_time is null then 'null' else S.start_time end AS SERVICE_startTime,
		case when S.start_time is null then 'null' else ADDTIME(S.start_time, SE.duration) end AS SERVICE_endTime,
		case when SE.uuid is null then 'null' else SE.uuid end AS SERVICE_uuid,
		case when SE.price is null then 'null' else SE.price end AS SERVICE_price,
		case when SE.description is null then 'null' else SE.description end AS SERVICE_description,
		case when SE.imgPath is null then 'null' else SE.imgPath end AS SERVICE_imgPath,
		case when SE.duration is null then 'null' else SE.duration end AS SERVICE_duration,
		case when SE.account is null then 'null' else SE.account end AS SERVICE_account,
		case when SE.service_type is null then 'null' else SE.service_type end AS SERVICE_serviceType
	FROM BASKET AS B	
		LEFT JOIN basket_service AS S ON S.basket = B.uuid
		LEFT JOIN services AS SE ON SE.uuid = S.service`
	serviceParams := []interface{}{}
	equipmentRequest := `SELECT	
		case when BEQ.number is null then 'null' else BEQ.number end AS EQUIPMENT_number,
		case when E.uuid is null then 'null' else E.uuid end AS EQUIPMENT_uuid,
		case when E.name is null then 'null' else E.name end AS EQUIPMENT_name,
		case when E.description is null then 'null' else E.description end AS EQUIPMENT_description,
		case when E.price is null then 'null' else E.price end AS EQUIPMENT_price,
		case when E.number is null then 'null' else E.number end AS EQUIPMENT_numberTotal,
		case when E.imgPath is null then 'null' else E.imgPath end AS EQUIPMENT_imgPath,
		case when E.equipment_type is null then 'null' else E.equipment_type end AS EQUIPMENT_equipmentType,
		case when E.housing is null then 'null' else E.housing end AS EQUIPMENT_housing
	FROM BASKET AS B
		
		LEFT JOIN basket_equipment AS BEQ ON BEQ.basket = B.uuid
		LEFT JOIN equipment AS E ON E.uuid = BEQ.equipment`
	equipmentParams := []interface{}{}
	countRequest := "SELECT COUNT(*) FROM `BASKET` AS B"
	var countParams []interface{}

	if query["all"] != "true" {
		mainRequest += " WHERE "
		housingRequest += " WHERE "
		bedroomRequest += " WHERE "
		serviceRequest += " WHERE "
		equipmentRequest += " WHERE "
		countRequest += " WHERE "
		strictSearch := query["strictSearch"] == "true"

		for key, value := range query {
			// tools.AppendCondition(&request, &params, "B." + key, value, strictSearch)
			if key != "limit" && key != "offset" && key != "strictSearch" {
				tools.AppendCondition(&mainRequest, &mainParams, "B." + key, value, strictSearch)
				tools.AppendCondition(&housingRequest, &housingParams, "B." + key, value, strictSearch)
				tools.AppendCondition(&bedroomRequest, &bedroomParams, "B." + key, value, strictSearch)
				tools.AppendCondition(&serviceRequest, &serviceParams, "B." + key, value, strictSearch)
				tools.AppendCondition(&equipmentRequest, &equipmentParams, "B." + key, value, strictSearch)
			}
			tools.AppendCondition(&countRequest, &countParams, key, value, strictSearch)
		}

		// Removing the last "AND"
		mainRequest = mainRequest[:len(mainRequest)-3]
		housingRequest = housingRequest[:len(housingRequest)-3]
		bedroomRequest = bedroomRequest[:len(bedroomRequest)-3]
		serviceRequest = serviceRequest[:len(serviceRequest)-3]
		equipmentRequest = equipmentRequest[:len(equipmentRequest)-3]
		countRequest = countRequest[:len(countRequest)-3]
	}

	if query["limit"] != "" {
		mainRequest += " LIMIT " + query["limit"]
		housingRequest += " LIMIT " + query["limit"]
		bedroomRequest += " LIMIT " + query["limit"]
		serviceRequest += " LIMIT " + query["limit"]
		equipmentRequest += " LIMIT " + query["limit"]

		if query["offset"] != "" {
			mainRequest += " OFFSET " + query["offset"]
			housingRequest += " OFFSET " + query["offset"]
			bedroomRequest += " OFFSET " + query["offset"]
			serviceRequest += " OFFSET " + query["offset"]
			equipmentRequest += " OFFSET " + query["offset"]
		}
	}

	countResult, err := tools.ExecuteQuery(db, countRequest, countParams...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	var count string
	for countResult.Next() {
		err := countResult.Scan(&count)
		if err != nil {
			tools.ErrorLog(err.Error())
			tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
			return
		}
	}
	countResult.Close()

	jsonResponse := `{"count": ` + count + `, "baskets": [`

	mainResult, err := tools.ExecuteQuery(db, mainRequest, mainParams...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer mainResult.Close()

	for mainResult.Next() {
		var (
			account_ string
			paid_ string
		)

		err := mainResult.Scan(&account_, &paid_)
		if err != nil {
			tools.ErrorLog(err.Error())
			tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
			return
		}

		jsonResponse += `{"account": "` + account_ + `", "paid": "` + paid_ + `", "HOUSINGS": [`

	}

	housingResult, err := tools.ExecuteQuery(db, housingRequest, housingParams...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer housingResult.Close()

	for housingResult.Next() {
		var (
			HOUSING_startTime string
			HOUSING_endTime string
			HOUSING_uuid string
			HOUSING_surface string
			HOUSING_price string
			HOUSING_validated string
			HOUSING_streetNb string
			HOUSING_city string
			HOUSING_zipCode string
			HOUSING_street string
			HOUSING_description string
			HOUSING_imgPath string
			HOUSING_houseType string
			HOUSING_account string
		)

		err := housingResult.Scan(&HOUSING_startTime, &HOUSING_endTime, &HOUSING_uuid, &HOUSING_surface, &HOUSING_price, &HOUSING_validated, &HOUSING_streetNb, &HOUSING_city, &HOUSING_zipCode, &HOUSING_street, &HOUSING_description, &HOUSING_imgPath, &HOUSING_houseType, &HOUSING_account)
		if err != nil {
			tools.ErrorLog(err.Error())
			tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
			return
		}

		if HOUSING_startTime != "NULL" {
			jsonResponse += `{"startTime": "` + HOUSING_startTime + `", "endTime": "` + HOUSING_endTime + `", "uuid": "` + HOUSING_uuid + `", "surface": "` + HOUSING_surface + `", "price": "` + HOUSING_price + `", "validated": "` + HOUSING_validated + `", "streetNb": "` + HOUSING_streetNb + `", "city": "` + HOUSING_city + `", "zipCode": "` + HOUSING_zipCode + `", "street": "` + HOUSING_street + `", "description": "` + HOUSING_description + `", "imgPath": "` + HOUSING_imgPath + `", "houseType": "` + HOUSING_houseType + `", "account": "` + HOUSING_account + `"},`
		}

	}

	if jsonResponse[len(jsonResponse)-1] == ',' {
		jsonResponse = jsonResponse[:len(jsonResponse)-1] // Removing the last ","
	}	
	jsonResponse += `], "BEDROOMS": [`

	bedroomResult, err := tools.ExecuteQuery(db, bedroomRequest, bedroomParams...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer bedroomResult.Close()

	for bedroomResult.Next() {
		var (
			BEDROOM_startTime string
			BEDROOM_endTime string
			BEDROOM_uuid string
			BEDROOM_nbPlaces string
			BEDROOM_price string
			BEDROOM_description string
			BEDROOM_validated string
			BEDROOM_imgPath string
			BEDROOM_housing string
		)

		err := bedroomResult.Scan(&BEDROOM_startTime, &BEDROOM_endTime, &BEDROOM_uuid, &BEDROOM_nbPlaces, &BEDROOM_price, &BEDROOM_description, &BEDROOM_validated, &BEDROOM_imgPath, &BEDROOM_housing)
		if err != nil {
			tools.ErrorLog(err.Error())
			tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
			return
		}

		if BEDROOM_startTime != "NULL" {
			jsonResponse += `{"startTime": "` + BEDROOM_startTime + `", "endTime": "` + BEDROOM_endTime + `", "uuid": "` + BEDROOM_uuid + `", "nbPlaces": "` + BEDROOM_nbPlaces + `", "price": "` + BEDROOM_price + `", "description": "` + BEDROOM_description + `", "validated": "` + BEDROOM_validated + `", "imgPath": "` + BEDROOM_imgPath + `", "housing": "` + BEDROOM_housing + `"},`
		}

	}

	if jsonResponse[len(jsonResponse)-1] == ',' {
		jsonResponse = jsonResponse[:len(jsonResponse)-1] // Removing the last ","
	}
	jsonResponse += `], "SERVICES": [`

	serviceResult, err := tools.ExecuteQuery(db, serviceRequest, serviceParams...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	for serviceResult.Next() {
		var (
			SERVICE_startTime string
			SERVICE_endTime string
			SERVICE_uuid string
			SERVICE_price string
			SERVICE_description string
			SERVICE_imgPath string
			SERVICE_duration string
			SERVICE_account string
			SERVICE_serviceType string
		)

		err := serviceResult.Scan(&SERVICE_startTime, &SERVICE_endTime, &SERVICE_uuid, &SERVICE_price, &SERVICE_description, &SERVICE_imgPath, &SERVICE_duration, &SERVICE_account, &SERVICE_serviceType)
		if err != nil {
			tools.ErrorLog(err.Error())
			tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
			return
		}

		if SERVICE_startTime != "NULL" {
		
			jsonResponse += `{"startTime": "` + SERVICE_startTime + `", "endTime": "` + SERVICE_endTime + `", "uuid": "` + SERVICE_uuid + `", "price": "` + SERVICE_price + `", "description": "` + SERVICE_description + `", "imgPath": "` + SERVICE_imgPath + `", "duration": "` + SERVICE_duration + `", "account": "` + SERVICE_account + `", "serviceType": "` + SERVICE_serviceType + `"},`
		
		}

	}

	if jsonResponse[len(jsonResponse)-1] == ',' {
		jsonResponse = jsonResponse[:len(jsonResponse)-1] // Removing the last ","
	}
	jsonResponse += `], "EQUIPMENTS": [`

	equipmentResult, err := tools.ExecuteQuery(db, equipmentRequest, equipmentParams...)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	for equipmentResult.Next() {
		var (
			EQUIPMENT_number string
			EQUIPMENT_uuid string
			EQUIPMENT_name string
			EQUIPMENT_description string
			EQUIPMENT_price string
			EQUIPMENT_numberTotal string
			EQUIPMENT_imgPath string
			EQUIPMENT_equipmentType string
			EQUIPMENT_housing string
		)

		err := equipmentResult.Scan(&EQUIPMENT_number, &EQUIPMENT_uuid, &EQUIPMENT_name, &EQUIPMENT_description, &EQUIPMENT_price, &EQUIPMENT_numberTotal, &EQUIPMENT_imgPath, &EQUIPMENT_equipmentType, &EQUIPMENT_housing)
		if err != nil {
			tools.ErrorLog(err.Error())
			tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
			return
		}

		if EQUIPMENT_number != "NULL" {

			jsonResponse += `{"number": "` + EQUIPMENT_number + `", "uuid": "` + EQUIPMENT_uuid + `", "name": "` + EQUIPMENT_name + `", "description": "` + EQUIPMENT_description + `", "price": "` + EQUIPMENT_price + `", "numberTotal": "` + EQUIPMENT_numberTotal + `", "imgPath": "` + EQUIPMENT_imgPath + `", "equipmentType": "` + EQUIPMENT_equipmentType + `", "housing": "` + EQUIPMENT_housing + `"},`

		}

	}

	if jsonResponse[len(jsonResponse)-1] == ',' {
		jsonResponse = jsonResponse[:len(jsonResponse)-1] // Removing the last ","
	}
	jsonResponse += `]}]}`

	tools.JsonResponse(w, 200, jsonResponse)

}

func BasketDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	query := tools.ReadQuery(r)

	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.AtLeastOneValueInQuery(query, `uuid`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	uuid := query["uuid"]

	// Checking if the basket exists
	if !tools.ElementExists(db, "BASKET", "uuid", uuid) {
		tools.JsonResponse(w, 404, `{"message": "Basket not found"}`)
		return
	}

	// Checking if the account has the right to delete the basket
	if tools.GetUUID(r, db) != tools.GetElement(db, "BASKET", "account", "uuid", uuid) && !tools.IsAdmin(r, db) {
		tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		return
	}

	// Deleting the linked tables using a transaction
	tx, err := db.Begin()
	if err != nil {
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	_, err = tx.Exec("DELETE FROM BASKET_HOUSING WHERE basket = ?", uuid)
	if err != nil {
		tx.Rollback()
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	_, err = tx.Exec("DELETE FROM BASKET_BEDROOM WHERE basket = ?", uuid)
	if err != nil {
		tx.Rollback()
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	_, err = tx.Exec("DELETE FROM BASKET_SERVICE WHERE basket = ?", uuid)
	if err != nil {
		tx.Rollback()
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	_, err = tx.Exec("DELETE FROM BASKET_EQUIPMENT WHERE basket = ?", uuid)
	if err != nil {
		tx.Rollback()
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	_, err = tx.Exec("DELETE FROM BASKET WHERE uuid = ?", uuid)
	if err != nil {
		tx.Rollback()
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	tx.Commit()

	// Sending the response
	tools.JsonResponse(w, 200, `{"message": "Basket deleted"}`)

}

func BasketGetAll(db *sql.DB, uuid_ string, arrayOutput bool) (string, error) {
	result, err := tools.ExecuteQuery(db, `SELECT
	B.account AS ACCOUNT,
	B.paid AS PAID,
	
	case when H.start_time is null then 'null' else H.start_time end AS HOUSING_startTime,
	case when H.end_time is null then 'null' else H.end_time end AS HOUSING_endTime,
	case when HO.uuid is null then 'null' else HO.uuid end AS HOUSING_uuid,
	case when HO.surface is null then 'null' else HO.surface end AS HOUSING_surface,
	case when HO.price is null then 'null' else HO.price end AS HOUSING_price,
	case when HO.validated is null then 'null' else HO.validated end AS HOUSING_validated,
	case when HO.street_nb is null then 'null' else HO.street_nb end AS HOUSING_streetNb,
	case when HO.city is null then 'null' else HO.city end AS HOUSING_city,
	case when HO.zip_code is null then 'null' else HO.zip_code end AS HOUSING_zipCode,
	case when HO.street is null then 'null' else HO.street end AS HOUSING_street,
	case when HO.description is null then 'null' else HO.description end AS HOUSING_description,
	case when HO.imgPath is null then 'null' else HO.imgPath end AS HOUSING_imgPath,
	case when HO.house_type is null then 'null' else HO.house_type end AS HOUSING_houseType,
	case when HO.account is null then 'null' else HO.account end AS HOUSING_account,
	
	case when BE.start_time is null then 'null' else BE.start_time end AS BEDROOM_startTime,
	case when BE.end_time is null then 'null' else BE.end_time end AS BEDROOM_endTime,
	case when BED.uuid is null then 'null' else BED.uuid end AS BEDROOM_uuid,
	case when BED.nbPlaces is null then 'null' else BED.nbPlaces end AS BEDROOM_nbPlaces,
	case when BED.price is null then 'null' else BED.price end AS BEDROOM_price,
	case when BED.description is null then 'null' else BED.description end AS BEDROOM_description,
	case when BED.validated is null then 'null' else BED.validated end AS BEDROOM_validated,
	case when BED.imgPath is null then 'null' else BED.imgPath end AS BEDROOM_imgPath,
	case when BED.housing is null then 'null' else BED.housing end AS BEDROOM_housing,
	
	case when S.start_time is null then 'null' else S.start_time end AS SERVICE_startTime,
	case when S.start_time is null then 'null' else ADDTIME(S.start_time, SE.duration) end AS SERVICE_endTime,
	case when SE.uuid is null then 'null' else SE.uuid end AS SERVICE_uuid,
	case when SE.price is null then 'null' else SE.price end AS SERVICE_price,
	case when SE.description is null then 'null' else SE.description end AS SERVICE_description,
	case when SE.imgPath is null then 'null' else SE.imgPath end AS SERVICE_imgPath,
	case when SE.duration is null then 'null' else SE.duration end AS SERVICE_duration,
	case when SE.account is null then 'null' else SE.account end AS SERVICE_account,
	case when SE.service_type is null then 'null' else SE.service_type end AS SERVICE_serviceType,
	
	case when BEQ.number is null then 'null' else BEQ.number end AS EQUIPMENT_number,
	case when E.uuid is null then 'null' else E.uuid end AS EQUIPMENT_uuid,
	case when E.name is null then 'null' else E.name end AS EQUIPMENT_name,
	case when E.description is null then 'null' else E.description end AS EQUIPMENT_description,
	case when E.price is null then 'null' else E.price end AS EQUIPMENT_price,
	case when E.number is null then 'null' else E.number end AS EQUIPMENT_numberTotal,
	case when E.imgPath is null then 'null' else E.imgPath end AS EQUIPMENT_imgPath,
	case when E.equipment_type is null then 'null' else E.equipment_type end AS EQUIPMENT_equipmentType,
	case when E.housing is null then 'null' else E.housing end AS EQUIPMENT_housing
FROM BASKET AS B
	LEFT JOIN BASKET_HOUSING AS H ON H.basket = B.uuid
	LEFT JOIN housing AS HO ON HO.uuid = H.housing
	
	LEFT JOIN BASKET_BEDROOM AS BE ON BE.basket = B.uuid
	LEFT JOIN bed_room AS BED ON BED.uuid = BE.bedroom
	
	LEFT JOIN basket_service AS S ON S.basket = B.uuid
	LEFT JOIN services AS SE ON SE.uuid = S.service
	
	LEFT JOIN basket_equipment AS BEQ ON BEQ.basket = B.uuid
	LEFT JOIN equipment AS E ON E.uuid = BEQ.equipment
WHERE
	B.uuid = ?`, uuid_)
	if err != nil {
		return "", err
	}
	defer result.Close()

	return BasketGetAllAssociation(result, arrayOutput)
}

func BasketGetAllAssociation(result *sql.Rows, arrayOutput bool) (string, error) {
	var (
		account_ string
		paid_ string
		HOUSING_startTime string
		HOUSING_endTime string
		HOUSING_uuid string
		HOUSING_surface string
		HOUSING_price string
		HOUSING_validated string
		HOUSING_streetNb string
		HOUSING_city string
		HOUSING_zipCode string
		HOUSING_street string
		HOUSING_description string
		HOUSING_imgPath string
		HOUSING_houseType string
		HOUSING_account string
		BEDROOM_startTime string
		BEDROOM_endTime string
		BEDROOM_uuid string
		BEDROOM_nbPlaces string
		BEDROOM_price string
		BEDROOM_description string
		BEDROOM_validated string
		BEDROOM_imgPath string
		BEDROOM_housing string
		SERVICE_startTime string
		SERVICE_endTime string
		SERVICE_uuid string
		SERVICE_price string
		SERVICE_description string
		SERVICE_imgPath string
		SERVICE_duration string
		SERVICE_account string
		SERVICE_serviceType string
		EQUIPMENT_number string
		EQUIPMENT_uuid string
		EQUIPMENT_name string
		EQUIPMENT_description string
		EQUIPMENT_price string
		EQUIPMENT_numberTotal string
		EQUIPMENT_imgPath string
		EQUIPMENT_equipmentType string
		EQUIPMENT_housing string
	)

	switch arrayOutput {
	case true:
		var jsonResponse string
		jsonResponse += `[`
		for result.Next() {
			err := result.Scan(&account_, &paid_, &HOUSING_startTime, &HOUSING_endTime, &HOUSING_uuid, &HOUSING_surface, &HOUSING_price, &HOUSING_validated, &HOUSING_streetNb, &HOUSING_city, &HOUSING_zipCode, &HOUSING_street, &HOUSING_description, &HOUSING_imgPath, &HOUSING_houseType, &HOUSING_account, &BEDROOM_startTime, &BEDROOM_endTime, &BEDROOM_uuid, &BEDROOM_nbPlaces, &BEDROOM_price, &BEDROOM_description, &BEDROOM_validated, &BEDROOM_imgPath, &BEDROOM_housing, &SERVICE_startTime, &SERVICE_endTime, &SERVICE_uuid, &SERVICE_price, &SERVICE_description, &SERVICE_imgPath, &SERVICE_duration, &SERVICE_account, &SERVICE_serviceType, &EQUIPMENT_number, &EQUIPMENT_uuid, &EQUIPMENT_name, &EQUIPMENT_description, &EQUIPMENT_price, &EQUIPMENT_numberTotal, &EQUIPMENT_imgPath, &EQUIPMENT_equipmentType, &EQUIPMENT_housing)
			if err != nil {
				return "", err
			}
			jsonResponse += `{"account": "` + account_ + `", "paid": "` + paid_ + `", "HOUSING": {"startTime": "` + HOUSING_startTime + `", "endTime": "` + HOUSING_endTime + `", "uuid": "` + HOUSING_uuid + `", "surface": "` + HOUSING_surface + `", "price": "` + HOUSING_price + `", "validated": "` + HOUSING_validated + `", "streetNb": "` + HOUSING_streetNb + `", "city": "` + HOUSING_city + `", "zipCode": "` + HOUSING_zipCode + `", "street": "` + HOUSING_street + `", "description": "` + HOUSING_description + `", "imgPath": "` + HOUSING_imgPath + `", "houseType": "` + HOUSING_houseType + `", "account": "` + HOUSING_account + `"}, "BEDROOM": {"startTime": "` + BEDROOM_startTime + `", "endTime": "` + BEDROOM_endTime + `", "uuid": "` + BEDROOM_uuid + `", "nbPlaces": "` + BEDROOM_nbPlaces + `", "price": "` + BEDROOM_price + `", "description": "` + BEDROOM_description + `", "validated": "` + BEDROOM_validated + `", "imgPath": "` + BEDROOM_imgPath + `", "housing": "` + BEDROOM_housing + `"}, "SERVICE": {"startTime": "` + SERVICE_startTime + `", "endTime": "` + SERVICE_endTime + `", "uuid": "` + SERVICE_uuid + `", "price": "` + SERVICE_price + `", "description": "` + SERVICE_description + `", "imgPath": "` + SERVICE_imgPath + `", "duration": "` + SERVICE_duration + `", "account": "` + SERVICE_account + `", "serviceType": "` + SERVICE_serviceType + `"}, "EQUIPMENT": {"number": "` + EQUIPMENT_number + `", "uuid": "` + EQUIPMENT_uuid + `", "name": "` + EQUIPMENT_name + `", "description": "` + EQUIPMENT_description + `", "price": "` + EQUIPMENT_price + `", "numberTotal": "` + EQUIPMENT_numberTotal + `", "imgPath": "` + EQUIPMENT_imgPath + `", "equipmentType": "` + EQUIPMENT_equipmentType + `", "housing": "` + EQUIPMENT_housing + `"}},`
		}
		if len(jsonResponse) > 1 {
			jsonResponse = jsonResponse[:len(jsonResponse)-1]
		}
		jsonResponse += `]`
		return jsonResponse, nil
	default:
		for result.Next() {
			err := result.Scan(&account_, &paid_, &HOUSING_startTime, &HOUSING_endTime, &HOUSING_uuid, &HOUSING_surface, &HOUSING_price, &HOUSING_validated, &HOUSING_streetNb, &HOUSING_city, &HOUSING_zipCode, &HOUSING_street, &HOUSING_description, &HOUSING_imgPath, &HOUSING_houseType, &HOUSING_account, &BEDROOM_startTime, &BEDROOM_endTime, &BEDROOM_uuid, &BEDROOM_nbPlaces, &BEDROOM_price, &BEDROOM_description, &BEDROOM_validated, &BEDROOM_imgPath, &BEDROOM_housing, &SERVICE_startTime, &SERVICE_endTime, &SERVICE_uuid, &SERVICE_price, &SERVICE_description, &SERVICE_imgPath, &SERVICE_duration, &SERVICE_account, &SERVICE_serviceType, &EQUIPMENT_number, &EQUIPMENT_uuid, &EQUIPMENT_name, &EQUIPMENT_description, &EQUIPMENT_price, &EQUIPMENT_numberTotal, &EQUIPMENT_imgPath, &EQUIPMENT_equipmentType, &EQUIPMENT_housing)
			if err != nil {
				return "", err
			}
		}
		return `"account": "` + account_ + `", "paid": "` + paid_ + `", "HOUSING": {"startTime": "` + HOUSING_startTime + `", "endTime": "` + HOUSING_endTime + `", "uuid": "` + HOUSING_uuid + `", "surface": "` + HOUSING_surface + `", "price": "` + HOUSING_price + `", "validated": "` + HOUSING_validated + `", "streetNb": "` + HOUSING_streetNb + `", "city": "` + HOUSING_city + `", "zipCode": "` + HOUSING_zipCode + `", "street": "` + HOUSING_street + `", "description": "` + HOUSING_description + `", "imgPath": "` + HOUSING_imgPath + `", "houseType": "` + HOUSING_houseType + `", "account": "` + HOUSING_account + `"}, "BEDROOM": {"startTime": "` + BEDROOM_startTime + `", "endTime": "` + BEDROOM_endTime + `", "uuid": "` + BEDROOM_uuid + `", "nbPlaces": "` + BEDROOM_nbPlaces + `", "price": "` + BEDROOM_price + `", "description": "` + BEDROOM_description + `", "validated": "` + BEDROOM_validated + `", "imgPath": "` + BEDROOM_imgPath + `", "housing": "` + BEDROOM_housing + `"}, "SERVICE": {"startTime": "` + SERVICE_startTime + `", "endTime": "` + SERVICE_endTime + `", "uuid": "` + SERVICE_uuid + `", "price": "` + SERVICE_price + `", "description": "` + SERVICE_description + `", "imgPath": "` + SERVICE_imgPath + `", "duration": "` + SERVICE_duration + `", "account": "` + SERVICE_account + `", "serviceType": "` + SERVICE_serviceType + `"}, "EQUIPMENT": {"number": "` + EQUIPMENT_number + `", "uuid": "` + EQUIPMENT_uuid + `", "name": "` + EQUIPMENT_name + `", "description": "` + EQUIPMENT_description + `", "price": "` + EQUIPMENT_price + `", "numberTotal": "` + EQUIPMENT_numberTotal + `", "imgPath": "` + EQUIPMENT_imgPath + `", "equipmentType": "` + EQUIPMENT_equipmentType + `", "housing": "` + EQUIPMENT_housing + `"}`, nil
	}

}