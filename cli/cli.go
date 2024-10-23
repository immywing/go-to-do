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
)

var (
	post       = flag.Bool("post", false, "Add new Todo")
	put        = flag.Bool("put", false, "updateTodo")
	get        = flag.Bool("get", false, "Get existing Todo")
	id         = flag.String("id", "", "UUID of ToDo item")
	userId     = flag.String("user-id", "", "UUID representing user id")
	title      = flag.String("title", "", "Title of ToDo item")
	priority   = flag.String("priority", "", "Priority of ToDo item")
	complete   = flag.Bool("complete", false, "Completion status of ToDo item")
	version    = flag.String("version", "", "version of the api to use")
	cliactions = []CliAction{
		{flag: post, do: cliPost},
		{flag: put, do: cliPut},
		{flag: get, do: cliGet},
	}
)

type CliAction struct {
	flag *bool
	do   func(todoflags map[string]string, client apiclient.APIClient, ctx context.Context)
}

func cliPost(todoflags map[string]string, client apiclient.APIClient, ctx context.Context) {
	item, err := client.Req(ctx, "POST", todoflags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("POST success! API response:\n", item)
}

func cliPut(todoflags map[string]string, client apiclient.APIClient, ctx context.Context) {
	item, err := client.Req(ctx, "PUT", todoflags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("PUT success! API response:\n", item)
}

func cliGet(todoflags map[string]string, client apiclient.APIClient, ctx context.Context) {
	item, err := client.Req(ctx, "GET", todoflags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("GET success! API response:\n", item)
}

func cli() {

	var err error

	flag.Parse()
	todoflags := map[string]string{
		"user-id":  *userId,
		"id":       *id,
		"title":    *title,
		"priority": *priority,
		"complete": strconv.FormatBool(*complete),
		"version":  *version,
	}
	ctx := logging.AddTraceID(context.Background())
	client := apiclient.NewAPIClient("http://localhost:8081/")

	if *version != "v1" && *version != "v2" {
		err = errors.New("missing required flag of --version=<v1|v2>")
		fmt.Println(err)
		os.Exit(1)
	}

	for _, action := range cliactions {
		if *action.flag {
			action.do(todoflags, client, ctx)
			os.Exit(0)
		}
	}

	err = errors.New("no method flag provided. requires 1 of --<post|put|get>")
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	cli()
}
