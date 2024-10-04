package order_grpc

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"english-ai-full/quanqr/proto_qr/order"
	"english-ai-full/token"

	"github.com/go-chi/chi/v5"
)


type OrderHandlerController struct {
	ctx        context.Context
	client     order.OrderServiceClient
	TokenMaker *token.JWTMaker // Update this to match your actual gRPC client interface
}

func NewOrderHandler(client order.OrderServiceClient, secretKey string) *OrderHandlerController {
	return &OrderHandlerController{
		ctx:    context.Background(),
		client: client,

		TokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (h *OrderHandlerController) CreateOrders(w http.ResponseWriter, r *http.Request) {
	var req order.CreateOrdersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	log.Println("Creating orders for guest:", req.GuestId)
	createdOrders, err := h.client.CreateOrders(h.ctx, &req)
	if err != nil {
		log.Println("Error creating orders:", err)
		http.Error(w, "error creating orders", http.StatusInternalServerError)
		return
	}

	res := createdOrders // Modify this based on how you want to structure your response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *OrderHandlerController) GetOrders(w http.ResponseWriter, r *http.Request) {
	var req order.GetOrdersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	log.Println("Fetching orders from", req.FromDate, "to", req.ToDate)
	orders, err := h.client.GetOrders(h.ctx, &req)
	if err != nil {
		log.Println("Error fetching orders:", err)
		http.Error(w, "error fetching orders", http.StatusInternalServerError)
		return
	}

	res := orders // Modify this based on your response structure
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *OrderHandlerController) GetOrderDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	orderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing ID", http.StatusBadRequest)
		return
	}

	orderDetail, err := h.client.GetOrderDetail(h.ctx, &order.OrderDetailIdParam{Id: orderID})
	if err != nil {
		log.Println("Error fetching order detail:", err)
		http.Error(w, "error fetching order detail", http.StatusInternalServerError)
		return
	}

	res := orderDetail // Modify this based on your response structure
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *OrderHandlerController) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var req order.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	updatedOrder, err := h.client.UpdateOrder(h.ctx, &req)
	if err != nil {
		log.Println("Error updating order:", err)
		http.Error(w, "error updating order", http.StatusInternalServerError)
		return
	}

	res := updatedOrder // Modify this based on your response structure
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *OrderHandlerController) PayGuestOrders(w http.ResponseWriter, r *http.Request) {
	var req order.PayGuestOrdersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	paidOrders, err := h.client.PayGuestOrders(h.ctx, &req)
	if err != nil {
		log.Println("Error paying guest orders:", err)
		http.Error(w, "error paying guest orders", http.StatusInternalServerError)
		return
	}

	res := paidOrders // Modify this based on your response structure
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
