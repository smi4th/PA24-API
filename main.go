package main

import (
	"net/http"
	"requests"
	"tools"
)

func main() {

	// Initialize the database connection
	db := tools.InitDatabaseConnection()
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

		// Check if the user is authenticated
		if !tools.IsAuthenticated(r, db) {
			tools.JsonResponse(w, 401, `{"message": "Unauthorized"}`)
			return
		}

		switch r.URL.Path {
		case "/account_bedroom":
			requests.AccountBedroom(w, r, db)
		case "/account_housing":
			requests.AccountHousing(w, r, db)
		case "/account_services":
			requests.AccountServices(w, r, db)
		case "/account_subscription":
			requests.AccountSubscription(w, r, db)
		case "/account_type":
			requests.AccountType(w, r, db)
		case "/account":
			requests.Account(w, r, db)
		case "/bed_room":
			requests.BedRoom(w, r, db)
		case "/consume":
			requests.Consume(w, r, db)
		case "/disponibility_account":
			requests.DisponibilityAccount(w, r, db)
		case "/disponibility":
			requests.Disponibility(w, r, db)
		case "/equipment_type":
			requests.EquipmentType(w, r, db)
		case "/equipment":
			requests.Equipment(w, r, db)
		case "/house_type":
			requests.HouseType(w, r, db)
		case "/housing_equipment":
			requests.HousingEquipment(w, r, db)
		case "/housing":
			requests.Housing(w, r, db)
		case "/message":
			requests.Message(w, r, db)
		case "/provider_account":
			requests.ProviderAccount(w, r, db)
		case "/provider":
			requests.Provider(w, r, db)
		case "/reservation_berdroom":
			requests.ReservationBedroom(w, r, db)
		case "/reservation_housing":
			requests.ReservationHousing(w, r, db)
		case "/services_types":
			requests.ServicesTypes(w, r, db)
		case "/services":
			requests.Services(w, r, db)
		case "/subscription":
			requests.Subscription(w, r, db)
		default:
			tools.JsonResponse(w, 404, `{"message": "Not found"}`)
		}
	})

	tools.InfoLog("Server is running on port 80")
	http.ListenAndServe(":80", nil)

}
