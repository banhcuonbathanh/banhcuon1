cd english-app-fe-nextjs

cd golang

go get -u github.com/go-chi/chi/v5

go run cmd/server/main.go
go run cmd/grcp-server/main.go
go run cmd/client/main.go
======================================= postgres ======================
psql -U myuser -d mydatabase

# psql -U myuser -d mydatabase

DROP DATABASE mydatabase;
TRUNCATE TABLE schema*migrations, users; delete all data
\dt : list all table
\d users
\d sessions
SELECT * FROM users;
SELECT \* FROM sessions;
DELETE FROM sessions;
\d order_items
mydatabase=# \d users
SELECT \* FROM sessions;
DROP TABLE sessions;
DELETE FROM schema_migrations;
=================================================== docker =======================
docker-compose up -d
docker-compose up
docker compose build go_app_ai
docker compose down
docker-compose up go_app_ai
//
========================================= golang ==============================

go run cmd/server/main.go

Run the desired commands using make <target>. For example:

To run the server: make run-server
To run the client: make run-client
To run all tests: make test
To run only the CreateUser test: make test-create
To run only the GetUser test: make test-getf
To clean build artifacts: make clean
To see available commands: make help

make stop-server

go test -v test/test-api/test-api.go
golang/
============================================== git hub ================================
git branch dev
git checkout golang-new-server-for-grpc

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ecomm-grpc/proto/user.proto

git checkout dev
git merge golang-new-server-for-grpc
git commit
git push origin dev
