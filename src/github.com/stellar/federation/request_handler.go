package federation

import(
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"
)

type RequestHandler struct {
  app *App
}

func (rh *RequestHandler) Main(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  requestType := r.URL.Query().Get("type")
  switch requestType {
    case "name":
      rh.MakeDatabaseRequest(rh.app.config.FederationQuery, w, r);
    case "id":
      rh.MakeDatabaseRequest(rh.app.config.ReverseFederationQuery, w, r);
    default:
      fmt.Fprint(w, "invalid request")
  }
}

func (rh *RequestHandler) MakeDatabaseRequest(query string, w http.ResponseWriter, r *http.Request) {
  record := Record{}
  err := rh.app.database.Get(&record, query, r.URL.Query().Get("q"))

  if err != nil {
    if err.Error() == "sql: no rows in result set" {
      http.Error(w, ErrorResponseString("not_found", "Account not found"), http.StatusNotFound)
    } else {
      http.Error(w, ErrorResponseString("server_error", "Server error"), http.StatusInternalServerError)
    }
    return
  }

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
