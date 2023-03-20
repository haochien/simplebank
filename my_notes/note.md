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

6. uuid for Golang:
    ```bash
    go get github.com/google/uuid
    ```


7. Golang JWT:
    ```bash
    go get github.com/dgrijalva/jwt-go
    ```

8. Golang Paseto:
    ```bash
    go get github.com/o1egl/paseto
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


# How to build docker file:
1. create dockerfile

2. build the image
    ```docker
    docker build -t simplebank:latest .

    # in case rerun and need to delete old image (c8f6cb16c708 is image id):
    docker rmi c8f6cb16c708 
    ```

3. run created image
    ```docker
    # --name: specify the name of container we are going to create
    docker run --name simplebank -p 8080:8080 simplebank:latest

    # to run not in debug mode:
    docker run --name simplebank -p 8080:8080 -e GIN_MODE=release simplebank:latest

    # to deal with issue that default ip address of postgres and our simplebank is different:
    docker run --name simplebank -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@172.17.0.2:5432/simple_bank?sslmode=disable" simplebank:latest
    # note: you can use following command to check the ip address
    docker container inspect simplebank
    docker container inspect postgres15  

    # better solution to deal with different ip in two container: create your own network, instead of default bridge network
    docker network create bank-network
    docker network connect bank-network postgres15
    docker network inspect bank-network
    docker run --name simplebank --network bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres15:5432/simple_bank?sslmode=disable" simplebank:latest


    # in case error occurs and need to rebuild:
    docker rm simplebank
    docker rmi c8f6cb16c708 
    docker build -t simplebank:latest .
    ```

4. create start.sh to able to run migration up command:

    remember to update migrate commands in dockerfile to enable migrate binary download and usage.

    also need update dockerfile to copy start.sh into image
    
    add ENTRYPOINT [ "/app/start.sh" ] in dockerfile to make start.sh as start point of the command

5. create docker-compose.yaml file to auto run step 3 commands:
    ```docker
    # to run docker compose file:
    # after run this, it will start from build image, and then connect to db container and access to api container in parallel
    docker compose up 

    # if need to revert, then run:
    docker compose down
    docker rmi simplebank-api
    ```

    Note that error might occur since db container and api container are build in parallel.

    The api should wait until db is set up (because api relies on db connection)

    Thus, depends_on should be setup in docker-compose.yaml.
    
    But depends_on does not wait for db and redis to be “ready” before starting web, healthcheck should be added to provide wait feature

