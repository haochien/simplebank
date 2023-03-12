# external lib required by this project:
1. operate db migration (for windows system):
    ```bash
    $ scoop install migrate
    ```

2. lib/pq: to let go sql to talk to specific db engine
    ```bash
    go get github.com/lib/pq
    ```

3. testify: testing framework
    ```bash
    go get github.com/stretchr/testify
    ```

4. Golang Gin (Web framework)
    ```bash
    go get -u github.com/gin-gonic/gin
    ```

5. Golang Viper (env config )
    ```bash
    go get github.com/spf13/viper
    ```

5. Golang GoMock (mock DB)
    ```bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```



# Set up and connect postgres from docker:
0. pull docker image:
    ```docker
    docker pull postgres:15-alpine
    ```


1. connect to postgres:
    ```docker
    docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine
    ```
2. create/drop db
    ```docker
    docker exec -it postgres15 createdb --username=root --owner=root simple_bank

    docker exec -it postgres15 dropdb simple_bank
    ```

3. others:
    ```docker
    # stop db container
    docker stop postgres15

    # start db container
    docker start postgres15

    # remove current db container
    docker rm postgres15

    # log db 
    docker log postgres15

    # execute the postgres image as user root
    docker exec -it postgres15 psql -U root
    select now()
    ```


# Prepare migration:
1. init migration folder:
    ```bash
    migrate create -ext sql -dir db/migration -seq init_schema
    ```

2. run migration script:
    ```bash
    migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

    migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
    ```

3. we can continue on create other migration scripts:
    ```bash
    migrate create -ext sql -dir db/migration -seq second_script
    ```


# How to set up go module and dependency:
1. init module:
    ```bash
    go mod init github.com/haochien/simplebank
    ```
2. create/update dependency:
    ```bash
    go mod tidy
    ```


# How to use sqlc in go project:
0. download sqlc:
    ```docker
    docker pull kjconroy/sqlc
    ```
1. init sqlc: create sqlc.yaml. Can be done via:
    ```docker
    docker run --rm -v ${pwd}:/src -w /src kjconroy/sqlc init
    ```
2. configure content in sqlc.yaml
3. prepare sql files under query folder
4. generate go files sqlc functions via:
    ```docker
    docker run --rm -v ${pwd}:/src -w /src kjconroy/sqlc generate
    ```

5. create transaction maintain file: store.go 


# How to set up go Mock DB environment:
0. download gomock:
    ```bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```
 

1. all function under store struct should be map to store interface.
   we can manually add all those function to new store interface created in the next step.
   But easier way would be add following parameter in sqlc.yaml:
   ```yaml
   emit_interface: true
   ```
   and regenerate sqlc
    ```bash
    docker run --rm -v ${pwd}:/src -w /src kjconroy/sqlc generate

    # querier.go will then be created 
    ```

    

2. create a store interface to differentiate the real store struct
    ```go
    // in store.go file: 

    // Store defines all functions to execute db queries and transactions
    type Store interface {
        Querier
        TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
    }

    // SQLStore provides all functions to execute SQL queries and transactions
    type SQLStore struct {
        db *sql.DB
        *Queries
    }

    ```

3. create a mock folder under db folder, and run following command:
    ```bash
    mockgen -package mockdb -destination db/mock/store.go github.com/haochien/simplebank/db/sqlc Store

    # store.go will then created under mock folder
    ```