package main

import (
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"

	comment_api "english-ai-full/ecomm-api/comment-api"
	python_api "english-ai-full/ecomm-api/python-api"
	python_ielts "english-ai-full/ecomm-api/python-ielts"
	reading_api "english-ai-full/ecomm-api/reading-api"
	user_api "english-ai-full/ecomm-api/user-api"
	websocket_handler "english-ai-full/ecomm-api/websocket/websocket_handler"
	image_upload "english-ai-full/upload/image"

	"github.com/go-chi/cors"

	"english-ai-full/ecomm-api/websocket/websocket_repository"
	"english-ai-full/ecomm-api/websocket/websocket_service"
	"english-ai-full/ecomm-grpc/config"
	pb "english-ai-full/ecomm-grpc/proto"
	pb_python "english-ai-full/ecomm-grpc/proto/python_proto"
	pb_python_ielts "english-ai-full/ecomm-grpc/proto/python_proto/claude"

	pb_comment "english-ai-full/ecomm-grpc/proto/comment"
	pb_reading "english-ai-full/ecomm-grpc/proto/reading"
	dish "english-ai-full/quanqr/dish"
	pb_dish "english-ai-full/quanqr/proto_qr/dish"

	// "github.com/go-chi/chi"

	"github.com/go-chi/chi"
	"github.com/ianschenck/envflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"english-ai-full/ecomm-api/handler"
	"english-ai-full/ecomm-api/route"
)

const minSecretKeySize = 32

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var (
		secretKey = envflag.String("SECRET_KEY", "01234567890123456789012345678901", "secret key for JWT signing")
		svcAddr   = envflag.String("GRPC_SVC_ADDR", cfg.GRPCAddress, "address where the ecomm-grpc service is listening on")
	)
	envflag.Parse()

	if len(*secretKey) < minSecretKeySize {
		log.Fatalf("SECRET_KEY must be at least %d characters", minSecretKeySize)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}


r := chi.NewRouter()

r.Use(cors.Handler(cors.Options{
    AllowedOrigins:   []string{"*"},  // Allow all origins for development
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposedHeaders:   []string{"Link"},
    AllowCredentials: true,
    MaxAge:           300,
}))
// python server ---------------------
python_conn, err := grpc.NewClient(":50052", opts...)
if err != nil {
	log.Fatalf("failed to connect to Python gRPC server: %v", err)
}
defer python_conn.Close()

// Python greeter
python_client := pb_python.NewGreeterClient(python_conn)


python_hdl := python_api.NewPythonHandler(python_client)


python_api.RegisterPythonRoutes(r, python_hdl)



// python ielts service 




python_client_ielts := pb_python_ielts.NewIELTSServiceClient(python_conn)

python_hdl_ielts := python_ielts.NewPythonIeltsHandler(python_client_ielts)
python_ielts.RegisterPythonIeltsRoutes(r,python_hdl_ielts)
//  ---------------------------
	conn, err := grpc.NewClient(*svcAddr, opts...)
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

// reading service

setupReadingService(r, conn, secretKey)
// client_reading := pb_reading.NewEcommReadingClient(conn)

// hdl_reading := reading_api.NewReadingHandler(client_reading, *secretKey)
// reading_api.RegisterReadingRoutes(r, hdl_reading)
//  user handler
client := pb.NewEcommUserClient(conn)
	hdl := handler.NewHandler(client, *secretKey)
	
	route.RegisterRoutes(r, hdl)


	// comment 
	client_comment := pb_comment.NewCommentServiceClient(conn)

	hdl_comment := comment_api.NewCommentHandler(client_comment, *secretKey)
	comment_api.RegisterCommentRoutes(r, hdl_comment)


// new user  handler
hdl_NewUser := user_api.NewHandlerUser(client, *secretKey)

user_api.RegisterRoutesUser(r, hdl_NewUser)
// web socket

setupWebSocketService(r, )

// websockrepo := websocket_repository.NewInMemoryMessageRepository()
// websocketService := websocket_service.NewWebSocketService(websockrepo)
// go websocketService.Run()

// websocketHandler := websocket_handler.NewWebSocketHandler(websocketService)

// r.Get("/ws", websocketHandler.HandleWebSocket)


// dish

dish_client := pb_dish.NewDishServiceClient(conn)
	dish_hdl := dish.NewDishHandler(dish_client, *secretKey)
	
	dish.RegisterDishRoutes(r, dish_hdl)

	//
    r.Get("/image", func(w http.ResponseWriter, r *http.Request) {


        file, err := os.Open("upload/quananqr/public/pexels-ella-olsson-572949-1640777.jpg")
        if err != nil {
            http.Error(w, "Image not found.", http.StatusNotFound)
            return
        }
        defer file.Close()

        img, _, err := image.Decode(file)
        if err != nil {
            http.Error(w, "Error decoding image.", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "image/jpeg")
        jpeg.Encode(w, img, nil)
    })


	// 

	hdl_image := image_upload.NewImageHandler( *secretKey)

	image_upload.RegisterImageRoutes(r, hdl_image)

route.Start(":8888", r)


}


func setupReadingService(r *chi.Mux, conn *grpc.ClientConn, secretKey *string) {
	client_reading := pb_reading.NewEcommReadingClient(conn)

	hdl_reading := reading_api.NewReadingHandler(client_reading, *secretKey)
	reading_api.RegisterReadingRoutes(r, hdl_reading)
}



func setupWebSocketService(r *chi.Mux) {
    websockrepo := websocket_repository.NewInMemoryMessageRepository()
    websocketService := websocket_service.NewWebSocketService(websockrepo)
    go websocketService.Run()

    websocketHandler := websocket_handler.NewWebSocketHandler(websocketService)
    r.Get("/ws", websocketHandler.HandleWebSocket)
}

