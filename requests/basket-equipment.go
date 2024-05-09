package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func BasketEquipment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		BasketEquipmentPost(w, r, db)
	case "DELETE":
		BasketEquipmentDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func BasketEquipmentPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	body := tools.ReadBody(r)

	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `basket`, `equipment`, `number`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	basket_ := tools.BodyValueToString(body, `basket`)
	equipment_ := tools.BodyValueToString(body, `equipment`)
	number_ := tools.BodyValueToString(body, `number`)

	// Check if basket exists
	if !tools.ElementExists(db, "BASKET", "uuid", basket_) {
		tools.JsonResponse(w, 404, `{"message": "Basket not found"}`)
		return
	}

	if tools.GetUUID(r, db) != tools.GetElement(db, "BASKET", "account", "uuid", basket_) {
		tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		return
	}

	// Checking if the value is empty
	if tools.ValueIsEmpty(basket_, equipment_, number_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the basket exists
	if !tools.ElementExists(db, "BASKET", "uuid", basket_) {
		tools.JsonResponse(w, 404, `{"message": "Basket not found"}`)
		return
	}

	// Checking if the account don't have a unpaid basket
	rows, err := db.Query("SELECT * FROM BASKET WHERE uuid = ? AND paid = false", basket_)
	if err != nil {
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		tools.JsonResponse(w, 400, `{"message": "Basket already paid"}`)
		return
	}

	// Checking if the equipment exists
	if !tools.ElementExists(db, "EQUIPMENT", "uuid", equipment_) {
		tools.JsonResponse(w, 404, `{"message": "Equipment not found"}`)
		return
	}

	// checking if the number choosed is not superior to the number of equipment available
	nbEquipment := tools.GetElement(db, "EQUIPMENT", "number", "uuid", equipment_)
	if number_ > nbEquipment {
		tools.JsonResponse(w, 400, `{"message": "Number of equipment choosed is superior to the number of equipment available"}`)
		return
	}

	// check if the equipment is already in the basket
	if tools.ElementExistsInLinkTable(db, "BASKET_EQUIPMENT", "basket", basket_, "equipment", equipment_) {
		tools.JsonResponse(w, 400, `{"message": "Equipment already in basket"}`)
		return
	}

	// Inserting the basket_housing
	_, err = db.Exec("INSERT INTO BASKET_EQUIPMENT (basket, equipment, number) VALUES (?, ?, ?)", basket_, equipment_, number_)
	if err != nil {
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	jsonResponse := `{"message": "Equipment added to basket"}`
	tools.JsonResponse(w, 201, jsonResponse)

}

func BasketEquipmentDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	query := tools.ReadQuery(r)

	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `basket`, `equipment`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	basket_ := query["basket"]
	equipment_ := query["equipment"]

	// Checking if the value is empty
	if tools.ValueIsEmpty(basket_, equipment_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the basket exists
	if !tools.ElementExists(db, "BASKET", "uuid", basket_) {
		tools.JsonResponse(w, 404, `{"message": "Basket not found"}`)
		return
	}

	if tools.GetUUID(r, db) != tools.GetElement(db, "BASKET", "account", "uuid", basket_) {
		tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
		return
	}

	// Checking if the housing exists
	if !tools.ElementExists(db, "EQUIPMENT", "uuid", equipment_) {
		tools.JsonResponse(w, 404, `{"message": "Equipment not found"}`)
		return
	}

	// Checking if the basket_housing exists
	if !tools.ElementExistsInLinkTable(db, "BASKET_EQUIPMENT", "basket", basket_, "equipment", equipment_) {
		tools.JsonResponse(w, 404, `{"message": "Equipment not in basket"}`)
		return
	}

	// Deleting the basket_housing
	_, err := db.Exec("DELETE FROM BASKET_EQUIPMENT WHERE basket = ? AND equipment = ?", basket_, equipment_)
	if err != nil {
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	jsonResponse := `{"message": "Equipment removed from basket"}`
	tools.JsonResponse(w, 200, jsonResponse)

}