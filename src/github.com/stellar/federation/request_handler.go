package federation

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type RequestHandler struct {
	config   *Config
	database Database
}

func (rh *RequestHandler) Main(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	requestType := r.URL.Query().Get("type")
	q := r.URL.Query().Get("q")
	switch {
	case requestType == "name" && q != "":
		rh.FedDBRequest(q, w)
	case requestType == "id" && q != "":
		rh.RevFedDBRequest(q, w)
	default:
		http.Error(w, ErrorResponseString("invalid_request", "Invalid request"), http.StatusBadRequest)
	}
}

func (rh *RequestHandler) RevFedDBRequest(accountID string, w http.ResponseWriter) {
	record := FedRecord{}

	record.AccountId = accountID

	revResult := RevFedRecord{}

	err := rh.database.Get(&revResult, rh.config.ReverseFederationQuery, accountID)

	if checkDBErr(err, w) {
		record.StellarAddress = revResult.Name + "*" + rh.config.Domain
		rh.WriteResponse(record, w)
	}
}

func (rh *RequestHandler) FedDBRequest(stellarAddress string, w http.ResponseWriter) {
	record := FedRecord{}
	var name string
	domain := ""

	if i := strings.Index(stellarAddress, "*"); i >= 0 {
		name = stellarAddress[:i]
		domain = stellarAddress[i+1:]
	}

	if name == "" || domain != rh.config.Domain {
		http.Error(w, ErrorResponseString("not_found", "Incorrect Domain"), http.StatusNotFound)
		return
	}

	err := rh.database.Get(&record, rh.config.FederationQuery, name)
	record.StellarAddress = stellarAddress

	if checkDBErr(err, w) {
		rh.WriteResponse(record, w)
	}
}

// returns false if there was an error
func checkDBErr(err error, w http.ResponseWriter) bool {
	if err != nil {
		if err.Error() == "sql: no rows in result set" {

			log.Print("Federation record NOT found")
			http.Error(w, ErrorResponseString("not_found", "Account not found"), http.StatusNotFound)
		} else {
			log.Print("Server error: ", err)
			http.Error(w, ErrorResponseString("server_error", "Server error"), http.StatusInternalServerError)
		}
		return false
	}
	return true
}

func (rh *RequestHandler) WriteResponse(record FedRecord, w http.ResponseWriter) {
	log.Print("Federation record found")

	json, err := json.Marshal(record)

	if err != nil {
		http.Error(w, ErrorResponseString("server_error", "Server error"), http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

func ErrorResponseString(code string, message string) string {
	error := Error{Code: code, Message: message}
	json, _ := json.Marshal(error)
	return string(json)
}
