package order_grpc

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"google.golang.org/protobuf/types/known/timestamppb"
	"english-ai-full/token"
	"english-ai-full/logger"
	"english-ai-full/quanqr/proto_qr/order"
)

type OrderHandlerController struct {
	ctx     context.Context
	service  order.OrderServiceClient
	logger  *logger.Logger
	TokenMaker *token.JWTMaker
}

func NewOrderHandler(service order.OrderServiceClient, secretKey string) *OrderHandlerController {
	return &OrderHandlerController{
		ctx:     context.Background(),
		service: service,
		logger:  logger.NewLogger(),
		TokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (h *OrderHandlerController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequestType
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	h.logger.Info("Creating new order")
	createdOrderResponse, err := h.service.CreateOrder(h.ctx, ToCreateOrderRequest(req))
	if err != nil {
		h.logger.Error("Error creating order: " + err.Error())
		http.Error(w, "error creating order", http.StatusInternalServerError)
		return
	}

	res := OrderResponse{Data: ToOrderType(createdOrderResponse.Data)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}


func (h *OrderHandlerController) GetOrders(w http.ResponseWriter, r *http.Request) {
	var req GetOrdersRequestType
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	h.logger.Info("Fetching orders")
	ordersResponse, err := h.service.GetOrders(h.ctx, ToGetOrdersRequest(req))
	if err != nil {
		h.logger.Error("Error fetching orders: " + err.Error())
		http.Error(w, "error fetching orders", http.StatusInternalServerError)
		return
	}

	res := OrderListResponse{Data: ToOrderTypeList(ordersResponse.Data)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *OrderHandlerController) GetOrderDetail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid order ID", http.StatusBadRequest)
		return
	}

	h.logger.Info("Fetching order detail")
	orderDetailResponse, err := h.service.GetOrderDetail(h.ctx, &order.OrderIdParam{Id: id})
	if err != nil {
		h.logger.Error("Error fetching order detail: " + err.Error())
		http.Error(w, "error fetching order detail", http.StatusInternalServerError)
		return
	}

	res := OrderResponse{Data: ToOrderType(orderDetailResponse.Data)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}


func (h *OrderHandlerController) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var req UpdateOrderRequestType
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	h.logger.Info("Updating order")
	updatedOrderResponse, err := h.service.UpdateOrder(h.ctx, ToUpdateOrderRequest(req))
	if err != nil {
		h.logger.Error("Error updating order: " + err.Error())
		http.Error(w, "error updating order", http.StatusInternalServerError)
		return
	}

	res := OrderResponse{Data: ToOrderType(updatedOrderResponse.Data)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *OrderHandlerController) PayOrders(w http.ResponseWriter, r *http.Request) {
	var req PayOrdersRequestType
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	h.logger.Info("Processing payment for orders")
	paidOrdersResponse, err := h.service.PayOrders(h.ctx, ToPayOrdersRequest(req))
	if err != nil {
		h.logger.Error("Error processing payment: " + err.Error())
		http.Error(w, "error processing payment", http.StatusInternalServerError)
		return
	}

	res := OrderListResponse{Data: ToOrderTypeList(paidOrdersResponse.Data)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// Conversion functions

func ToCreateOrderRequest(req CreateOrderRequestType) *order.CreateOrderRequest {
	return &order.CreateOrderRequest{
		GuestId:        req.GuestID,
		UserId:         req.UserID,
		IsGuest:        req.IsGuest,
		TableNumber:    req.TableNumber,
		OrderHandlerId: req.OrderHandlerID,
		Status:         req.Status,
		CreatedAt:      timestamppb.New(time.Now()),
		UpdatedAt:      timestamppb.New(time.Now()),
		TotalPrice:     req.TotalPrice,
		DishItems:      ToDishOrderItems(req.DishItems),
		SetItems:       ToSetOrderItems(req.SetItems),
	}
}

func ToGetOrdersRequest(req GetOrdersRequestType) *order.GetOrdersRequest {
	return &order.GetOrdersRequest{
		FromDate: timestamppb.New(req.FromDate),
		ToDate:   timestamppb.New(req.ToDate),
		UserId:   req.UserID,
		GuestId:  req.GuestID,
	}
}

func ToUpdateOrderRequest(req UpdateOrderRequestType) *order.UpdateOrderRequest {
	return &order.UpdateOrderRequest{
		Id:             req.ID,
		GuestId:        req.GuestID,
		UserId:         req.UserID,
		TableNumber:    req.TableNumber,
		OrderHandlerId: req.OrderHandlerID,
		Status:         req.Status,
		TotalPrice:     req.TotalPrice,
		DishItems:      ToDishOrderItemsProto(req.DishItems),
		SetItems:       ToSetOrderItemsProto(req.SetItems),
		IsGuest:        req.IsGuest,
	}
}

func ToPayOrdersRequest(req PayOrdersRequestType) *order.PayOrdersRequest {
	payReq := &order.PayOrdersRequest{}
	if req.GuestID != nil {
		payReq.Identifier = &order.PayOrdersRequest_GuestId{GuestId: *req.GuestID}
	} else if req.UserID != nil {
		payReq.Identifier = &order.PayOrdersRequest_UserId{UserId: *req.UserID}
	}
	return payReq
}

func ToOrderType(o *order.Order) OrderType {
	return OrderType{
		ID:             o.Id,
		GuestID:        o.GuestId,
		UserID:         o.UserId,
		IsGuest:        o.IsGuest,
		TableNumber:    o.TableNumber,
		OrderHandlerID: o.OrderHandlerId,
		Status:         o.Status,
		CreatedAt:      o.CreatedAt.AsTime(),
		UpdatedAt:      o.UpdatedAt.AsTime(),
		TotalPrice:     o.TotalPrice,
		DishItems:      ToDishOrderItemTypes(o.DishItems),
		SetItems:       ToSetOrderItemTypes(o.SetItems),
	}
}

func ToOrderTypeList(orders []*order.Order) []OrderType {
	result := make([]OrderType, len(orders))
	for i, o := range orders {
		result[i] = ToOrderType(o)
	}
	return result
}

func ToDishOrderItems(items []CreateOrderItemType) []*order.DishOrderItem {
	result := make([]*order.DishOrderItem, len(items))
	for i, item := range items {
		result[i] = &order.DishOrderItem{
			Quantity: item.Quantity,
		}
	}
	return result
}

func ToSetOrderItems(items []CreateOrderItemType) []*order.SetOrderItem {
	result := make([]*order.SetOrderItem, len(items))
	for i, item := range items {
		result[i] = &order.SetOrderItem{
			Quantity: item.Quantity,
		}
	}
	return result
}

func ToDishOrderItemTypes(items []*order.DishOrderItem) []DishOrderItem {
	result := make([]DishOrderItem, len(items))
	for i, item := range items {
		result[i] = DishOrderItem{
			ID:       item.Id,
			Quantity: item.Quantity,
			Dish:     ToDishOrderType(item.Dish),
		}
	}
	return result
}

func ToSetOrderItemTypes(items []*order.SetOrderItem) []SetOrderItemType {
	result := make([]SetOrderItemType, len(items))
	for i, item := range items {
		result[i] = SetOrderItemType{
			ID:       item.Id,
			Quantity: item.Quantity,
			Set:      ToSetProtoType(item.Set),
		}
	}
	return result
}

func ToDishOrderType(d *order.DishOrder) DishOrderType {
	return DishOrderType{
		ID:          d.Id,
		Name:        d.Name,
		Price:       d.Price,
		Description: d.Description,
		Image:       d.Image,
		Status:      d.Status,
		CreatedAt:   d.CreatedAt.AsTime(),
		UpdatedAt:   d.UpdatedAt.AsTime(),
	}
}

func ToSetProtoType(s *order.SetProto) SetProtoType {
	return SetProtoType{
		ID:          int64(s.Id),
		Name:        s.Name,
		Description: s.Description,
		UserID:      s.UserId,
		IsFavourite: s.IsFavourite,
		LikeBy:      s.LikeBy,
		CreatedAt:   s.CreatedAt.AsTime(),
		UpdatedAt:   s.UpdatedAt.AsTime(),
		IsPublic:    s.IsPublic,
		Image:       s.Image,
		Dishes:      ToSetProtoDishTypes(s.Dishes),
	}
}

func ToSetProtoDishTypes(dishes []*order.SetProtoDish) []SetProtoDishType {
	result := make([]SetProtoDishType, len(dishes))
	for i, d := range dishes {
		result[i] = SetProtoDishType{
			ID:    d.Id,
			Name:  d.Name,
			Price: d.Price,
		}
	}
	return result
}

func ToDishOrderItemsProto(items []DishOrderItem) []*order.DishOrderItem {
	result := make([]*order.DishOrderItem, len(items))
	for i, item := range items {
		result[i] = &order.DishOrderItem{
			Id:       item.ID,
			Quantity: item.Quantity,
			Dish:     ToDishOrderProto(item.Dish),
		}
	}
	return result
}

func ToSetOrderItemsProto(items []SetOrderItemType) []*order.SetOrderItem {
	result := make([]*order.SetOrderItem, len(items))
	for i, item := range items {
		result[i] = &order.SetOrderItem{
			Id:       item.ID,
			Quantity: item.Quantity,
			Set:      ToSetProtoProto(item.Set),
		}
	}
	return result
}

func ToDishOrderProto(dish DishOrderType) *order.DishOrder {
	return &order.DishOrder{
		Id:          dish.ID,
		Name:        dish.Name,
		Price:       dish.Price,
		Description: dish.Description,
		Image:       dish.Image,
		Status:      dish.Status,
		CreatedAt:   timestamppb.New(dish.CreatedAt),
		UpdatedAt:   timestamppb.New(dish.UpdatedAt),
	}
}

func ToSetProtoProto(set SetProtoType) *order.SetProto {
	return &order.SetProto{
		Id:          int32(set.ID),
		Name:        set.Name,
		Description: set.Description,
		UserId:      set.UserID,
		IsFavourite: set.IsFavourite,
		LikeBy:      set.LikeBy,
		CreatedAt:   timestamppb.New(set.CreatedAt),
		UpdatedAt:   timestamppb.New(set.UpdatedAt),
		IsPublic:    set.IsPublic,
		Image:       set.Image,
		Dishes:      ToSetProtoDishesProto(set.Dishes),
	}
}
func ToSetProtoDishesProto(dishes []SetProtoDishType) []*order.SetProtoDish {
	result := make([]*order.SetProtoDish, len(dishes))
	for i, dish := range dishes {
		result[i] = &order.SetProtoDish{
			Id:    dish.ID,
			Name:  dish.Name,
			Price: dish.Price,
		}
	}
	return result
}