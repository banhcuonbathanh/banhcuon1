package main

import (
	"english-ai-full/ecomm-grpc/config"
	"english-ai-full/ecomm-grpc/db"

	// "os"

	order "english-ai-full/quanqr/order"
	pb_order "english-ai-full/quanqr/proto_qr/order"

	//-----
	tables "english-ai-full/quanqr/tables"
	pb_table "english-ai-full/quanqr/proto_qr/table"
	//-----

//
	dish "english-ai-full/quanqr/dish"
	dishPb "english-ai-full/quanqr/proto_qr/dish"
	guests "english-ai-full/quanqr/qr_guests"
	pb_guests "english-ai-full/quanqr/proto_qr/guest"
	comment_repository "english-ai-full/ecomm-grpc/repository/comment_repository"
	reading_repository "english-ai-full/ecomm-grpc/repository/reading_repository"
	repository "english-ai-full/ecomm-grpc/repository/user_repository"
	comment_service "english-ai-full/ecomm-grpc/service/comment_service"
	reading_service "english-ai-full/ecomm-grpc/service/reading_service"
	service "english-ai-full/ecomm-grpc/service/user_service"
	"log"
	"net"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"

	pb "english-ai-full/ecomm-grpc/proto"
	comment_pb "english-ai-full/ecomm-grpc/proto/comment"
	reading_pb "english-ai-full/ecomm-grpc/proto/reading"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Run migrations
	
	if err := db.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}


	//--------------------------

// 	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
// if err != nil {
//     log.Fatalf("Failed to connect: %v", err)
// }
// defer conn.Close()

// client := reading_pb.NewEcommReadingClient(conn)
	/// -------------------
	// Connect to the database
	dbConn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()
// intiate userRepo userService
	userRepo := repository.NewUserRepository(dbConn)
	userService := service.NewUserServer(userRepo)
// intiate userRepo userService
readingRepo := reading_repository.NewReadingRepository(dbConn)
readingService := reading_service.NewReadingServer(readingRepo)

// intiate commentRepo commentService
commentRepo := comment_repository.NewCommentRepository(dbConn)
commentService := comment_service.NewCommentServer(commentRepo)
dishrepo := dish.NewDishRepository(dbConn)
dishService := dish.NewDishService(dishrepo)



guestsrepo := guests.NewGuestRepository(dbConn)
guestsService := guests.NewGuestService(guestsrepo)
// intiate commentRepo commentService

// order "english-ai-full/quanqr/order"
// pb_order "english-ai-full/quanqr/proto_qr/order"
orderrepo := order.NewOrderRepository(dbConn)
orderService := order.NewOrderService(orderrepo)
// pb_order "english-ai-full/quanqr/proto_qr/order"
tablerepo := tables.NewTableRepository(dbConn)
tableService := tables.NewTableService(tablerepo)

//

	lis, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb_table.RegisterTableServiceServer(grpcServer,tableService)

	dishPb.RegisterDishServiceServer(grpcServer,dishService)
	pb.RegisterEcommUserServer(grpcServer, userService)

	pb_guests.RegisterGuestServiceServer(grpcServer,guestsService)

	reading_pb.RegisterEcommReadingServer(grpcServer, readingService)

	comment_pb.RegisterCommentServiceServer(grpcServer, commentService)
	pb_order.RegisterOrderServiceServer(grpcServer,orderService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

