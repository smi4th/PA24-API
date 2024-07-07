package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func BasketServices(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		BasketServicesPost(w, r, db)
	case "DELETE":
		BasketServicesDelete(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func BasketServicesPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	body := tools.ReadBody(r)

	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `basket`, `services`, `start_time`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	basket_ := tools.BodyValueToString(body, `basket`)
	services_ := tools.BodyValueToString(body, `services`)
	start_time_ := tools.BodyValueToString(body, `start_time`)

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
	if tools.ValueIsEmpty(basket_, services_, start_time_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

	// Checking if the basket exists
	if !tools.ElementExists(db, "BASKET", "uuid", basket_) {
		tools.JsonResponse(w, 404, `{"message": "Basket not found"}`)
		return
	}

	// Checking if the account don't have a unpaid basket
	rows, err := db.Query("SELECT * FROM BASKET WHERE uuid = ? AND paid = 0", basket_)
	if err != nil {
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		tools.JsonResponse(w, 400, `{"message": "Basket already paid"}`)
		return
	}

	// Checking if the services exists
	if !tools.ElementExists(db, "SERVICES", "uuid", services_) {
		tools.JsonResponse(w, 404, `{"message": "Services not found"}`)
		return
	}

	// Checking if the basket_services already exists
	if tools.ElementExistsInLinkTable(db, "BASKET_SERVICE", "basket", basket_, "service", services_) {
		tools.JsonResponse(w, 400, `{"message": "Services already in basket"}`)
		return
	}

	// Inserting the basket_housing
	_, err = db.Exec("INSERT INTO BASKET_SERVICE (basket, service, start_time) VALUES (?, ?, ?)", basket_, services_, start_time_)
	if err != nil {
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	jsonResponse := `{"message": "Services added to basket"}`
	tools.JsonResponse(w, 201, jsonResponse)

}

func BasketServicesDelete(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	query := tools.ReadQuery(r)

	tools.RequestLog(r, tools.ReadBody(r))

	// Checking if the query contains the required fields
	if tools.ValuesNotInQuery(query, `basket`, `services`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

	basket_ := query["basket"]
	services_ := query["services"]

	// Checking if the value is empty
	if tools.ValueIsEmpty(basket_, services_) {
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
	if !tools.ElementExists(db, "SERVICES", "uuid", services_) {
		tools.JsonResponse(w, 404, `{"message": "Services not found"}`)
		return
	}

	// Checking if the basket_housing exists
	if !tools.ElementExistsInLinkTable(db, "BASKET_SERVICE", "basket", basket_, "service", services_) {
		tools.JsonResponse(w, 404, `{"message": "Services not in basket"}`)
		return
	}

	// Deleting the basket_housing
	_, err := db.Exec("DELETE FROM BASKET_SERVICE WHERE basket = ? AND service = ?", basket_, services_)
	if err != nil {
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}

	jsonResponse := `{"message": "Services removed from basket"}`
	tools.JsonResponse(w, 200, jsonResponse)

}