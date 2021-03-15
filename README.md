# :memo: ToDo

A simple ToDo REST API.

## Getting started

### Requirements

* Go 1.15 or Docker
* MariaDB or Docker

### Run the application with Docker Compose

If you just want to get the application up and running, use Docker Compose:

```
$ docker-compose up
```

The REST API is exposed on `localost:8000`. Note that the service configuration
is suited for demonstration purposes only.

### Compile and run the application with Go

First, get a MariaDB server with administrator permissions up and running. If
you don't have MariaDB installed locally, just launch a Docker container:

```
$ docker container run -p 3306:3306 -e MYSQL_ROOT_PASSWORD=test123 mariadb
```

Once MariaDB is running, use `go run` to compile and run the application. Make
sure to pass the correct MariaDB credentials.

```
$ go run . --mariadb-user root --mariadb-password test123
```

The REST API is exposed on `localost:8000`.

### Configuration

The example above uses `test123` as MariaDB root password. However, the access
to MariaDB and other parameters can be configured exactly to your needs. These
are the configuration values available:

|Configuration Value|Default|Environment Variable|CLI Flag|
|-|-|-|-|
|MariaDB user|`admin`|`TODO_MARIADB_USER`|`--mariadb-user`|
|MariaDB password|`admin`|`TODO_MARIADB_PASSWORD`|`--mariadb-password`|
|MariaDB address|`0.0.0.0:3306`|`TODO_MARIADB_ADDRESS`|`--mariadb-address`|
|MariaDB DB name|`todo_app`|`TODO_MARIADB_DBNAME`|`--mariadb-dbname`|
|ToDo API port|`8000`|`TODO_PORT`|`--port`|

## REST API

For a detailed overview, see the [OpenAPI definition](swagger.yaml).

### Models

The API expects and returns ToDo items looking as follows:

```json
{
  "id": 1,
  "name": "My ToDo",
  "description": "My ToDo Description",
  "tasks": [
    {
      "id": 1,
      "name": "A Task",
      "description": "A Task Description"
    }
  ]
}
```

The `ID` fields have to be empty when the respective item doesn't exist yet,
e.g. when calling `POST /todos`.

### Endpoints

|Method|Route|Description|Expected Body|
|POST|`/todos`|Creates a new ToDo|A ToDo item without ID|
|GET|`/todos`|Returns a list of all ToDos|-|
|GET|`/todos/{id}`|Returns a ToDo|-|
|PUT|`/todos/{id}`|Overwrites an existing Todo|An updated ToDo item|
|DELETE|`/todos/{id}`|Deletes a ToDo|-|
