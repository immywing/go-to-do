package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"go-to-do-app/to-do-lib/datastores"
	"go-to-do-app/to-do-lib/logging"
	"go-to-do-app/to-do-server/server"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	_ "github.com/lib/pq"
)

var (
	mode         = flag.String("mode", "", "set the mode the application should run in (in-mem, json-store, pgdb)")
	addr         = flag.String("address", ":8081", "set the address for the server. Default is :8081")
	jsonPath     = flag.String("json", "", "filepath of json file to use as datastore")
	password     = flag.String("password", "", "database password")
	user         = flag.String("user", "postgres", "database username")
	create       = flag.Bool("pg-create", false, "Create ToDo database & items table with postgres connection")
	shutdownChan = make(chan bool)
	dbname       = "todo"
)

func createPostgresDB() {
	connStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", *user, *password, "")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	_, err = db.Exec("CREATE DATABASE todo")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	todoConnStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", *user, *password, "todo")
	tododb, err := sql.Open("postgres", todoConnStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer tododb.Close()
	tododb.Exec("CREATE TABLE IF NOT EXISTS items (user_id TEXT, item_id TEXT, title TEXT, priority TEXT, complete BOOLEAN);")
	os.Exit(0)
}

func run() {
	flag.Parse()

	var store datastores.DataStore
	if *create {
		createPostgresDB()
	}
	if *mode == "" {
		logging.LogWithTrace(
			context.Background(),
			map[string]interface{}{},
			"no valid mode provided to start server with datastore",
		)
		os.Exit(1)
	}
	if *mode == "pgdb" {
		store, _ = datastores.NewPGDatastore(*user, *password, dbname)
		defer store.Close()
	}
	if *mode == "in-mem" {
		store = datastores.NewInMemDataStore()
	}
	if *mode == "json-store" {
		if filepath.Ext(*jsonPath) != ".json" {
			logging.LogWithTrace(
				context.Background(),
				map[string]interface{}{"path": *jsonPath},
				"no valid path to json file provided",
			)
			os.Exit(1)
		}
		store = datastores.NewJsonDatastore(*jsonPath)
		defer store.Close()
	}
	srv := server.NewToDoServer(*addr, shutdownChan, store)
	go srv.Start()
	fmt.Println("server running @", *addr)
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	<-interruptChannel
	srv.Shutdown()
	srv.AwaitShutdown()
}

func main() {
	run()
}
