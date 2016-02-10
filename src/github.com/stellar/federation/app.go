package federation

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/alice"
	_ "github.com/lib/pq"
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
		config.Database.Type,
		config.Database.Url,
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
	).Then(http.HandlerFunc(requestHandler.Main))

	http.Handle("/federation", handler)
	http.Handle("/federation/", handler)

	portString := fmt.Sprintf(":%d", a.config.Port)

	log.Println("Starting server on port: ", a.config.Port)

	var err error
	if a.config.Tls.CertificateFile != "" && a.config.Tls.PrivateKeyFile != "" {
		err = http.ListenAndServeTLS(
			portString,
			a.config.Tls.CertificateFile,
			a.config.Tls.PrivateKeyFile,
			nil,
		)
	} else {
		err = http.ListenAndServe(portString, nil)
	}

	if err != nil {
		log.Fatal(err)
	}
}
