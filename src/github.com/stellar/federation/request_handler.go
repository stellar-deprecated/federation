package federation

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/stellar/federation/config"
	"github.com/stellar/federation/db"
)

type RequestHandler struct {
	config *config.Config
	driver db.Driver
}

func (rh *RequestHandler) Main(w http.ResponseWriter, r *http.Request) {
	requestType := r.URL.Query().Get("type")
	q := r.URL.Query().Get("q")

	switch {
	case requestType == "name" && q != "":
		rh.FederationRequest(q, w)
	case requestType == "id" && q != "":
		rh.ReverseFederationRequest(q, w)
	case requestType == "txid" && q != "":
		rh.writeErrorResponse(w, ErrorResponseString("not_implemented", "txid requests are not supported"), http.StatusNotImplemented)
	default:
		rh.writeErrorResponse(w, ErrorResponseString("invalid_request", "Invalid request"), http.StatusBadRequest)
	}
}

func (rh *RequestHandler) ReverseFederationRequest(accountID string, w http.ResponseWriter) {
	response := Response{
		AccountId: accountID,
	}

	record, err := rh.driver.GetByAccountId(accountID, rh.config.Queries.ReverseFederation)

	if err != nil {
		log.Print("Server error: ", err)
		rh.writeErrorResponse(w, ErrorResponseString("server_error", "Server error"), http.StatusInternalServerError)
	} else if record == nil {
		log.Print("Federation record NOT found")
		rh.writeErrorResponse(w, ErrorResponseString("not_found", "Account not found"), http.StatusNotFound)
	} else {
		response.StellarAddress = record.Name + "*" + rh.config.Domain
		rh.writeResponse(w, response)
	}
}

func (rh *RequestHandler) FederationRequest(stellarAddress string, w http.ResponseWriter) {
	var name, domain string

	if i := strings.Index(stellarAddress, "*"); i >= 0 {
		name = stellarAddress[:i]
		domain = stellarAddress[i+1:]
	}

	if name == "" || domain != rh.config.Domain {
		rh.writeErrorResponse(w, ErrorResponseString("not_found", "Incorrect Domain"), http.StatusNotFound)
		return
	}

	response := Response{
		StellarAddress: stellarAddress,
	}

	record, err := rh.driver.GetByStellarAddress(name, rh.config.Queries.Federation)

	if err != nil {
		log.Print("Server error: ", err)
		rh.writeErrorResponse(w, ErrorResponseString("server_error", "Server error"), http.StatusInternalServerError)
	} else if record == nil {
		log.Print("Federation record NOT found")
		rh.writeErrorResponse(w, ErrorResponseString("not_found", "Account not found"), http.StatusNotFound)
	} else {
		response.AccountId = record.AccountId
		response.MemoType = record.MemoType
		response.Memo = record.Memo
		rh.writeResponse(w, response)
	}
}

func (rh *RequestHandler) writeResponse(w http.ResponseWriter, response Response) {
	log.Print("Federation record found")

	json, err := json.Marshal(response)

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
