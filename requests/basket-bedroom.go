package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func BasketBedroom(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		BasketBedroomPost(w, r, db)
	case "DELETE":
		BasketBedroomDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func BasketBedroomPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	body := tools.ReadBody(r)

	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `basket`, `bedroom`, `start_time`, `end_time`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	basket_ := tools.BodyValueToString(body, `basket`)
	bedroom_ := tools.BodyValueToString(body, `bedroom`)
	start_time_ := tools.BodyValueToString(body, `start_time`)
	end_time_ := tools.BodyValueToString(body, `end_time`)

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
	if tools.ValueIsEmpty(basket_, bedroom_, start_time_, end_time_) {
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

	// Checking if the bedroom exists
	if !tools.ElementExists(db, "BED_ROOM", "uuid", bedroom_) {
		tools.JsonResponse(w, 404, `{"message": "Bedroom not found"}`)
		return
	}

	// Checking if the housing is available during the period
	if tools.PeriodeOverlap(db, "BASKET_HOUSING", "start_time", "end_time", "bedroom", bedroom_, start_time_, end_time_) {
		tools.JsonResponse(w, 400, `{"message": "Bedroom not available during this period"}`)
		return
	}

	// Checking if the basket_bedroom already exists
	if tools.ElementExistsInLinkTable(db, "BASKET_BEDROOM", "basket", basket_, "bedroom", bedroom_) {
		tools.JsonResponse(w, 400, `{"message": "Bedroom already in basket"}`)
		return
	}

	// Inserting the basket_housing
	_, err = db.Exec("INSERT INTO BASKET_BEDROOM (basket, bedroom, start_time, end_time) VALUES (?, ?, ?, ?)", basket_, bedroom_, start_time_, end_time_)
	if err != nil {
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	jsonResponse := `{"message": "Bedroom added to basket"}`
	tools.JsonResponse(w, 201, jsonResponse)

}

func BasketBedroomDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	query := tools.ReadQuery(r)

	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `basket`, `bedroom`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	basket_ := query["basket"]
	bedroom_ := query["bedroom"]

	// Checking if the value is empty
	if tools.ValueIsEmpty(basket_, bedroom_) {
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
	if !tools.ElementExists(db, "BED_ROOM", "uuid", bedroom_) {
		tools.JsonResponse(w, 404, `{"message": "Bedroom not found"}`)
		return
	}

	// Checking if the basket_housing exists
	if !tools.ElementExistsInLinkTable(db, "BASKET_BEDROOM", "basket", basket_, "bedroom", bedroom_) {
		tools.JsonResponse(w, 404, `{"message": "Bedroom not in basket"}`)
		return
	}

	// Deleting the basket_housing
	_, err := db.Exec("DELETE FROM BASKET_BEDROOM WHERE basket = ? AND bedroom = ?", basket_, bedroom_)
	if err != nil {
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	jsonResponse := `{"message": "Bedroom removed from basket"}`
	tools.JsonResponse(w, 200, jsonResponse)

}