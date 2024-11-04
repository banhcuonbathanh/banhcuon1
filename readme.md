http://localhost:3000/table/1?token=MTp0YWJsZTo0ODgzNzk4NTY1.mI7st71i-AQ
fmt.Printf("golang/quanqr/order/order_handler.go ordersResponse %v\n", ordersResponse)
docker compose up

cd quananqr1
npm run dev

cd english-app-fe-nextjs

cd golang

go get -u github.com/go-chi/chi/v5
cd golang
go run cmd/server/main.go
cd golang
go run cmd/grcp-server/main.go

cd golang && cd cmd && cd python && source env/bin/activate
python server/python_server.py
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. python_proto/claude/claude.proto
go run cmd/client/main.go
======================================= postgres ======================
psql -U myuser -d mydatabase

# psql -U myuser -d mydatabase

DROP DATABASE mydatabase;
TRUNCATE TABLE schema\*migrations, users; delete all data
\dt : list all table
\d guests
\d users
\d comments
\d sessions
\d reading_test_models;

\d orders

SELECT _ FROM tables;
SELECT _ FROM set*dishes;
SELECT * FROM sets;
SELECT _ FROM dishes;
SELECT \* FROM orders;
SELECT _ FROM users;
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
git branch web-sokcert
git checkout web-sokcert

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ecomm-grpc/proto/python_proto/claude/claude.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ecomm-grpc/proto/python_proto/helloworld.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ecomm-grpc-python/ielts/proto/ielts.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ecomm-grpc/proto/claude/claude.proto

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

Jump back to the golang branch:
git checkout test_isadmin

Merge the golang branch with the python branch:
Jump back to the golang branch:
git checkout test_isadmin

Merge the golang branch with the python branch:
git merge guest
git merge --no-ff guest

Update the changes to the remote repository:
git push origin test_isadmin

Jump back to the python branch:
git checkout guest

git branch
========================================= golang ==============================

====================================== project proto ============================

cd project_protos

go mod init project_proto

source env/bin/activate

cd python
python server/greeter_server.py

python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. python_proto/helloworld.proto

python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. python_proto/claude/claude.proto

------------------------------------- quan an qr ------------
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative quanqr/proto_qr/delivery/delivery.proto


protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative quanqr/proto_qr/set/set.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative quanqr/proto_qr/account/account.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative quanqr/proto_qr/dish/dish.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative quanqr/proto_qr/dishsnapshot/dishsnapshot.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative quanqr/proto_qr/guest/guest.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative quanqr/proto_qr/order/order.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative quanqr/proto_qr/table/table.proto

http://localhost:8888/images/image?filename=Screenshot%202024-02-20%20at%2014.37.22.png&path=folder1/folder2

=============================== test ========================

stand at python

git checkout -b test_isadmin

http://localhost:3000/admin/dished

Exit the editor: If you’re using vim (which is the default editor for Git), you can quit by:
Pressing Esc to ensure you’re in normal mode.
Typing :q! and pressing Enter to quit without saving changes.
Abort the merge: If you want to abort the merge entirely, you can run:
git merge --abort

If you need to write a proper commit message, you can edit the message above the lines starting with #. For example:

Merge branch 'test_isadmin' into python

This merge is necessary to integrate the latest changes from the 'test_isadmin' branch into the 'python' branch.

git checkout -b testferetur---set--add--database

const ws = new WebSocket(`ws://your-server/ws?userId=${userId}&userName=${userName}`);
// For User 1
const ws1 = new WebSocket('ws://your-server/ws?userId=user1&userName=John');

// For User 2
const ws2 = new WebSocket('ws://your-server/ws?userId=user2&userName=Jane');


Connect first user: ws://your-server/ws?userId=user1&userName=John
Connect second user: ws://your-server/ws?userId=user2&userName=Jane

{
  "fromUser": "user1",
  "toUser": "user2",
  "content": "Test direct message"
}