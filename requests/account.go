package requests

import (
	"net/http"
	"tools"
	"slices"
)

func Account(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		post(w, r)
	case "GET":
		get(w, r)
	case "PUT":
		put(w, r)
	case "DELETE":
		delete(w, r)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func post(w http.ResponseWriter, r *http.Request) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, "username", "password") {
		tools.JsonResponse(w, 400, `{"message": "Missing username or password"}`)
		return
	}

	username := tools.BodyValueToString(body, "username")
	password := tools.BodyValueToString(body, "password")

	// Checking if the values are empty
	if tools.ValueIsEmpty(username, password) {
		tools.JsonResponse(w, 400, `{"message": "Empty username or password"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooShort(8, username, password) {
		tools.JsonResponse(w, 400, `{"message": "Username or password too short"}`)
		return
	}

	// Checking if the values are too short or too long
	if tools.ValueTooLong(32, username, password) {
		tools.JsonResponse(w, 400, `{"message": "Username or password too long"}`)
		return
	}

	// Checking if the password is strong enough
	if tools.PasswordNotStrong(password) {
		tools.JsonResponse(w, 400, `{"message": "Password not strong enough"}`)
		return
	}
	
	// Creating the response
	jsonResponse := `{"message": "Account created"`
	returnFields := tools.GetReturnFields(r)
	
	// Adding the return fields asked by the user
	if slices.Contains(returnFields, "username") {
		jsonResponse += `, "username": "` + username + `"`
	}
	if slices.Contains(returnFields, "password") {
		jsonResponse += `, "password": "` + password + `"`
	}

	jsonResponse += "}"

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse)

}

func get(w http.ResponseWriter, r *http.Request) {
	tools.JsonResponse(w, 200, `{"message": "Account get"}`)
}

func put(w http.ResponseWriter, r *http.Request) {
	tools.JsonResponse(w, 200, `{"message": "Account put"}`)
}

func delete(w http.ResponseWriter, r *http.Request) {
	tools.JsonResponse(w, 200, `{"message": "Account delete"}`)
}