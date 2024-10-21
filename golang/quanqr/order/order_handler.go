package order_grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"english-ai-full/logger"
	"english-ai-full/quanqr/proto_qr/order"
	"english-ai-full/token"

	"github.com/go-chi/chi"
	"google.golang.org/protobuf/types/known/timestamppb"
)
type OrderHandlerController struct {
    ctx        context.Context
    client     order.OrderServiceClient
    TokenMaker *token.JWTMaker
    logger     *logger.Logger
}

func NewOrderHandler(client order.OrderServiceClient, secretKey string) *OrderHandlerController {
    return &OrderHandlerController{
        ctx:        context.Background(),
        client:     client,
        TokenMaker: token.NewJWTMaker(secretKey),
        logger:     logger.NewLogger(),
    }
}

func (h *OrderHandlerController) CreateOrder(w http.ResponseWriter, r *http.Request) {
    var orderReq CreateOrderRequestType

    if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
        http.Error(w, "error decoding request body", http.StatusBadRequest)
        return
    }

    h.logger.Info(fmt.Sprintf("Creating new order: %+v", orderReq))
    createdOrderResponse, err := h.client.CreateOrder(h.ctx, ToPBCreateOrderRequest(orderReq))
    if err != nil {
        h.logger.Error("Error creating order: " + err.Error())
        http.Error(w, "error creating order", http.StatusInternalServerError)
        return
    }

    res := ToOrderResFromPbOrderResponse(createdOrderResponse)
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
    ordersResponse, err := h.client.GetOrders(h.ctx, ToPBGetOrdersRequest(req))
    if err != nil {
        h.logger.Error("Error fetching orders: " + err.Error())
        http.Error(w, "error getting orders", http.StatusInternalServerError)
        return
    }

    res := ToOrderListResFromPbOrderListResponse(ordersResponse)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(res)
}

func (h *OrderHandlerController) GetOrderDetail(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    i, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        http.Error(w, "error parsing ID", http.StatusBadRequest)
        return
    }

    h.logger.Info(fmt.Sprintf("Fetching order detail for ID: %d", i))
    orderResponse, err := h.client.GetOrderDetail(h.ctx, &order.OrderIdParam{Id: i})
    if err != nil {
        h.logger.Error("Error fetching order detail: " + err.Error())
        http.Error(w, "error getting order", http.StatusInternalServerError)
        return
    }

    res := ToOrderResFromPbOrderResponse(orderResponse)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(res)
}

func (h *OrderHandlerController) UpdateOrder(w http.ResponseWriter, r *http.Request) {
    var orderReq UpdateOrderRequestType
    if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
        http.Error(w, "error decoding request body", http.StatusBadRequest)
        return
    }

    h.logger.Info(fmt.Sprintf("Updating order: %d", orderReq.ID))
    updatedOrderResponse, err := h.client.UpdateOrder(h.ctx, ToPBUpdateOrderRequest(orderReq))
    if err != nil {
        h.logger.Error("Error updating order: " + err.Error())
        http.Error(w, "error updating order", http.StatusInternalServerError)
        return
    }

    res := ToOrderResFromPbOrderResponse(updatedOrderResponse)
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
    paidOrdersResponse, err := h.client.PayOrders(h.ctx, ToPBPayOrdersRequest(req))
    if err != nil {
        h.logger.Error("Error processing payment: " + err.Error())
        http.Error(w, "error processing payment", http.StatusInternalServerError)
        return
    }

    res := ToOrderListResFromPbOrderListResponse(paidOrdersResponse)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(res)
}

// Conversion functions
func ToPBCreateOrderRequest(req CreateOrderRequestType) *order.CreateOrderRequest {
    return &order.CreateOrderRequest{
        GuestId:        req.GuestID,
        UserId:         req.UserID,
        IsGuest:        req.IsGuest,
        TableNumber:    req.TableNumber,
        OrderHandlerId: req.OrderHandlerID,
        Status:         req.Status,
        CreatedAt:      timestamppb.New(req.CreatedAt),
        UpdatedAt:      timestamppb.New(req.UpdatedAt),
        TotalPrice:     req.TotalPrice,
        DishItems:      ToPBDishOrderItems(req.DishItems),
        SetItems:       ToPBSetOrderItems(req.SetItems),
        BowChili:       req.BowChili,
        BowNoChili:     req.BowNoChili,
    }
}


