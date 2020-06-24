package main

import (
	"fmt"
	_ "net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"github.com/swaggo/echo-swagger"
	"golang.org/x/crypto/acme/autocert"

	_ "nomsu-api/docs" // docs is generated by Swag CLI, but you have to import it.
	"nomsu-api/handlers"
)

// @title Schrodinger's Cup API
// @version 1.0
// @description A game of Peter Suber's Nomic running on the Tezos network.
// @BasePath /
func main() {

	// config
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("error reading config file: %v\n", err))
	}



	// set some defaults
	viper.SetDefault("PORT", "1323")
	viper.SetDefault("NOMSU", "nomsu")
	viper.SetDefault("JWT_SECRET", "secret")
	viper.SetDefault("PRODUCTION", "0")

	e := echo.New()

	// TLS
	isProduction, ok := viper.Get("PRODUCTION").(string)
	if !ok {
		panic(fmt.Sprintf("Error reading environment variable PRODUCTION: %v\n", ok))
	}
	if isProduction == "1" {
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS default
	// Allows requests from any origin wth GET, HEAD, PUT, POST or DELETE method.
	e.Use(middleware.CORS())

	// Routes
	e.GET("/", handlers.HelloWorld)
	e.GET("/ping", handlers.Ping)
	e.GET("/players", handlers.Players)
	e.GET("/round", handlers.Round)
	e.POST("/test", handlers.TestNomsu)
	e.POST("/auth", handlers.Auth)

	// Game requires auth
	g := e.Group("/game")
	g.Use(middleware.JWT([]byte(viper.GetString("JWT_SECRET"))))
	g.POST("/vote", handlers.CastVote)
	e.GET("/proposal/:round", handlers.Proposal)
	g.POST("/propose", handlers.SubmitProposal)

	// swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	if isProduction != "1" {
		e.Logger.Fatal(e.Start(":" + viper.GetString("PORT")))
	} else {
		e.Logger.Fatal(e.StartAutoTLS(":443"))	
	}	
}
