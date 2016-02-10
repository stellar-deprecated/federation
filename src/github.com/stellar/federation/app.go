package federation

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/justinas/alice"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

type Database interface {
	Get(dest interface{}, query string, args ...interface{}) error
}

type App struct {
	config   Config
	database Database
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
	requestHandler := &RequestHandler{config: &a.config, database: a.database}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET"},
	})

	handler := alice.New(
		c.Handler,
		headersMiddleware(),
		stripTrailingSlashMiddleware(),
	).Then(http.HandlerFunc(requestHandler.Main))

	http.Handle("/federation", handler)
	http.Handle("/federation/", handler)

	portString := fmt.Sprintf(":%d", a.config.Port)
	err := http.ListenAndServe(portString, nil)
	if err != nil {
		log.Fatal(err)
	}
}
