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
\d reading_test_models;
SELECT * FROM reading_test_models;
SELECT \* FROM sessions;
DELETE FROM sessions;
\d order_items
mydatabase=# \d users
SELECT \* FROM reading_tests;
DROP TABLE schema_migrations;
DELETE FROM schema_migrations;
DELETE FROM reading_tests;
\l
\c testdb
testdb=# \dT+ paragraph_content
UPDATE users
SET is_admin = true
WHERE id = 1;
migrate -database postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable force 7

-- List all tables in the public schema
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
AND table_type = 'BASE TABLE';

-- List all custom types (including ENUMs)
SELECT t.typname AS enum_name,
e.enumlabel AS enum_value
FROM pg_type t
JOIN pg_enum e ON t.oid = e.enumtypid
JOIN pg_catalog.pg_namespace n ON n.oid = t.typnamespace
WHERE n.nspname = 'public';

-- Drop all tables in the public schema
DO $$
DECLARE
r RECORD;
BEGIN
FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
END LOOP;
END $$;

-- Drop the question_type ENUM
DROP TYPE IF EXISTS question_type CASCADE;

-- Verify that all tables are dropped
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
AND table_type = 'BASE TABLE';

-- Verify that the question_type ENUM is dropped
SELECT t.typname AS enum_name,
e.enumlabel AS enum_value
FROM pg_type t
JOIN pg_enum e ON t.oid = e.enumtypid
JOIN pg_catalog.pg_namespace n ON n.oid = t.typnamespace
WHERE n.nspname = 'public';
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

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ecomm-grpc/proto/comment/comment.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ecomm-grpc/proto/user.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ecomm-grpc/proto/reading/reading.proto
git checkout nextjs-fe-readiding-add-more-clean-architextture
git merge golang-new-server-for-grpc
git commit
git push origin dev

golang/ecomm-grpc/proto/reading/reading.proto

git checkout -b golang: create new branch

reading_test_models
section_models
passage_models
schema_migrations
paragraph_content_models
question_models
users
sessions
