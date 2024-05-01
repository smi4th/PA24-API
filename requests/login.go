package requests

import (
	"net/http"
	"tools"
	"database/sql"
)

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case "POST":
		LoginPost(w, r, db)
	default:
		tools.JsonResponse(w, 405, `{"message": "Method not allowed"}`)
	}
}

func LoginPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// Getting the body of the request
	body := tools.ReadBody(r)
	
	tools.RequestLog(r, body)

	// Checking if the body contains the required fields
	if tools.ValuesNotInBody(body, `email`, `password`) {
		tools.JsonResponse(w, 400, `{"message": "Missing fields"}`)
		return
	}

    email_ := tools.BodyValueToString(body, `email`)
    password_ := tools.BodyValueToString(body, `password`)
	

	// Checking if the values are empty
	if tools.ValueIsEmpty(email_, password_) {
		tools.JsonResponse(w, 400, `{"message": "Fields cannot be empty"}`)
		return
	}

    // Checking if the email exists in the database
    result, err := tools.ExecuteQuery(db, "SELECT email, password FROM account WHERE email = ?", email_)
    if err != nil {
        tools.ErrorLog(err.Error())
        tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
        return
    }
    defer result.Close()

    // Checking if the email exists
    if !result.Next() {
        tools.JsonResponse(w, 404, `{"message": "Email not found"}`)
        return
    }

    // Getting the password from the database
    var password, email string
    result.Scan(&email, &password)
    
    // Checking if the password is correct
    if tools.ComparePassword(password, password_) == false {
        tools.JsonResponse(w, 401, `{"message": "Invalid password"}`)
        return
    }

    // Generating the token
	token_ := tools.GenerateToken()

	// Inserting the Housing in the database
	result, err = tools.ExecuteQuery(db, "UPDATE account SET token = ? WHERE email = ?", token_, email_)
	if err != nil {
		tools.ErrorLog(err.Error())
		tools.JsonResponse(w, 500, `{"message": "Internal server error"}`)
		return
	}
	defer result.Close()

	// Creating the response
	jsonResponse := `{"message": "Login successful", "token": "` + token_ + `"}`

	// Sending the response
	tools.JsonResponse(w, 201, jsonResponse)

}