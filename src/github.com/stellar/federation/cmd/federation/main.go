package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stellar/federation"
)

var app *federation.App
var rootCmd *cobra.Command

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rootCmd.Execute()
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	rootCmd = &cobra.Command{
		Use:   "federation",
		Short: "stellar federation server",
		Long: `stellar federation server
=========================

Make sure config.toml file is in the working folder.
Required config values:
  - domain
  - database.type
  - database.url
  - queries.federation
  - queries.reverse-federation`,
		Run: run,
	}
}

func run(cmd *cobra.Command, args []string) {
	log.Print("Reading config.toml file")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if viper.GetString("database.type") == "" ||
		viper.GetString("database.url") == "" ||
		viper.GetString("domain") == "" ||
		viper.GetString("queries.federation") == "" ||
		viper.GetString("queries.reverse-federation") == "" {
		rootCmd.Help()
		os.Exit(1)
	}

	config := federation.Config{
		Port:                   viper.GetInt("port"),
		Domain:                 viper.GetString("domain"),
		DatabaseType:           viper.GetString("database.type"),
		DatabaseUrl:            viper.GetString("database.url"),
		FederationQuery:        viper.GetString("queries.federation"),
		ReverseFederationQuery: viper.GetString("queries.reverse-federation"),
	}

	app, err = federation.NewApp(config)

	if err != nil {
		log.Fatal(err.Error())
	}

	app.Serve()
}
