package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	api "doubleslash.de/coding-dojo-api/app/api"
	"doubleslash.de/coding-dojo-api/app/database"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

//go:embed api/swaggerui/*
var swaggerUIPage embed.FS

//go:embed api/coding-dojo-api.yaml
var openAPISpec embed.FS

func runApi(port int, serverImpl string, db *database.Database) {
	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec\n: %s", err)
	}
	swagger.Servers = nil

	var server api.ServerInterface
	switch serverImpl {
	case "mem":
		server = NewInMemoryServer()
	case "db":
		server = NewDatabaseServer(*db)
	default:
		server = NewInMemoryServer()
	}

	e := echo.New()
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	e.Use(echomiddleware.Logger())
	// e.Use(echomiddleware.AddTrailingSlash())
	// e.Use(middleware.OapiRequestValidator(swagger))

	specRoot, err := fs.Sub(openAPISpec, "api")
	if err != nil {
		log.Fatal(err)
	}
	e.GET("/coding-dojo-api.yaml", echo.WrapHandler(http.FileServer(http.FS(specRoot))))

	e.GET("/coding-dojo-api.json", func(ctx echo.Context) error {
		ctx.JSON(http.StatusOK, &swagger)
		return nil
	})

	serverRoot, err := fs.Sub(swaggerUIPage, "api/swaggerui")
	if err != nil {
		log.Fatal(err)
	}
	e.GET("/swaggerui/*", echo.WrapHandler(http.StripPrefix("/swaggerui/", http.FileServer(http.FS(serverRoot)))))

	api.RegisterHandlers(e, server)

	log.Print("REST API server listening on port ", port)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", port)))
}

func connectDB(hostname string, user string) (*database.Database, error) {
	connection, err := database.Connect(hostname, user)
	if err != nil {
		return nil, err
	}

	db := database.Database{DB: connection}

	err = db.InitTables()
	if err != nil {
		panic(err)
	}

	return &db, nil
}

func main() {
	var apiPort int
	var store string
	var postgresHostname, postgresUser string
	flag.IntVar(&apiPort, "api-port", 8008, "Port on which REST API will listen")
	flag.StringVar(&store, "store", "mem", "Server implementation (mem/db)")
	flag.StringVar(&postgresHostname, "postgres-hostname", "localhost", "Hostname for PostgreSQL")
	flag.StringVar(&postgresUser, "postgres-user", "postgres", "Username for PostgreSQL")
	flag.Parse()

	var db *database.Database
	if store == "db" {
		var err error
		db, err = connectDB(postgresHostname, postgresUser)
		if err != nil {
			panic(err)
		}
		defer db.Close()
	}

	runApi(apiPort, store, db)
}
