package datastores

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	todoerrors "go-to-do-app/to-do-lib/errors"
	"go-to-do-app/to-do-lib/logging"
	"go-to-do-app/to-do-lib/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type DataStore interface {
	AddItem(item models.ToDo) models.ToDo
	GetItem(userId string, itemId uuid.UUID) (models.ToDo, error)
	UpdateItem(item models.ToDo) (models.ToDo, error)
	Close()
}

type inMemDatastore struct {
	Items map[string]map[uuid.UUID]models.ToDo
	mut   sync.Mutex
}

// this function is only retuning (models.ToDo, error) to avoid code duplicatio on endpoints, which seems kinda bad
func (ds *inMemDatastore) AddItem(item models.ToDo) models.ToDo {
	item.Id = uuid.New()
	ds.mut.Lock()
	defer ds.mut.Unlock()

	if user, exists := ds.Items[item.UserId]; exists {
		user[item.Id] = item
	} else {
		ds.Items[item.UserId] = map[uuid.UUID]models.ToDo{item.Id: item}
	}
	return ds.Items[item.UserId][item.Id]
}

func (ds *inMemDatastore) GetItem(userId string, itemId uuid.UUID) (models.ToDo, error) {
	ds.mut.Lock()
	defer ds.mut.Unlock()
	if item, exists := ds.Items[userId][itemId]; exists {
		return item, nil
	}
	return models.ToDo{}, &todoerrors.NotFoundError{Message: "ToDo Not Found"}
}

func (ds *inMemDatastore) UpdateItem(item models.ToDo) (models.ToDo, error) {
	ds.mut.Lock()
	defer ds.mut.Unlock()

	if user, exists := ds.Items[item.UserId]; exists {
		if _, iexist := user[item.Id]; iexist {
			user[item.Id] = item
			return ds.Items[item.UserId][item.Id], nil
		}
	}
	return models.ToDo{}, &todoerrors.NotFoundError{Message: "ToDo Not Found"}
}

func (ds *inMemDatastore) Close() {
	//no action for in mem
}

func NewInMemDataStore() DataStore {
	return &inMemDatastore{Items: make(map[string]map[uuid.UUID]models.ToDo), mut: sync.Mutex{}}
}

func LoadJsonStore(fpath string) map[string]map[uuid.UUID]models.ToDo {
	file, err := os.Open(fpath)
	ctx := context.Background()
	logging.AddTraceID(ctx)
	if err != nil {
		logging.LogWithTrace(ctx, map[string]interface{}{}, err.Error())
	}
	defer file.Close()
	var todos []models.ToDo
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&todos)
	if err != nil {
		logging.LogWithTrace(ctx, map[string]interface{}{}, err.Error())
	}
	items := make(map[string]map[uuid.UUID]models.ToDo)
	for _, item := range todos {
		items[item.UserId] = map[uuid.UUID]models.ToDo{item.Id: item}
		if err != nil {
			logging.LogWithTrace(ctx, map[string]interface{}{}, fmt.Sprintf("error with todo: %+v", item))
		}
	}
	return items
}

type JsonDatastore struct {
	fpath string
	mut   sync.Mutex
	items map[string]map[uuid.UUID]models.ToDo
}

func (ds *JsonDatastore) AddItem(item models.ToDo) models.ToDo {
	item.Id = uuid.New()
	ds.mut.Lock()
	defer ds.mut.Unlock()
	if user, exists := ds.items[item.UserId]; exists {
		user[item.Id] = item
	} else {
		ds.items[item.UserId] = map[uuid.UUID]models.ToDo{item.Id: item}
	}
	ds.Close()
	return ds.items[item.UserId][item.Id]
}

func (ds *JsonDatastore) GetItem(userId string, itemId uuid.UUID) (models.ToDo, error) {
	ds.mut.Lock()
	defer ds.mut.Unlock()
	if item, exists := ds.items[userId][itemId]; exists {
		return item, nil
	}
	return models.ToDo{}, &todoerrors.NotFoundError{Message: "ToDo Not Found"}
}

func (ds *JsonDatastore) UpdateItem(item models.ToDo) (models.ToDo, error) {
	ds.mut.Lock()
	defer ds.mut.Unlock()
	if user, exists := ds.items[item.UserId]; exists {
		if _, iexist := user[item.Id]; iexist {
			ds.items[item.UserId][item.Id] = item
			ds.Close()
			return ds.items[item.UserId][item.Id], nil
		}
	}
	return models.ToDo{}, &todoerrors.NotFoundError{Message: "ToDo Not Found"}
}

func (ds *JsonDatastore) Close() {
	items := make([]models.ToDo, 0)
	for _, user := range ds.items {
		for _, item := range user {
			items = append(items, item)
		}
	}
	bytes, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	// Write the JSON to the file
	err = os.WriteFile(ds.fpath, bytes, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func NewJsonDatastore(path string) DataStore {
	return &JsonDatastore{fpath: path, items: LoadJsonStore(path), mut: sync.Mutex{}}
}

type PGDB struct {
	db      *sql.DB
	connStr string
	mut     sync.Mutex
}

func (p *PGDB) AddItem(item models.ToDo) models.ToDo {
	p.mut.Lock()
	defer p.mut.Unlock()
	id := uuid.New()
	p.db.Exec(
		"INSERT INTO items (user_id, item_id, title, priority, complete) VALUES($1, $2, $3, $4, $5)",
		item.UserId, id, item.Title, item.Priority, item.Complete,
	)
	rec, _ := p.GetItem(item.UserId, id)
	return rec
}
func (p *PGDB) GetItem(userId string, itemId uuid.UUID) (models.ToDo, error) {
	var item models.ToDo
	var (
		user_id  string
		item_id  string
		title    string
		priority string
		complete bool
	)
	if err := p.db.QueryRow(
		"SELECT * FROM items WHERE user_id = $1 AND item_id = $2",
		userId, itemId,
	).Scan(&user_id, &item_id, &title, &priority, &complete); err != nil {
		return models.ToDo{}, &todoerrors.NotFoundError{}
	}
	id, _ := uuid.Parse(item_id)
	item = models.ToDo{UserId: user_id, Id: id, Title: title, Priority: priority, Complete: complete}
	return item, nil
}
func (p *PGDB) UpdateItem(item models.ToDo) (models.ToDo, error) {
	p.mut.Lock()
	defer p.mut.Unlock()
	p.db.Exec(
		"UPDATE items SET user_id = $1, title = $3, priority = $4, complete = $5 WHERE user_id = $1 AND item_id = $2",
		item.UserId, item.Id, item.Title, item.Priority, item.Complete,
	)
	rec, err := p.GetItem(item.UserId, item.Id)
	if err != nil {
		return models.ToDo{}, err
	}
	return rec, nil
}
func (p *PGDB) Close() {
	p.db.Close()
}

func NewPGDatastore(user string, password string, database string) (DataStore, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", user, password, database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return &PGDB{db: nil, connStr: ""}, err
	}
	err = db.Ping()
	if err != nil {
		return &PGDB{db: nil, connStr: ""}, err
	}
	return &PGDB{db: db, connStr: connStr}, nil
}
