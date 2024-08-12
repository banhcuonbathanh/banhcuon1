go get -u github.com/jackc/pgx/v4
go get -u github.com/golang-migrate/migrate/v4
go get -u google.golang.org/grpc
go get -u github.com/spf13/viper

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/user.proto

go run cmd/server/main.go

chek env path.
monghoaivu@192 ~ % echo $(go env GOPATH)/bin
/Users/monghoaivu/go/bin
monghoaivu@192 ~ %

which protoc-gen-go
which protoc-gen-go-grpc
