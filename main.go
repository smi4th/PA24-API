package main

import (
	"net/http"
	"requests"
	"tools"
	"fmt"
)

func main() {

	// Initialize the database connection
	db := tools.InitDatabaseConnection()
	if db == nil {
		fmt.Println("Failed to connect to the database")
		tools.ErrorLog("Failed to connect to the database")
		return
	}
	defer tools.CloseDatabaseConnection(db)

	// Handle the requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/login"{
			requests.Login(w, r, db)
			return
		}

		if r.URL.Path == "/account" && r.Method == "POST" {
			requests.AccountPost(w, r, db)
			return
		}

		if r.URL.Path == "/account_type" && r.Method == "GET" && r.URL.Query().Get("private") == "false" {
			requests.AccountTypeGet(w, r, db)
			return
		}

		// Check if the user is authenticated
		if !tools.IsAuthenticated(r, db) {
			tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
			return
		}

		switch r.URL.Path {
		case "/account_subscription":
			requests.AccountSubscription(w, r, db)
		case "/account_type":
			requests.AccountType(w, r, db)
		case "/account":
			requests.Account(w, r, db)
		case "/account/verifyPassword":
			requests.VerifyPassword(w, r, db)
		case "/admin":
			requests.Admin(w, r, db)
		case "/basket":
			requests.Basket(w, r, db)
		case "/basket/housing":
			requests.BasketHousing(w, r, db)
		case "/basket/bedroom":
			requests.BasketBedroom(w, r, db)
		case "/basket/services":
			requests.BasketServices(w, r, db)
		case "/basket/equipment":
			requests.BasketEquipment(w, r, db)
		case "/bed_room":
			requests.BedRoom(w, r, db)
		case "/bed_room/available":
			requests.BedroomAvailable(w, r, db)
		case "/bed_room/reservation":
			requests.BedroomReservation(w, r, db)
		case "/disponibility":
			requests.Disponibility(w, r, db)
		case "/equipment_type":
			requests.EquipmentType(w, r, db)
		case "/equipment":
			requests.Equipment(w, r, db)
		case "/house_type":
			requests.HouseType(w, r, db)
		case "/housing":
			requests.Housing(w, r, db)
		case "/housing/available":
			requests.HousingAvailable(w, r, db)
		case "/housing/reservation":
			requests.HousingReservation(w, r, db)
		case "/message":
			requests.Message(w, r, db)
		case "/provider":
			requests.Provider(w, r, db)
		case "/review":
			requests.Review(w, r, db)
		case "/services_types":
			requests.ServicesTypes(w, r, db)
		case "/services":
			requests.Services(w, r, db)
		case "/status":
			requests.Status(w, r, db)
		case "/subscription":
			requests.Subscription(w, r, db)
		case "/ticket":
			requests.Ticket(w, r, db)
		case "/tmessage":
			requests.TMessage(w, r, db)
		default:
			tools.JsonResponse(w, 404, `{"message": "Not found"}`)
		}
	})

	tools.InfoLog("Server is running on port 80")
	http.ListenAndServe(":80", nil)

}
