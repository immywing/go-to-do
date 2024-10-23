package datastores_test

import (
	"log"
	"math/rand/v2"
	"os"
	"sync"
	"testing"

	"go-to-do-app/to-do-lib/datastores"
	todoerrors "go-to-do-app/to-do-lib/errors"
	"go-to-do-app/to-do-lib/models"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestNewInMemDataStore(t *testing.T) {
	store := datastores.NewInMemDataStore()
	if store == nil {
		t.Error("Expected a non-nil DataStore, but got nil")
	}
}

func TestInMemAddToDo(t *testing.T) {
	store := datastores.NewInMemDataStore()
	expected := models.ToDo{Id: uuid.Max, Title: "test", Priority: "Low", Complete: false, UserId: "TestToDoUser"}
	actual := store.AddItem(expected)
	ok := actual.UserId == expected.UserId && actual.Title == expected.Title && actual.Priority == expected.Priority && actual.Complete == expected.Complete
	if !ok {
		t.Errorf("Expected: %+v, Got: %+v", expected, actual)
	}
}

func TestInMemUpdateToDo(t *testing.T) {
	store := datastores.NewInMemDataStore()
	item := models.ToDo{Id: uuid.Max, Title: "test", Priority: "Low", Complete: false, UserId: "TestToDoUser"}
	expected := store.AddItem(item)
	expected.Priority = "High"
	expected.Complete = true
	actual, _ := store.UpdateItem(expected)
	if actual != expected {
		t.Errorf("Expected: %+v, Got: %+v", expected, actual)
	}
}

func TestInMemUpdateNonExistientToDo(t *testing.T) {
	store := datastores.NewInMemDataStore()
	td := models.ToDo{Id: uuid.Max, Title: "test", Priority: "Low", Complete: false, UserId: "TestToDoUser"}
	expected := todoerrors.NotFoundError{Message: "ToDo Not Found"}
	_, actual := store.UpdateItem(td)
	_, ok := actual.(*todoerrors.NotFoundError)
	if !ok {
		t.Errorf("Expected: %T, Got: %T", expected, actual)
	}
}

func TestInMemGetToDo(t *testing.T) {
	store := datastores.NewInMemDataStore()
	item := models.ToDo{Id: uuid.Max, Title: "test", Priority: "Low", Complete: false, UserId: "TestToDoUser"}
	expected := store.AddItem(item)
	actual, err := store.GetItem(expected.UserId, expected.Id)
	if err != nil {
		t.Errorf("datastore unable to find item that was created with uuid: %s", expected.Id)
	}
	if actual != expected {
		t.Errorf("Expected: %+v, Got: %+v", expected, actual)
	}
}

func TestJSONMemDataStore(t *testing.T) {
	store := datastores.NewInMemDataStore()
	if store == nil {
		t.Error("Expected a non-nil DataStore, but got nil")
	}
}

func TestJSONAddToDo(t *testing.T) {
	store := datastores.NewJsonDatastore("store.json")
	expected := models.ToDo{Id: uuid.Max, Title: "test", Priority: "Low", Complete: false, UserId: "TestToDoUser"}
	actual := store.AddItem(expected)
	ok := actual.UserId == expected.UserId && actual.Title == expected.Title && actual.Priority == expected.Priority && actual.Complete == expected.Complete
	if !ok {
		t.Errorf("Expected: %+v, Got: %+v", expected, actual)
	}
}

func TestJSONUpdateToDo(t *testing.T) {
	store := datastores.NewInMemDataStore()
	item := models.ToDo{Id: uuid.Max, Title: "test", Priority: "Low", Complete: false, UserId: "TestToDoUser"}
	expected := store.AddItem(item)
	expected.Priority = "High"
	expected.Complete = true
	actual, _ := store.UpdateItem(expected)
	if actual != expected {
		t.Errorf("Expected: %+v, Got: %+v", expected, actual)
	}
}

func TestJSONUpdateNonExistientToDo(t *testing.T) {
	store := datastores.NewJsonDatastore("store.Json")
	td := models.ToDo{Id: uuid.Max, Title: "test", Priority: "Low", Complete: false, UserId: "TestToDoUser"}
	expected := todoerrors.NotFoundError{Message: "ToDo Not Found"}
	_, actual := store.UpdateItem(td)
	_, ok := actual.(*todoerrors.NotFoundError)
	if !ok {
		t.Errorf("Expected: %T, Got: %T", expected, actual)
	}
}

func TestJSONGetToDo(t *testing.T) {
	store := datastores.NewJsonDatastore("store.Json")
	item := models.ToDo{Id: uuid.Max, Title: "test", Priority: "Low", Complete: false, UserId: "TestToDoUser"}
	expected := store.AddItem(item)
	actual, err := store.GetItem(expected.UserId, expected.Id)
	if err != nil {
		t.Errorf("datastore unable to find item that was created with uuid: %s", expected.Id)
	}
	if actual != expected {
		t.Errorf("Expected: %+v, Got: %+v", expected, actual)
	}
}

func TestPostgresAddToDo(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	password := os.Getenv("DB_PASSWORD")
	store, err := datastores.NewPGDatastore("postgres", password, "todo")
	if err != nil {
		log.Fatal("unable to connect to PG databased with credentials")
	}
	expected := models.ToDo{Id: uuid.Max, Title: "test", Priority: "Low", Complete: false, UserId: "TestToDoUser"}
	actual := store.AddItem(expected)
	ok := actual.UserId == expected.UserId && actual.Title == expected.Title && actual.Priority == expected.Priority && actual.Complete == expected.Complete
	if !ok {
		t.Errorf("Expected: %+v, Got: %+v", expected, actual)
	}
}

func TestPostgresGetToDo(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	password := os.Getenv("DB_PASSWORD")
	store, err := datastores.NewPGDatastore("postgres", password, "todo")
	if err != nil {
		log.Fatal("unable to connect to PG databased with credentials")
	}
	item := models.ToDo{Id: uuid.Max, Title: "test", Priority: "Low", Complete: false, UserId: "TestToDoUser"}
	expected := store.AddItem(item)
	actual, err := store.GetItem(expected.UserId, expected.Id)
	if err != nil {
		t.Errorf("datastore unable to find item that was created with uuid: %s", expected.Id)
	}
	if actual != expected {
		t.Errorf("Expected: %+v, Got: %+v", expected, actual)
	}
}

func TestPostgresUpdateToDo(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	password := os.Getenv("DB_PASSWORD")
	store, err := datastores.NewPGDatastore("postgres", password, "todo")
	if err != nil {
		log.Fatal("unable to connect to PG databased with credentials")
	}
	item := models.ToDo{Id: uuid.Max, Title: "test", Priority: "Low", Complete: false, UserId: "TestToDoUser"}
	expected := store.AddItem(item)
	expected.Priority = "High"
	expected.Complete = true
	actual, _ := store.UpdateItem(expected)
	if actual != expected {
		t.Errorf("Expected: %+v, Got: %+v", expected, actual)
	}
}

func TestPostgresUpdateNonExistientToDo(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	password := os.Getenv("DB_PASSWORD")
	store, err := datastores.NewPGDatastore("postgres", password, "todo")
	if err != nil {
		log.Fatal("unable to connect to PG databased with credentials")
	}
	td := models.ToDo{Id: uuid.Max, Title: "test", Priority: "Low", Complete: false, UserId: "TestToDoUser"}
	expected := todoerrors.NotFoundError{Message: "ToDo Not Found"}
	_, actual := store.UpdateItem(td)
	_, ok := actual.(*todoerrors.NotFoundError)
	if !ok {
		t.Errorf("Expected: %T, Got: %T", expected, actual)
	}
}

func TestConcurrentPutRequests(t *testing.T) {
	stores := []datastores.DataStore{
		datastores.NewInMemDataStore(),
		datastores.NewJsonDatastore("store.json"),
	}
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	password := os.Getenv("DB_PASSWORD")
	pg, err := datastores.NewPGDatastore("postgres", password, "todo")
	if err != nil {
		log.Fatal("unable to connect to PG databased with credentials")
	}
	stores = append(stores, pg)
	statuses := []bool{true, false}
	priorities := []string{models.PriorityLow, models.PriorityMedium, models.PriorityHigh}
	itmev1 := models.ToDo{Id: uuid.Max, Title: "test", Priority: "High", Complete: false, UserId: ""}
	itemv2 := models.ToDo{Id: uuid.Max, Title: "test", Priority: "High", Complete: false, UserId: "TestToDoUser"}
	versions := make(map[string]models.ToDo)
	for _, datastore := range stores {
		expectedV1 := datastore.AddItem(itmev1)
		expectedV2 := datastore.AddItem(itemv2)
		versions["v1"] = expectedV1
		versions["v2"] = expectedV2

		var wg sync.WaitGroup
		numRequests := 1000
		for _, expectedItem := range versions {
			for i := 1; i < numRequests; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					expected := expectedItem
					expected.Priority = priorities[rand.IntN(len(statuses))]
					expected.Complete = statuses[rand.IntN(len(statuses))]
					actual, _ := datastore.UpdateItem(expected)
					if actual != expected {
						t.Errorf("Expected %+v, Got %+v", expected, actual)
					}
				}(i)
			}
		}
		wg.Wait()
		datastore.Close()
	}
}
