package federation

import (
	"flag"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"github.com/zenazn/goji"
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

	portString := fmt.Sprintf(":%d", a.config.Port)
	flag.Set("bind", portString)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET"},
	})
	goji.Use(c.Handler)
	goji.Use(headersMiddleware())
	goji.Use(stripTrailingSlashMiddleware())

	goji.Get("/federation", requestHandler.Main)
	goji.Serve()
}
