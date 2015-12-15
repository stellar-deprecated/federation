package federation

import (
  "fmt"
  "net/http"
  "log"

  _ "github.com/go-sql-driver/mysql"
  "github.com/jmoiron/sqlx"
  _ "github.com/lib/pq"
  _ "github.com/mattn/go-sqlite3"
)

type App struct {
  config   Config
  database *sqlx.DB
}

// NewApp constructs an new App instance from the provided config.
func NewApp(config Config) (*App, error) {
  database, err := sqlx.Connect(
    config.DatabaseType,
    config.DatabaseUrl,
  )

  if err != nil {
    log.Panic(err)
  }

  result := &App{config: config, database: database}
  return result, nil
}

func (a *App) Serve() {
  requestHandler := &RequestHandler{app: a}

  http.HandleFunc("/federation", requestHandler.Main)
  http.HandleFunc("/federation/", requestHandler.Main)
  portString := fmt.Sprintf(":%d", a.config.Port)
  http.ListenAndServe(portString, nil)
}

func (a *App) Close() {
  a.database.Close();
}
