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
  viper.SetConfigType("yaml")
  viper.AddConfigPath(".")

  rootCmd = &cobra.Command{
    Use:   "federation",
    Short: "stellar federation server",
    Long:
`stellar federation server
=========================

Make sure config.yaml file is in the current folder.
Required config values:
  - database-type
  - database-url
  - domain
  - federation-query
  - reverse-federation-query`,
    Run:   run,
  }
}

func run(cmd *cobra.Command, args []string) {
  err := viper.ReadInConfig()
  if err != nil {
    panic(fmt.Errorf("Fatal error config file: %s \n", err))
  }

  if viper.GetString("database-type") == "" ||
     viper.GetString("database-url") == "" ||
     viper.GetString("domain") == "" ||
     viper.GetString("federation-query") == "" ||
     viper.GetString("reverse-federation-query") == "" {
    rootCmd.Help()
    os.Exit(1)
  }

  config := federation.Config{
    Port:                   viper.GetInt("port"),
    Domain:                 viper.GetString("domain"),
    DatabaseType:           viper.GetString("database-type"),
    DatabaseUrl:            viper.GetString("database-url"),
    FederationQuery:        viper.GetString("federation-query"),
    ReverseFederationQuery: viper.GetString("reverse-federation-query"),
  }

  app, err = federation.NewApp(config)

  if err != nil {
    log.Fatal(err.Error())
  }

  app.Serve()
}
