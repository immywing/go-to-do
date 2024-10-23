package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"go-to-do-app/to-do-lib/apiclient"
	"go-to-do-app/to-do-lib/logging"
	"go-to-do-app/to-do-lib/models"

	"github.com/google/uuid"
)

var (
	post     = flag.Bool("post", false, "Add new Todo")
	put      = flag.Bool("put", false, "updateTodo")
	get      = flag.Bool("get", false, "Get existing Todo")
	id       = flag.String("id", "", "UUID of ToDo item")
	userId   = flag.String("user-id", "", "UUID representing user id")
	title    = flag.String("title", "", "Title of ToDo item")
	priority = flag.String("priority", "", "Priority of ToDo item")
	complete = flag.Bool("complete", false, "Completion status of ToDo item")
	version  = flag.String("version", "", "version of the api to use")
)

func cli() {
	flag.Parse()
	todoflags := map[string]string{
		"user-id":  *userId,
		"id":       *id,
		"title":    *title,
		"priority": *priority,
		"complete": strconv.FormatBool(*complete),
		"version":  *version,
	}
	var item models.ToDo
	var err error
	ctx := logging.AddTraceID(context.Background())
	client := apiclient.NewAPIClient("http://localhost:8081/")
	if !*post && !*put && !*get {
		err = errors.New("no method flag provided. requires 1 of --<post|put|get>")
		fmt.Println(err)
		os.Exit(1)
	}
	if *version != "v1" && *version != "v2" {
		err = errors.New("missing required flag of --version=<v1|v2>")
		fmt.Println(err)
		os.Exit(1)
	}
	if *post {
		item = models.ToDo{
			Id:       uuid.Max,
			UserId:   *userId,
			Title:    *title,
			Priority: *priority,
			Complete: *complete,
		}
		item, err = client.Req(ctx, "POST", item, todoflags)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("POST success! API response:\n", item)
	}
	if *put {
		item, err = models.NewToDo(userId, id, title, priority, complete)
		if err != nil {
			err = fmt.Errorf("failed to create to do model with flags of %s. ERR: %s", todoflags, err)
			fmt.Println(err)
			os.Exit(1)
		}
		item, err = client.Req(ctx, "PUT", item, todoflags)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("PUT success! API response:\n", item)
	}
	if *get {
		item = models.ToDo{
			Id:       uuid.Max,
			UserId:   *userId,
			Title:    *title,
			Priority: *priority,
			Complete: *complete,
		}
		item, err = client.Req(ctx, "GET", item, todoflags)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("GET success! API response:\n", item)
	}
}

func main() {
	cli()
}
