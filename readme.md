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
SELECT * FROM orders;
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
git branch make_order
git checkout make_order

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

{
"guest_id": null,
"user_id": 1,
"is_guest": false,
"table_number": 5,
"order_handler_id": 1,
"status": "pending",
"created_at": "2024-10-19T10:00:00Z",
"updated_at": "2024-10-19T10:00:00Z",
"total_price": 63,
"dish_items": [
{
"id": 1,
"quantity": 2,
"dish": {
"id": 1,
"name": "trung",
"price": 9,
"description": "Classic Italian pasta dish with eggs, cheese, pancetta, and black pepper",
"image": "https://example.com/spaghetti-carbonara.jpg",
"status": "available",
"created_at": "2024-10-17T08:50:27.304909Z",
"updated_at": "2024-10-17T08:50:27.304909Z"
}
},
{
"id": 2,
"quantity": 1,
"dish": {
"id": 2,
"name": "do",
"price": 9,
"description": "Classic Italian pasta dish with eggs, cheese, pancetta, and black pepper",
"image": "https://example.com/spaghetti-carbonara.jpg",
"status": "available",
"created_at": "2024-10-17T08:50:32.470107Z",
"updated_at": "2024-10-17T08:50:32.470107Z"
}
}
],
"set_items": [
{
"id": 1,
"quantity": 1,
"set": {
"id": 1,
"name": "My New Set 1212",
"description": "A delicious combination of dishes",
"dishes": [
{
"id": 1,
"name": "trung",
"price": 9,
"description": "Classic Italian pasta dish with eggs, cheese, pancetta, and black pepper",
"image": "https://example.com/spaghetti-carbonara.jpg",
"status": "available",
"created_at": "2024-10-17T08:50:27.304909Z",
"updated_at": "2024-10-17T08:50:27.304909Z"
},
{
"id": 2,
"name": "do",
"price": 9,
"description": "Classic Italian pasta dish with eggs, cheese, pancetta, and black pepper",
"image": "https://example.com/spaghetti-carbonara.jpg",
"status": "available",
"created_at": "2024-10-17T08:50:32.470107Z",
"updated_at": "2024-10-17T08:50:32.470107Z"
}
],
"user_id": 1,
"created_at": "2024-10-17T08:50:48.249115Z",
"updated_at": "2024-10-17T08:50:48.249115Z",
"is_favourite": false,
"like_by": [],
"is_public": true,
"image": "asdfasdfasdgfasdg"
}
}
]
}

import {
OrderDetailedResponse,
OrderSetDetailed,
OrderDetailedDish
} from "@/schemaValidations/interface/type_order";
import { ColumnDef } from "@tanstack/react-table";
import {
Select,
SelectContent,
SelectItem,
SelectTrigger,
SelectValue
} from "@/components/ui/select";
import { useState } from "react";
import { Input } from "@/components/ui/input";

const ORDER_STATUSES = ["ORDERING", "SERVING", "WAITING", "DONE"] as const;
type OrderStatus = (typeof ORDER_STATUSES)[number];

const PAYMENT_METHODS = ["CASH", "TRANSFER"] as const;
type PaymentMethod = (typeof PAYMENT_METHODS)[number];

interface TableMeta {
onStatusChange?: (orderId: number, newStatus: string) => void;
onPaymentMethodChange?: (orderId: number, newMethod: string) => void;
}