func ToPBUpdateOrderRequest(req UpdateOrderRequestType) *order.UpdateOrderRequest {
    return &order.UpdateOrderRequest{
        Id:             req.ID,
        GuestId:        req.GuestID,
        UserId:         req.UserID,
        TableNumber:    req.TableNumber,
        OrderHandlerId: req.OrderHandlerID,
        Status:         req.Status,
        TotalPrice:     req.TotalPrice,
        DishItems:      ToPBDishOrderItems(req.DishItems),
        SetItems:       ToPBSetOrderItems(req.SetItems),
        IsGuest:        req.IsGuest,
        BowChili:       req.BowChili,
        BowNoChili:     req.BowNoChili,
    }
}


func ToPBGetOrdersRequest(req GetOrdersRequestType) *order.GetOrdersRequest {
    return &order.GetOrdersRequest{
        FromDate: timestamppb.New(req.FromDate),
        ToDate:   timestamppb.New(req.ToDate),
        UserId:   req.UserID,
        GuestId:  req.GuestID,
    }
}

func ToPBPayOrdersRequest(req PayOrdersRequestType) *order.PayOrdersRequest {
    pbReq := &order.PayOrdersRequest{}
    if req.GuestID != nil {
        pbReq.Identifier = &order.PayOrdersRequest_GuestId{GuestId: *req.GuestID}
    } else if req.UserID != nil {
        pbReq.Identifier = &order.PayOrdersRequest_UserId{UserId: *req.UserID}
    }
    return pbReq
}

func ToPBDishOrderItems(items []DishOrderItem) []*order.DishOrderItem {
    pbItems := make([]*order.DishOrderItem, len(items))
    for i, item := range items {
        pbItems[i] = &order.DishOrderItem{
            Id:       item.ID,
            Quantity: item.Quantity,
        }
    }
    return pbItems
}

func ToPBSetOrderItems(items []SetOrderItem) []*order.SetOrderItem {
    pbItems := make([]*order.SetOrderItem, len(items))
    for i, item := range items {
        pbItems[i] = &order.SetOrderItem{
            Id:       item.ID,
            Quantity: item.Quantity,
        }
    }
    return pbItems
}






func ToOrderResFromPbOrderResponse(pbRes *order.OrderResponse) OrderResponse {
    return OrderResponse{
        Data: ToOrderFromPbOrder(pbRes.Data),
    }
}

func ToOrderListResFromPbOrderListResponse(pbRes *order.OrderListResponse) OrderListResponse {
    orders := make([]OrderType, len(pbRes.Data))
    for i, pbOrder := range pbRes.Data {
        orders[i] = ToOrderFromPbOrder(pbOrder)
    }
    return OrderListResponse{
        Data: orders,
    }
}

func ToOrderFromPbOrder(pbOrder *order.Order) OrderType {
    return OrderType{
        ID:             pbOrder.Id,
        GuestID:        pbOrder.GuestId,
        UserID:         pbOrder.UserId,
        IsGuest:        pbOrder.IsGuest,
        TableNumber:    pbOrder.TableNumber,
        OrderHandlerID: pbOrder.OrderHandlerId,
        Status:         pbOrder.Status,
        CreatedAt:      pbOrder.CreatedAt.AsTime(),
        UpdatedAt:      pbOrder.UpdatedAt.AsTime(),
        TotalPrice:     pbOrder.TotalPrice,
        DishItems:      ToDishOrderItemsFromPb(pbOrder.DishItems),
        SetItems:       ToSetOrderItemsFromPb(pbOrder.SetItems),
        BowChili:       pbOrder.BowChili,
        BowNoChili:     pbOrder.BowNoChili,
    }
}


func ToDishOrderItemsFromPb(pbItems []*order.DishOrderItem) []DishOrderItem {
    items := make([]DishOrderItem, len(pbItems))
    for i, pbItem := range pbItems {
        items[i] = DishOrderItem{
            ID:       pbItem.Id,
            Quantity: pbItem.Quantity,
        }
    }
    return items
}

func ToSetOrderItemsFromPb(pbItems []*order.SetOrderItem) []SetOrderItem {
    items := make([]SetOrderItem, len(pbItems))
    for i, pbItem := range pbItems {
        items[i] = SetOrderItem{
            ID:       pbItem.Id,
            Quantity: pbItem.Quantity,
        }
    }
    return items
}




