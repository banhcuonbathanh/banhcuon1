package main

import (
	"context"
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
	ws2 "english-ai-full/quanqr/ws2"
	"english-ai-full/token"

	image_upload "english-ai-full/upload/image"

	"github.com/go-chi/cors"

	pb "english-ai-full/ecomm-grpc/proto"
	pb_python "english-ai-full/ecomm-grpc/proto/python_proto"
	pb_python_ielts "english-ai-full/ecomm-grpc/proto/python_proto/claude"
	"english-ai-full/util"

	pb_comment "english-ai-full/ecomm-grpc/proto/comment"
	pb_reading "english-ai-full/ecomm-grpc/proto/reading"
	dish "english-ai-full/quanqr/dish"
	pb_dish "english-ai-full/quanqr/proto_qr/dish"

	pb_set "english-ai-full/quanqr/proto_qr/set"
	set "english-ai-full/quanqr/set"

	pb_guests "english-ai-full/quanqr/proto_qr/guest"
	guests "english-ai-full/quanqr/qr_guests"

	order "english-ai-full/quanqr/order"
	pb_order "english-ai-full/quanqr/proto_qr/order"

	// delivery
	delivery "english-ai-full/quanqr/delivery"
	pb_delivery "english-ai-full/quanqr/proto_qr/delivery"

	// "github.com/go-chi/chi"

	pb_tables "english-ai-full/quanqr/proto_qr/table"
	tables "english-ai-full/quanqr/tables"

	//----------

	"github.com/go-chi/chi"
	"github.com/ianschenck/envflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"english-ai-full/ecomm-api/handler"
	"english-ai-full/ecomm-api/route"
)



func main() {

	cfg, err := util.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}


	
	envflag.Parse()


	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}


r := chi.NewRouter()

r.Use(cors.Handler(cors.Options{
	AllowedOrigins:   []string{cfg.QuanAnAddress, "http://localhost:*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
	AllowedHeaders:   []string{
		"Accept",
		"Authorization",
		"Content-Type", 
		"X-CSRF-Token",
		"X-Table-Token",
		"X-Requested-With",
	},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
}))

// Use environment variable with a default value
if getEnvWithDefault("GO_ENV", "development") == "development" {
	r.Use(debugMiddleware)
}

setupGlobalMiddleware(r, cfg)
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
conn, err := grpc.DialContext(
    context.Background(),
    cfg.GRPCAddress,
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
if err != nil {
    log.Fatalf("failed to connect to server: %v", err)
}
defer conn.Close()

// reading service

setupReadingService(r, conn, &cfg.QuanAnJWTsecretKey)
// client_reading := pb_reading.NewEcommReadingClient(conn)

// hdl_reading := reading_api.NewReadingHandler(client_reading, *secretKey)
// reading_api.RegisterReadingRoutes(r, hdl_reading)
//  user handler
client := pb.NewEcommUserClient(conn)
hdl := handler.NewHandler(client, cfg.QuanAnJWTsecretKey)
	
	route.RegisterRoutes(r, hdl)


	// comment 
	client_comment := pb_comment.NewCommentServiceClient(conn)

	hdl_comment := comment_api.NewCommentHandler(client_comment, cfg.QuanAnJWTsecretKey)
	comment_api.RegisterCommentRoutes(r, hdl_comment)


// new user  handler
hdl_NewUser := user_api.NewHandlerUser(client, cfg.QuanAnJWTsecretKey)

user_api.RegisterRoutesUser(r, hdl_NewUser)
// web socket





set_client := pb_set.NewSetServiceClient(conn)
set_hdl := set.NewSetHandler(set_client, cfg.QuanAnJWTsecretKey)
	
set.RegisterSetRoutes(r, set_hdl)

// dish

dish_client := pb_dish.NewDishServiceClient(conn)
	dish_hdl := dish.NewDishHandler(dish_client, cfg.QuanAnJWTsecretKey)
	
	dish.RegisterDishRoutes(r, dish_hdl)

	// table


	table_client := pb_tables.NewTableServiceClient(conn)
	table_hdl := tables.NewTableHandler(table_client)
	
	tables.RegisterTablesRoutes(r, table_hdl)

	// guest
	guests_client := pb_guests.NewGuestServiceClient(conn)
	guests_hdl := guests.NewGuestHandler(guests_client,cfg.QuanAnJWTsecretKey)
	
	guests.RegisterGuestRoutes(r, guests_hdl)
// order

order_client := pb_order.NewOrderServiceClient(conn)
order_hdl := order.NewOrderHandler(order_client, cfg.QuanAnJWTsecretKey)

order.RegisterOrderRoutes(r, order_hdl)


// delivery

delivery_client := pb_delivery.NewDeliveryServiceClient(conn)
delivery_hdl := delivery.NewDeliveryHandler(delivery_client, cfg.QuanAnJWTsecretKey)

delivery.RegisterDeliveryRoutes(r, delivery_hdl)

SetupWs2(r, order_hdl, delivery_hdl, cfg)
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

	hdl_image := image_upload.NewImageHandler( cfg.QuanAnJWTsecretKey)

	image_upload.RegisterImageRoutes(r, hdl_image)
	
route.Start(cfg.ServerAddress, r)


}


func setupReadingService(r *chi.Mux, conn *grpc.ClientConn, secretKey *string) {
	client_reading := pb_reading.NewEcommReadingClient(conn)

	hdl_reading := reading_api.NewReadingHandler(client_reading, *secretKey)
	reading_api.RegisterReadingRoutes(r, hdl_reading)
}





func setupGlobalMiddleware(r *chi.Mux, cfg *util.Config) {
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Set CORS headers for every response
            w.Header().Set("Access-Control-Allow-Origin", cfg.QuanAnAddress)
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, X-Table-Token")
            w.Header().Set("Access-Control-Allow-Credentials", "true")

            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }

            next.ServeHTTP(w, r)
        })
    })
}

func debugMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
        // log.Printf("Headers: %v", r.Header)
        next.ServeHTTP(w, r)
    })
}

func getEnvWithDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}






func SetupWs2(r chi.Router, orderHandler *order.OrderHandlerController, deliveryHandler *delivery.DeliveryHandlerController, cfg *util.Config) {
    log.Println("golang/cmd/server/main.go")

    // Initialize the JWT token maker
 
    tokenMaker := token.NewJWTMaker(cfg.QuanAnJWTsecretKey)

    // Create message handlers
    orderMsgHandler := ws2.NewOrderMessageHandler(orderHandler)
    deliveryMsgHandler := ws2.NewDeliveryMessageHandler(deliveryHandler)

    // Create a combined message handler
    combinedHandler := ws2.NewCombinedMessageHandler(orderMsgHandler, deliveryMsgHandler)

    // Create and setup the hub
    hub := ws2.NewHub(combinedHandler)
    broadcaster := ws2.NewBroadcaster(hub)

    // Set broadcasters
    orderMsgHandler.SetBroadcaster(broadcaster)
    deliveryMsgHandler.SetBroadcaster(broadcaster)

    // Setup router with token maker
    // wsRouter := ws2.NewWebSocketRouter(hub)

	wsRouter := ws2.NewWebSocketRouter(hub, tokenMaker)
    wsRouter.RegisterRoutes(r)

    go hub.Run()
}