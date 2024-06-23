package tools

import (
	"os"
	"time"
	"net/http"
	"strconv"
	"io/ioutil"
	"database/sql"
	"encoding/json"
	"strings"
	"fmt"
	"unicode"
	"crypto/rand"
	
	"golang.org/x/crypto/bcrypt"
	_ "github.com/go-sql-driver/mysql"

)

/*
#######################################
########## Logging Functions ##########
#######################################
*/

func log(color string, message string) {
	os.Stdout.WriteString(color + time.Now().Format("01-02-2006 15:04:05") + " " + message + "\n")
	os.Stdout.WriteString("\033[0m")

	file := "logs/" + time.Now().Format("01-02-2006") + ".log"
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		ErrorLog(err.Error())
	}
	defer f.Close()

	f.WriteString(time.Now().Format("01-02-2006 15:04:05") + " " + message + "\n")

}

func ErrorLog(message string) {
	if configMap == nil {
		getConfig()
	}
	// if the log level contains "error"
	if strings.Contains(configMap["logs"]["level"], "error") {
		log("\033[31m", "[ERROR] " + message)
	}
}

func InfoLog(message string) {
	if configMap == nil {
		getConfig()
	}
	// if the log level contains "info"
	if strings.Contains(configMap["logs"]["level"], "info") {
		log("\033[32m", "[INFO] " + message)
	}
}

func ResponseLog(status int, body string) {
	if configMap == nil {
		getConfig()
	}
	// if the log level contains "response"
	if strings.Contains(configMap["logs"]["level"], "response") {
		log("\033[33m", "[RESPONSE] " + strconv.Itoa(status) + " " + body)
	}
}

func RequestLog(r *http.Request, body map[string]interface{}) {
	if configMap == nil {
		getConfig()
	}
	// if the log level contains "request"
	if strings.Contains(configMap["logs"]["level"], "request") {
		log("\033[34m", "[REQUEST] " + r.Method + " " + r.URL.Path + "\n" + fmt.Sprintf("%v", body))
	}
}

func SQLLog(query string) {
	if configMap == nil {
		getConfig()
	}
	// if the log level contains "sql"
	if strings.Contains(configMap["logs"]["level"], "sql") {
		log("\033[35m", "[SQL] " + query)
	}
}

/*
#######################################
########## Config Functions ###########
#######################################
*/

