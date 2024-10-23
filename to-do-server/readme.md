# TO DO Server Application

## Quickstart

Running the server application can be done from the to-do-server directory with `go run .` followed by the required flags that provide detail to the application about which datastore implementation it should utilise.


> `--mode=<in-mem|json-store|pgdb>` instructs the server the type of datastore to use.

> `--json=<path_to_.json>` specifies the *.json* store that a *json-store* datastore should load and save data to & from. As expected, this flag is not required with an *in-mem* datastore instance. 

> `--pg-create` can be used in combination with `--password=<db-password>` & `--user=db-username` to instruct the application to create the expected database & table required for the application. (currently defaulted to a localhost postgres) . Naturally the pre-requisite to using this command, or `--mode=pgdb` is to ensure that you have postgres installed in a local environment that is ready to be connected to. 

> *NOTE* Because credentials are required for testing the postgres implementation, a `.env` file should be added to the [datastores](../to-do-lib/datastores/) directory, following the `.env.example` file.

> A caveat to the above flags is that they are subject to change as development continues. A more universally appropriate flag structure may be applied when all datastore [Interfaces](../to-do-lib/datastores/datastores.go#L30)

## Implemented Datastores

- [x] In Mem
- [x] Json Store
- [x] Postgres DB

## API

Once the server is running, you can view the api spec for the relative versions with the below links:

- v1 <pr>The API spec can found at http://localhost:8081/v1/swagger-ui</pr>
- v2 <pr>The API spec can found at http://localhost:8081/v2/swagger-ui</pr>
