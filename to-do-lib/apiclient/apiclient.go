package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go-to-do-app/to-do-lib/models"

	"github.com/google/uuid"
)

type APIClient struct {
	BaseURL    string
	httpClient *http.Client
}

func (c *APIClient) Req(
	ctx context.Context, m string, args map[string]string) (models.ToDo, error) {

	var apiURL string
	var req *http.Request
	var itemIn models.ToDo
	var buffer []byte
	var err error

	userid := args["user-id"]
	itemid := args["id"]
	version := args["version"]
	title := args["title"]
	priority := args["priority"]
	complete := args["complete"] == "true"

	if m == http.MethodGet {
		apiURL = fmt.Sprintf("http://localhost:8081/%s/todo?user_id=%s&id=%s",
			version, userid, itemid)
	}
	if m == http.MethodPut {
		apiURL = fmt.Sprintf("http://localhost:8081/%s/todo", args["version"])
		itemIn, err = models.NewToDo(&userid, &itemid, &title, &priority, &complete)
		if err != nil {
			return models.ToDo{}, err
		}
		buffer, err = json.Marshal(itemIn)
		if err != nil {
			return models.ToDo{}, err
		}
	}
	if m == http.MethodPost {
		apiURL = fmt.Sprintf("http://localhost:8081/%s/todo", args["version"])
		itemIn = models.ToDo{Id: uuid.Max, UserId: userid, Title: title, Priority: priority, Complete: complete}
		buffer, err = json.Marshal(itemIn)
		if err != nil {
			return models.ToDo{}, err
		}
	}
	req, err = http.NewRequest(m, apiURL, bytes.NewBuffer(buffer))
	var item models.ToDo
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.ToDo{}, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return models.ToDo{}, err
	}
	return item, nil
}

func NewAPIClient(baseURL string) APIClient {
	return APIClient{BaseURL: baseURL, httpClient: &http.Client{}}
}