type configJson struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"database"`

	Logs struct {
		Path string `json:"path"`
		Level []string `json:"level"`
	} `json:"logs"`
}

var configMap map[string]map[string]string

func getConfig() map[string]map[string]string {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		ErrorLog(err.Error())
		return nil
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		ErrorLog(err.Error())
		return nil
	}

	var config configJson
	json.Unmarshal(byteValue, &config)

	configMap = make(map[string]map[string]string)
	configMap["database"] = make(map[string]string)
	configMap["database"]["host"] = os.Getenv("DB_HOST")
	configMap["database"]["port"] = os.Getenv("DB_PORT")
	configMap["database"]["username"] = os.Getenv("MYSQL_USER")
	configMap["database"]["database"] = os.Getenv("MYSQL_DATABASE")

	// the password is a file stored "/run/secrets/mysql_password"
	password, err := ioutil.ReadFile(os.Getenv("MYSQL_PASSWORD"))
	if err != nil {
		ErrorLog(err.Error())
		return nil
	}
	configMap["database"]["password"] = string(password)

	configMap["logs"] = make(map[string]string)
	configMap["logs"]["path"] = config.Logs.Path
	configMap["logs"]["level"] = ""
	for _, level := range config.Logs.Level {
		configMap["logs"]["level"] += level + " "
	}

	InfoLog("Config loaded")

	return configMap

}

/*
#######################################
########## JSON Functions #############
#######################################
*/

func JsonResponse(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(body))

	ResponseLog(status, body)
}

func ReadBody(r *http.Request) map[string]interface{} {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorLog(err.Error())
		return nil
	}

	var jsonMap map[string]interface{}
	json.Unmarshal(body, &jsonMap)

	return jsonMap
}

func ReadQuery(r *http.Request) map[string]string {
	query := r.URL.Query()
	queryMap := make(map[string]string)
	for key, value := range query {
		queryMap[key] = value[0]
	}
	return queryMap
}

func BodyValueToString(body map[string]interface{}, key string) string {
	value, ok := body[key].(string)
	if !ok {
		return ""
	}
	return value
}

func BodyValueToInt(body map[string]interface{}, key string) int {
	value, ok := body[key].(float64)
	if !ok {
		return 0
	}
	return int(value)
}

func BodyValueToFloat(body map[string]interface{}, key string) float64 {
	value, ok := body[key].(float64)
	if !ok {
		return 0
	}
	return value
}

func BodyValueToBool(body map[string]interface{}, key string) bool {
	value, ok := body[key].(bool)
	if !ok {
		return false
	}
	return value
}

func BodyValueToMap(body map[string]interface{}, key string) map[string]interface{} {
	value, ok := body[key].(map[string]interface{})
	if !ok {
		return nil
	}
	return value
}

func BodyValueToArray(body map[string]interface{}, key string) []interface{} {
	value, ok := body[key].([]interface{})
	if !ok {
		return nil
	}
	return value
}

/*
#######################################
########## Global Functions ###########
#######################################
*/

func getLimit(query map[string]string) int {
	limit, err := strconv.Atoi(query["limit"])
	if err != nil {
		return 0
	}
	return limit
}

func getOffset(query map[string]string) int {
	offset, err := strconv.Atoi(query["offset"])
	if err != nil {
		return 0
	}
	return offset
}

func IsAuthenticated(r *http.Request, db *sql.DB) bool {
	token := r.Header.Get("Authorization")
	if token == "" {
		InfoLog("No token provided")
		return false
	}
	return ElementExists(db, "account", "token", strings.Replace(token, "Bearer ", "", 1))
}

func IsAdmin(r *http.Request, db *sql.DB) bool {
	uuid := GetUUID(r, db)
	request := "SELECT `admin` FROM `ACCOUNT_TYPE` WHERE `uuid` = (SELECT `account_type` FROM `ACCOUNT` WHERE `uuid` = ?)"
	rows, err := ExecuteQuery(db, request, uuid)
	if err != nil {
		ErrorLog(err.Error())
		return false
	}
	defer rows.Close()

	if rows.Next() {
		var admin string
		err := rows.Scan(&admin)
		if err != nil {
			ErrorLog(err.Error())
			return false
		}
		return admin == "true"
	}

	return false
}

func GetElement(db *sql.DB, table string, attribute string, pkAttribute string, pkValue string) string {
	// Execute the query to get the element from the specified table and attribute.
	rows, err := ExecuteQuery(db, "SELECT `" + attribute + "` FROM `" + table + "` WHERE `" + pkAttribute + "` = ?", pkValue)
	if err != nil {
		ErrorLog(err.Error())
		return ""
	}
	defer rows.Close()

	// Check if the element exists.
	if rows.Next() {
		var element string
		err := rows.Scan(&element)
		if err != nil {
			ErrorLog(err.Error())
			return ""
		}
		return element
	}

	return ""

}

func GetElementFromLinkTable(db *sql.DB, table string, attribute string, idAttribute1 string, idValue1 string, idAttribute2 string, idValue2 string) string {
	// Execute the query to get the element from the specified table and attribute.
	rows, err := ExecuteQuery(db, "SELECT `" + attribute + "` FROM `" + table + "` WHERE `" + idAttribute1 + "` = ? AND `" + idAttribute2 + "` = ?", idValue1, idValue2)
	if err != nil {
		ErrorLog(err.Error())
		return ""
	}
	defer rows.Close()

	// Check if the element exists.
	if rows.Next() {
		var element string
		err := rows.Scan(&element)
		if err != nil {
			ErrorLog(err.Error())
			return ""
		}
		return element
	}

	return ""

}

func GetUUID(r *http.Request, db *sql.DB) string {
	token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	rows, err := ExecuteQuery(db, "SELECT `uuid` FROM `ACCOUNT` WHERE `token` = ?", token)
	if err != nil {
		ErrorLog(err.Error())
		return ""
	}
	defer rows.Close()

	if rows.Next() {
		var uuid string
		err := rows.Scan(&uuid)
		if err != nil {
			ErrorLog(err.Error())
			return ""
		}
		return uuid
	}

	return ""
}

/*
#######################################
########## Database Functions #########
#######################################
*/

func InitDatabaseConnection() *sql.DB {
	config := getConfig()
	db, err := sql.Open("mariadb", config["database"]["username"] + ":" + config["database"]["password"] + "@tcp(" + config["database"]["host"] + ":" + config["database"]["port"] + ")/" + config["database"]["database"])
	if err != nil {
		ErrorLog(err.Error())
		return nil
	}
	err = db.Ping()
	if err != nil {
		ErrorLog(err.Error())
		return nil
	}

	InfoLog("Connection with " + config["database"]["database"] + "@" + config["database"]["host"] + " is successful")

	return db
}

func CloseDatabaseConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		ErrorLog(err.Error())
		return
	}
	InfoLog("Connection with database is closed")
}

func ExecuteQuery(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	query = strings.ToUpper(query)
	SQLLog("Preparing query: " + query)
	stmt, err := db.Prepare(query)
	if err != nil {
		ErrorLog(err.Error())
		return nil, err
	}
	defer stmt.Close()

	if len(args) == 0 {
		SQLLog("Executing query without arguments")
		rows, err := stmt.Query()
		if err != nil {
			ErrorLog(err.Error())
			return nil, err
		}
		return rows, nil
	}

	SQLLog("Executing query with arguments ")
	for i, arg := range args {
		SQLLog("Arg " + strconv.Itoa(i) + ": " + arg.(string))
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		ErrorLog(err.Error())
		return nil, err
	}

	return rows, nil
}

func RowsToJson(rows *sql.Rows) string {
	columns, err := rows.Columns()
	if err != nil {
		ErrorLog(err.Error())
		return ""
	}
	
	var result []map[string]string

	// for each row
	for rows.Next() {
		// create a map of columns
		columnsMap := make(map[string]string)
		// create a slice of pointers to interface{}
		columnsPointers := make([]interface{}, len(columns))
		// fill the slice with pointers to the values of the columns
		for i := range columns {
			columnsPointers[i] = new(interface{})
		}
		// scan the row and fill the columnsPointers with the values
		err = rows.Scan(columnsPointers...)
		if err != nil {
			ErrorLog(err.Error())
			return ""
		}
		// fill the columnsMap with the values of the columns
		for i, column := range columns {
			val := columnsPointers[i].(*interface{})
			// int64
			if v, ok := (*val).(int64); ok {
				columnsMap[column] = strconv.FormatInt(v, 10)
			// float64
			} else if v, ok := (*val).(float64); ok {
				columnsMap[column] = strconv.FormatFloat(v, 'f', -1, 64)
			// bool
			} else if v, ok := (*val).(bool); ok {
				columnsMap[column] = strconv.FormatBool(v)
			// string
			} else if v, ok := (*val).(string); ok {
				columnsMap[column] = v
			// []byte
			} else {
				ErrorLog("Unknown type")
				return ""
			}
		}
		// append the columnsMap to the result
		result = append(result, columnsMap)
	}

	// convert the result to a JSON string
	jsonResult, err := json.Marshal(result)
	if err != nil {
		ErrorLog(err.Error())
		return ""
	}

	return string(jsonResult)

}

/*
#######################################
########## Data functions #############
#######################################
*/

func GetDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func DateBefore(date1 string, date2 string) bool {
	layout := "2006-01-02 15:04:05"
	t1, err := time.Parse(layout, date1)
	if err != nil {
		ErrorLog(err.Error())
		return false
	}
	t2, err := time.Parse(layout, date2)
	if err != nil {
		ErrorLog(err.Error())
		return false
	}
	return t1.Before(t2)
}

func DateAfter(date1 string, date2 string) bool {
	layout := "2006-01-02 15:04:05"
	t1, err := time.Parse(layout, date1)
	if err != nil {
		ErrorLog(err.Error())
		return false
	}
	t2, err := time.Parse(layout, date2)
	if err != nil {
		ErrorLog(err.Error())
		return false
	}
	return t1.After(t2)
}

func PeriodeOverlap(db *sql.DB, table string, start1Attribute string, end1Attribute string, pkAttribute string, pkValue string, start2Value string, end2Value string) bool {
	// Execute the query to check if the periode overlaps in the specified table.
	rows, err := ExecuteQuery(db, "SELECT COUNT(*) as `count` FROM `" + table + "` WHERE `" + start1Attribute + "` <= ? AND `" + end1Attribute + "` >= ? AND `" + pkAttribute + "` = ?", end2Value, start2Value, pkValue)
	if err != nil {
		ErrorLog(err.Error())
		return false
	}
	defer rows.Close()

	// Check if the periode overlaps.
	if rows.Next() {
		var count int
		err := rows.Scan(&count)
		if err != nil {
			ErrorLog(err.Error())
			return false
		}
		return count > 0
	}

	return false
}

func GenerateUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		ErrorLog(err.Error())
		return ""
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func GenerateToken() string {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		ErrorLog(err.Error())
		return ""
	}
	return fmt.Sprintf("%x", token)
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ErrorLog(err.Error())
		return ""
	}
	return string(hash)
}

func ComparePassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		ErrorLog(err.Error())
		return false
	}
	return true
}

func ValuesNotInBody(body map[string]interface{}, keys ...string) bool {
	for _, key := range keys {
		if _, ok := body[key]; !ok {
			return true
		}
	}
	return false
}

func ValuesNotInQuery(query map[string]string, keys ...string) bool {
	for _, key := range keys {
		if query[key] == "" {
			return true
		}
	}
	return false
}

func AtLeastOneValueInQuery(query map[string]string, keys ...string) bool {
	for _, key := range keys {
		if query[key] != "" {
			return false
		}
	}
	return true
}

func AtLeastOneValueInBody(body map[string]interface{}, keys ...string) bool {
	for _, key := range keys {
		if _, ok := body[key]; ok {
			return false
		}
	}
	return true
}

func AtLeastOneValueNotEmpty(values ...string) bool {
	for _, value := range values {
		if value != "" {
			return false
		}
	}
	return true
}

func ValueIsEmpty(values ...string) bool {
	for _, value := range values {
		if value == "" {
			return true
		}
	}
	return false
}

func ValueTooShort(length int, values ...string) bool {
	for _, value := range values {
		if len(value) < length && value != "" {
			return true
		}
	}
	return false
}

func ValueTooLong(length int, values ...string) bool {
	for _, value := range values {
		if len(value) > length {
			return true
		}
	}
	return false
}

func PasswordNotStrong(password string) bool {
	var (
        hasUpper   = false
        hasLower   = false
        hasNumber  = false
        hasSpecial = false
    )
	
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
		}
    }

    return !(!ValueTooShort(8, password) && hasUpper && hasLower && hasNumber && hasSpecial)
}

func ValueInArray(value string, array ...string) bool {
	for _, element := range array {
		if element == value {
			return true
		}
	}
	return false
}

func ElementExists(db *sql.DB, table string, attribute string, value string) bool {
	// Execute the query to count occurrences of the value in the specified table and attribute.
	InfoLog("Testing if " + value + " exists in " + table + "." + attribute)
	result, err := ExecuteQuery(db , "SELECT COUNT(`" + attribute + "`) as `count` FROM `" + table + "` WHERE `" + attribute + "` = ?", value)
	if err != nil {
		ErrorLog(err.Error())
		return false
	}
	defer result.Close()

	// Check if the `count` is greater than 0.
	if result.Next() {
		var count int
		err := result.Scan(&count)
		if err != nil {
			ErrorLog(err.Error())
			return false
		}

		InfoLog("Count: " + strconv.Itoa(count))

		return count > 0
	}

	return false
}

func ElementExistsInLinkTable(db *sql.DB, table string, attribute1 string, value1 string, attribute2 string, value2 string) bool {
	// Execute the query to count occurrences of the value in the specified table and attributes.
	result, err := ExecuteQuery(db , "SELECT COUNT(`" + attribute1 + "`) as `count` FROM `" + table + "` WHERE `" + attribute1 + "` = ? AND `" + attribute2 + "` = ?", value1, value2)
	if err != nil {
		ErrorLog(err.Error())
		return false
	}
	defer result.Close()

	// Check if the `count` is greater than 0.
	if result.Next() {
		var count int
		err := result.Scan(&count)
		if err != nil {
			ErrorLog(err.Error())
			return false
		}
		return count > 0
	}

	return false
}

func EmailIsValid(email string) bool {
	// Check if the email contains an @.
	if !strings.Contains(email, "@") {
		return false
	}
	// Check if the email contains a dot.
	if !strings.Contains(email, ".") {
		return false
	}
	// Check if the email is not too short.
	if len(email) < 5 {
		return false
	}
	// Check if the email is not too long.
	if len(email) > 64 {
		return false
	}
	return true
}

// appendLikeCondition appends a 'LIKE' condition to the SQL query if the provided value is non-empty.
// It updates both the SQL query and the parameters slice accordingly.
func AppendLikeCondition(request *string, params *[]interface{}, field, value string) {
    if value != "" {
        // Append the 'LIKE' condition to the SQL query and add the parameter to the slice.
        *request += " " + field + " LIKE ? OR"
        *params = append(*params, "%"+value+"%")
    }
}

// appendCondition appends a condition to the SQL query if the provided value is non-empty or if it's a non-strict search.
// It updates both the SQL query and the parameters slice accordingly.
func AppendCondition(request *string, params *[]interface{}, field, value string, strictSearch bool) {
    if value != "" && !ValueInArray(field, "strictSearch", "limit", "offset") {
        // Determine the logical operator based on strict or non-strict search.
        operator := "OR"
        if strictSearch {
            operator = "AND"
        }

        // Additional condition for 'start_date' and 'end_date'.
        if field == "start_date" { 
            *request += " " + field + "<=? " + operator
        } else if field == "end_date" {
            *request += " " + field + ">=? " + operator
        } else {
            *request += " " + field + "=? " + operator
        }

        // Append the parameter to the slice.
        *params = append(*params, value)
    }
}

func AppendUpdate(request *string, params *[]interface{}, key string, value interface{}) {
	// Append the key and value to the SQL query and add the parameter to the slice.
	*request += "`" + key + "` = ?, "
	*params = append(*params, value)
}