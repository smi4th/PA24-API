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

/*
#######################################
########## Config Functions ###########
#######################################
*/

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
	configMap["database"]["host"] = config.Database.Host
	configMap["database"]["port"] = config.Database.Port
	configMap["database"]["username"] = config.Database.Username
	configMap["database"]["password"] = config.Database.Password
	configMap["database"]["database"] = config.Database.Database

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

func BodyValueToString(body map[string]interface{}, key string) string {
	value, ok := body[key].(string)
	if !ok {
		ErrorLog("Value is not a string")
		return ""
	}
	return value
}

func BodyValueToInt(body map[string]interface{}, key string) int {
	value, ok := body[key].(float64)
	if !ok {
		ErrorLog("Value is not an integer")
		return 0
	}
	return int(value)
}

func BodyValueToFloat(body map[string]interface{}, key string) float64 {
	value, ok := body[key].(float64)
	if !ok {
		ErrorLog("Value is not a float")
		return 0
	}
	return value
}

func BodyValueToBool(body map[string]interface{}, key string) bool {
	value, ok := body[key].(bool)
	if !ok {
		ErrorLog("Value is not a boolean")
		return false
	}
	return value
}

func BodyValueToMap(body map[string]interface{}, key string) map[string]interface{} {
	value, ok := body[key].(map[string]interface{})
	if !ok {
		ErrorLog("Value is not a map")
		return nil
	}
	return value
}

func BodyValueToArray(body map[string]interface{}, key string) []interface{} {
	value, ok := body[key].([]interface{})
	if !ok {
		ErrorLog("Value is not an array")
		return nil
	}
	return value
}

/*
#######################################
########## Database Functions #########
#######################################
*/

func InitDatabaseConnection() *sql.DB {
	config := getConfig()
	db, err := sql.Open("mysql", config["database"]["username"] + ":" + config["database"]["password"] + "@tcp(" + config["database"]["host"] + ":" + config["database"]["port"] + ")/" + config["database"]["database"])
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

func ExecuteQuery(db *sql.DB, query string, args ...interface{}) *sql.Rows {
	InfoLog("Preparing query: " + query)
	stmt, err := db.Prepare(query)
	if err != nil {
		ErrorLog(err.Error())
		return nil
	}
	defer stmt.Close()

	if len(args) == 0 {
		InfoLog("Executing query without arguments")
		rows, err := stmt.Query()
		if err != nil {
			ErrorLog(err.Error())
			return nil
		}
		return rows
	}

	InfoLog("Executing query with arguments ")
	for i, arg := range args {
		InfoLog("Arg " + strconv.Itoa(i) + ": " + arg.(string))
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		ErrorLog(err.Error())
		return nil
	}

	return rows
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
		if len(value) < length {
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

/*
#######################################
########## Request functions ##########
#######################################
*/

func GetReturnFields(r *http.Request) []string {
	// _return_fields is a comma separated list of fields
	fields := strings.Split(r.URL.Query().Get("_return_fields"), ",")
	return fields
}