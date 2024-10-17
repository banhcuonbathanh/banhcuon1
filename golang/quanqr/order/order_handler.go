package order_grpc

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"google.golang.org/protobuf/types/known/timestamppb"

	"english-ai-full/quanqr/proto_qr/order"
	"english-ai-full/token"
)

type OrderHandlerController struct {
	ctx        context.Context
	client     order.OrderServiceClient
	TokenMaker *token.JWTMaker
}

func NewOrderHandler(client order.OrderServiceClient, secretKey string) *OrderHandlerController {
	return &OrderHandlerController{
		ctx:        context.Background(),
		client:     client,
		TokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (h *OrderHandlerController) CreateOrders(w http.ResponseWriter, r *http.Request) {
	var orderReq CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	log.Println("handler CreateOrders before")
	createdOrders, err := h.client.CreateOrders(h.ctx, ToPBCreateOrderRequest(&orderReq))
	if err != nil {
		log.Println("handler CreateOrders err ", err)
		http.Error(w, "error creating orders in handler", http.StatusInternalServerError)
		return
	}
	log.Println("handler CreateOrders after")

	res := ToOrderResListFromPbOrderListResponse(createdOrders)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *OrderHandlerController) GetOrders(w http.ResponseWriter, r *http.Request) {
	fromDate := r.URL.Query().Get("from_date")
	toDate := r.URL.Query().Get("to_date")

	var fromTimestamp, toTimestamp *timestamppb.Timestamp

	if fromDate != "" {
		t, err := time.Parse(time.RFC3339, fromDate)
		if err != nil {
			http.Error(w, "invalid from_date format", http.StatusBadRequest)
			return
		}
		fromTimestamp = timestamppb.New(t)
	}

	if toDate != "" {
		t, err := time.Parse(time.RFC3339, toDate)
		if err != nil {
			http.Error(w, "invalid to_date format", http.StatusBadRequest)
			return
		}
		toTimestamp = timestamppb.New(t)
	}

	orders, err := h.client.GetOrders(h.ctx, &order.GetOrdersRequest{
		FromDate: fromTimestamp,
		ToDate:   toTimestamp,
	})
	if err != nil {
		http.Error(w, "failed to fetch orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := ToOrderResListFromPbOrderListResponse(orders)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandlerController) GetOrderDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing ID", http.StatusBadRequest)
		return
	}

	orderDetail, err := h.client.GetOrderDetail(h.ctx, &order.OrderDetailIdParam{Id: i})
	if err != nil {
		http.Error(w, "error getting order detail", http.StatusInternalServerError)
		return
	}

	res := ToOrderResFromPbOrderResponse(orderDetail)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *OrderHandlerController) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	updatedOrder, err := h.client.UpdateOrder(h.ctx, ToPBUpdateOrderRequest(&orderReq))
	if err != nil {
		http.Error(w, "error updating order", http.StatusInternalServerError)
		return
	}

	res := ToOrderResFromPbOrderResponse(updatedOrder)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *OrderHandlerController) PayGuestOrders(w http.ResponseWriter, r *http.Request) {
	guestID := chi.URLParam(r, "guest_id")
	id, err := strconv.ParseInt(guestID, 10, 64)
	if err != nil {
		http.Error(w, "error parsing guest ID", http.StatusBadRequest)
		return
	}

	paidOrders, err := h.client.PayGuestOrders(h.ctx, &order.PayGuestOrdersRequest{GuestId: id})
	if err != nil {
		http.Error(w, "error paying guest orders", http.StatusInternalServerError)
		return
	}

	res := ToOrderResListFromPbOrderListResponse(paidOrders)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// Helper functions

func ToPBCreateOrderRequest(req *CreateOrderRequest) *order.CreateOrderRequest {
	return &order.CreateOrderRequest{
		GuestId:         req.GuestID,
		TableNumber:     req.TableNumber,
		DishSnapshotId:  req.DishSnapshotID,
		OrderHandlerId:  req.OrderHandlerID,
		Status:          req.Status,
		CreatedAt:       timestamppb.New(req.CreatedAt),
		UpdatedAt:       timestamppb.New(req.UpdatedAt),
		TotalPrice:      req.TotalPrice,
		DishItems:       ToPBDishOrderItems(req.DishItems),
		SetItems:        ToPBSetOrderItems(req.SetItems),
	}
}

func ToPBUpdateOrderRequest(req *UpdateOrderRequest) *order.UpdateOrderRequest {
	return &order.UpdateOrderRequest{
		Id:              req.ID,
		GuestId:         req.GuestID,
		TableNumber:     req.TableNumber,
		DishSnapshotId:  req.DishSnapshotID,
		OrderHandlerId:  req.OrderHandlerID,
		Status:          req.Status,
		CreatedAt:       timestamppb.New(req.CreatedAt),
		UpdatedAt:       timestamppb.New(req.UpdatedAt),
		TotalPrice:      req.TotalPrice,
		DishItems:       ToPBDishOrderItems(req.DishItems),
		SetItems:        ToPBSetOrderItems(req.SetItems),
	}
}

func ToOrderResFromPbOrderResponse(pbRes *order.OrderResponse) OrderResponse {
	return OrderResponse{
		Data: ToOrderFromPbOrder(pbRes.Data),
	}
}

func ToOrderResListFromPbOrderListResponse(pbRes *order.OrderListResponse) OrderListResponse {
	orders := make([]Order, len(pbRes.Data))
	for i, pbOrder := range pbRes.Data {
		orders[i] = ToOrderFromPbOrder(pbOrder)
	}
	return OrderListResponse{
		Data: orders,
	}
}

func ToOrderFromPbOrder(pbOrder *order.Order) Order {
	return Order{
		ID:              pbOrder.Id,
		GuestID:         pbOrder.GuestId,
		TableNumber:     pbOrder.TableNumber,
		DishSnapshotID:  pbOrder.DishSnapshotId,
		OrderHandlerID:  pbOrder.OrderHandlerId,
		Status:          pbOrder.Status,
		CreatedAt:       pbOrder.CreatedAt.AsTime(),
		UpdatedAt:       pbOrder.UpdatedAt.AsTime(),
		TotalPrice:      pbOrder.TotalPrice,
		DishItems:       ToDishOrderItems(pbOrder.DishItems),
		SetItems:        ToSetOrderItems(pbOrder.SetItems),
	}
}

func ToPBDishOrderItems(items []DishOrderItem) []*order.DishOrderItem {
	pbItems := make([]*order.DishOrderItem, len(items))
	for i, item := range items {
		pbItems[i] = &order.DishOrderItem{
			Id:       item.ID,
			Quantity: item.Quantity,
			Dish:     ToPBDish(&item.Dish),
		}
	}
	return pbItems
}

func ToPBSetOrderItems(items []SetOrderItem) []*order.SetOrderItem {
	pbItems := make([]*order.SetOrderItem, len(items))
	for i, item := range items {
		pbItems[i] = &order.SetOrderItem{
			Id:             item.ID,
			Quantity:       item.Quantity,
			Set:            ToPBSetProto(&item.Set),
			ModifiedDishes: ToPBSetProtoDishes(item.ModifiedDishes),
		}
	}
	return pbItems
}

func ToPBDish(dish *Dish) *order.DishOrder {
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

func ToPBSetProto(set *SetProto) *order.SetProto {
	return &order.SetProto{
		Id:           int32(set.ID),
		Name:         set.Name,
		Description:  set.Description,
		Dishes:       ToPBSetProtoDishes(set.Dishes),
		UserId:       int32Ptr(set.UserID),
		CreatedAt:    timestamppb.New(set.CreatedAt),
		UpdatedAt:    timestamppb.New(set.UpdatedAt),
		IsFavourite:  set.IsFavourite,
		LikeBy:       set.LikeBy,
		IsPublic:     set.IsPublic,
		Image:        set.Image,
	}
}

func ToPBSetProtoDishes(dishes []SetProtoDish) []*order.SetProtoDish {
	pbDishes := make([]*order.SetProtoDish, len(dishes))
	for i, dish := range dishes {
		pbDishes[i] = &order.SetProtoDish{
			Id:    dish.ID,
			Name:  dish.Name,
			Price: dish.Price,
		}
	}
	return pbDishes
}

func ToDishOrderItems(pbItems []*order.DishOrderItem) []DishOrderItem {
	items := make([]DishOrderItem, len(pbItems))
	for i, pbItem := range pbItems {
		items[i] = DishOrderItem{
			ID:       pbItem.Id,
			Quantity: pbItem.Quantity,
			Dish:     ToDish(pbItem.Dish),
		}
	}
	return items
}

func ToSetOrderItems(pbItems []*order.SetOrderItem) []SetOrderItem {
	items := make([]SetOrderItem, len(pbItems))
	for i, pbItem := range pbItems {
		items[i] = SetOrderItem{
			ID:             pbItem.Id,
			Quantity:       pbItem.Quantity,
			Set:            ToSetProto(pbItem.Set),
			ModifiedDishes: ToSetProtoDishes(pbItem.ModifiedDishes),
		}
	}
	return items
}

func ToDish(pbDish *order.DishOrder) Dish {
	return Dish{
		ID:          pbDish.Id,
		Name:        pbDish.Name,
		Price:       pbDish.Price,
		Description: pbDish.Description,
		Image:       pbDish.Image,
		Status:      pbDish.Status,
		CreatedAt:   pbDish.CreatedAt.AsTime(),
		UpdatedAt:   pbDish.UpdatedAt.AsTime(),
	}
}

func ToSetProto(pbSet *order.SetProto) SetProto {
	return SetProto{
		ID:           int64(pbSet.Id),
		Name:         pbSet.Name,
		Description:  pbSet.Description,
		Dishes:       ToSetProtoDishes(pbSet.Dishes),
		UserID:       int32PtrToNullable(pbSet.UserId),
		CreatedAt:    pbSet.CreatedAt.AsTime(),
		UpdatedAt:    pbSet.UpdatedAt.AsTime(),
		IsFavourite:  pbSet.IsFavourite,
		LikeBy:       pbSet.LikeBy,
		IsPublic:     pbSet.IsPublic,
		Image:        pbSet.Image,
	}
}

func ToSetProtoDishes(pbDishes []*order.SetProtoDish) []SetProtoDish {
	dishes := make([]SetProtoDish, len(pbDishes))
	for i, pbDish := range pbDishes {
		dishes[i] = SetProtoDish{
			ID:    pbDish.Id,
			Name:  pbDish.Name,
			Price: pbDish.Price,
		}
	}
	return dishes
}

func int32Ptr(v *int32) *int32 {
	if v == nil {
		return nil
	}
	return v
}

func int32PtrToNullable(v *int32) *int32 {
	if v == nil {
		return nil
	}
	return v
}