export const columns: ColumnDef<OrderDetailedResponse, any>[] = [
{
accessorKey: "order_name",
header: "Name",
cell: ({ row }) => (

<div className="font-medium">#{row.getValue("order_name")}</div>
)
},
{
accessorKey: "table_number",
header: "Table/away",
cell: ({ row }) => (
<div
className={`text-center ${
          row.original.takeAway ? "bg-orange-600 rounded-md px-2 py-1" : ""
        }`} >
{row.getValue("table_number")}
</div>
)
},
{
accessorKey: "status",
header: "Status",
cell: ({ row, table }) => {
const [selectedStatus, setSelectedStatus] = useState<OrderStatus>(
row.getValue("status") as OrderStatus
);

      const statusStyles: Record<OrderStatus, string> = {
        ORDERING: "bg-blue-100 text-blue-800",
        SERVING: "bg-yellow-100 text-yellow-800",
        WAITING: "bg-orange-100 text-orange-800",
        DONE: "bg-green-100 text-green-800"
      };

      const meta = table.options.meta as TableMeta;

      return (
        <Select
          value={selectedStatus}
          onValueChange={(newStatus: OrderStatus) => {
            setSelectedStatus(newStatus);
            meta?.onStatusChange?.(row.original.id, newStatus);
          }}
        >
          <SelectTrigger
            className={`w-[120px] h-8 ${statusStyles[selectedStatus]}`}
          >
            <SelectValue>{selectedStatus}</SelectValue>
          </SelectTrigger>
          <SelectContent>
            {ORDER_STATUSES.map((orderStatus) => (
              <SelectItem
                key={orderStatus}
                value={orderStatus}
                className={statusStyles[orderStatus]}
              >
                {orderStatus}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      );
    }

},
{
accessorKey: "payment_method",
header: "Payment",
cell: ({ row, table }) => {
const [selectedPayment, setSelectedPayment] = useState<PaymentMethod>(
(row.getValue("payment_method") as PaymentMethod) || "CASH"
);

      const paymentStyles: Record<PaymentMethod, string> = {
        CASH: "bg-emerald-50 text-emerald-700",
        TRANSFER: "bg-indigo-50 text-indigo-700"
      };

      const meta = table.options.meta as TableMeta;

      return (
        <Select
          value={selectedPayment}
          onValueChange={(newMethod: PaymentMethod) => {
            setSelectedPayment(newMethod);
            meta?.onPaymentMethodChange?.(row.original.id, newMethod);
          }}
        >
          <SelectTrigger
            className={`w-[120px] h-8 ${paymentStyles[selectedPayment]}`}
          >
            <SelectValue>
              {selectedPayment === "CASH" ? "Cash" : "Transfer"}
            </SelectValue>
          </SelectTrigger>
          <SelectContent>
            {PAYMENT_METHODS.map((method) => (
              <SelectItem
                key={method}
                value={method}
                className={paymentStyles[method]}
              >
                {method === "CASH" ? "Cash" : "Transfer"}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      );
    }

},
{
accessorKey: "data_set",
header: "Sets",
cell: ({ row }) => {
const sets = row.getValue("data_set") as OrderSetDetailed[];
return (

<div className="space-y-1">
{sets.map((set) => (
<div key={set.id} className="text-sm">
{set.quantity}x {set.name} (${set.price})
            </div>
          ))}
        </div>
      );
    }
  },
  {
    accessorKey: "data_dish",
    header: "Individual Dishes",
    cell: ({ row }) => {
      const dishes = row.getValue("data_dish") as OrderDetailedDish[];
      return (
        <div className="space-y-1">
          {dishes.map((dish, index) => (
            <div key={`${dish.dish_id}-${index}`} className="text-sm">
              {dish.quantity}x {dish.name} (${dish.price})
</div>
))}
</div>
);
}
},
{
accessorKey: "bow_details",
header: "Bowl Details",
cell: ({ row }) => {
const withChili = row.original.bow_chili;
const noChili = row.original.bow_no_chili;
const total = withChili + noChili;
const isTakeAway = row.original.takeAway;
const chiliNumber = row.original.chiliNumber;

      return total > 0 || (isTakeAway && chiliNumber > 0) ? (
        <div className="space-y-1 text-sm">
          {withChili > 0 && <div>With Chili: {withChili}</div>}
          {noChili > 0 && <div>No Chili: {noChili}</div>}
          {isTakeAway && chiliNumber > 0 && (
            <div className="font-medium">Takeaway Chili: {chiliNumber}</div>
          )}
        </div>
      ) : null;
    }

},
{
accessorKey: "total_price",
header: "Total & Payment",
cell: ({ row }) => {
const totalPrice = row.getValue("total_price") as number;
const [amountPaid, setAmountPaid] = useState<string>("");
const [change, setChange] = useState<number | null>(null);

      const handlePaymentInput = (value: string) => {
        setAmountPaid(value);
        const numericValue = parseFloat(value) || 0;
        const changeAmount = numericValue - totalPrice;
        setChange(changeAmount >= 0 ? changeAmount : null);
      };

      return (
        <div className="space-y-2">
          <div className="font-medium text-right">Total: ${totalPrice}</div>
          <div className="flex items-center gap-2">
            <Input
              type="number"
              placeholder="Amount paid"
              value={amountPaid}
              onChange={(e) => handlePaymentInput(e.target.value)}
              className="w-24 h-8 text-right"
            />
            <span className="text-sm">$</span>
          </div>
          {change !== null && (
            <div
              className={`text-right text-sm ${
                change >= 0 ? "text-green-600" : "text-red-600"
              }`}
            >
              Change: ${change.toFixed(2)}
            </div>
          )}
        </div>
      );
    }

}
];
