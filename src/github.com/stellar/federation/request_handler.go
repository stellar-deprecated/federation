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
	requestType := r.URL.Query().Get("type")
	q := r.URL.Query().Get("q")
	switch {
	case requestType == "name" && q != "":
		rh.FedDBRequest(q, w)
	case requestType == "id" && q != "":
		rh.RevFedDBRequest(q, w)
	case requestType == "txid" && q != "":
		rh.writeErrorResponse(w, ErrorResponseString("not_implemented", "txid requests are not supported"), http.StatusNotImplemented)
	default:
		rh.writeErrorResponse(w, ErrorResponseString("invalid_request", "Invalid request"), http.StatusBadRequest)
	}
}

func (rh *RequestHandler) RevFedDBRequest(accountID string, w http.ResponseWriter) {
	record := FedRecord{}

	record.AccountId = accountID

	revResult := RevFedRecord{}

	err := rh.database.Get(&revResult, rh.config.ReverseFederationQuery, accountID)

	if rh.checkDBErr(err, w) {
		record.StellarAddress = revResult.Name + "*" + rh.config.Domain
		rh.writeResponse(w, record)
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
		rh.writeErrorResponse(w, ErrorResponseString("not_found", "Incorrect Domain"), http.StatusNotFound)
		return
	}

	err := rh.database.Get(&record, rh.config.FederationQuery, name)
	record.StellarAddress = stellarAddress

	if rh.checkDBErr(err, w) {
		rh.writeResponse(w, record)
	}
}

// returns false if there was an error
func (rh *RequestHandler) checkDBErr(err error, w http.ResponseWriter) bool {
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Print("Federation record NOT found")
			rh.writeErrorResponse(w, ErrorResponseString("not_found", "Account not found"), http.StatusNotFound)
		} else {
			log.Print("Server error: ", err)
			rh.writeErrorResponse(w, ErrorResponseString("server_error", "Server error"), http.StatusInternalServerError)
		}
		return false
	}
	return true
}

func (rh *RequestHandler) writeResponse(w http.ResponseWriter, record FedRecord) {
	log.Print("Federation record found")

	json, err := json.Marshal(record)

	if err != nil {
		rh.writeErrorResponse(w, ErrorResponseString("server_error", "Server error"), http.StatusInternalServerError)
		return
	}

	w.Write(json)
}

func (rh *RequestHandler) writeErrorResponse(w http.ResponseWriter, response string, errorCode int) {
	w.WriteHeader(errorCode)
	w.Write([]byte(response))
}

func ErrorResponseString(code string, message string) string {
	error := Error{Code: code, Message: message}
	json, _ := json.Marshal(error)
	return string(json)
}
