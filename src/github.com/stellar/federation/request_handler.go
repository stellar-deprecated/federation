package federation

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type RequestHandler struct {
	app *App
}

func (rh *RequestHandler) Main(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	requestType := r.URL.Query().Get("type")
	q := r.URL.Query().Get("q")
	switch {
	case requestType == "name" && q != "":
		rh.MakeDatabaseRequest(rh.app.config.FederationQuery, w, r)
	case requestType == "id" && q != "":
		rh.MakeDatabaseRequest(rh.app.config.ReverseFederationQuery, w, r)
	default:
		http.Error(w, ErrorResponseString("invalid_request", "Invalid request"), http.StatusBadRequest)
	}
}

func (rh *RequestHandler) MakeDatabaseRequest(query string, w http.ResponseWriter, r *http.Request) {
	record := Record{}

	q := r.URL.Query().Get("q")

	if r.URL.Query().Get("type") == "name" {
		name := ""
		domain := ""

		if i := strings.Index(q, "*"); i >= 0 {
			name = q[:i]
			domain = q[i+1:]
		}

		if name == "" || domain != rh.app.config.Domain {
			http.Error(w, ErrorResponseString("not_found", "Incorrect Domain"), http.StatusNotFound)
			return
		}

		q = name
	}

	err := rh.app.database.Get(&record, query, q)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Print("Federation record NOT found")
			http.Error(w, ErrorResponseString("not_found", "Account not found"), http.StatusNotFound)
		} else {
			log.Print("Server error: ", err)
			http.Error(w, ErrorResponseString("server_error", "Server error"), http.StatusInternalServerError)
		}
		return
	}

	log.Print("Federation record found")

	var usernameBuffer bytes.Buffer
	usernameBuffer.WriteString(record.Username)
	usernameBuffer.WriteString("*")
	usernameBuffer.WriteString(rh.app.config.Domain)

	record.Username = usernameBuffer.String()
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
