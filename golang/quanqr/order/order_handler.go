package order_grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"english-ai-full/logger"
	"english-ai-full/quanqr/proto_qr/order"
	"english-ai-full/token"

	"github.com/go-chi/chi"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)


type OrderHandlerController struct {
    ctx        context.Context
    client     order.OrderServiceClient
    TokenMaker *token.JWTMaker
    logger     *logger.Logger
    NewLoggerType     *logger.NewLoggerType  
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
    // Parse and validate the request body
    var orderReq CreateOrderRequestType
    if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
        h.logger.Error("Error decoding request body: " + err.Error())
        http.Error(w, "error decoding request body", http.StatusBadRequest)
        return
    }

    // Set default values for new fields if they're not provided
    if orderReq.Version == 0 {
        orderReq.Version = 1 // Initial version for new orders
    }

    // Validate the request
    if err := validateCreateOrderRequest(orderReq); err != nil {
        h.logger.Error("Validation error: " + err.Error())
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Add modification tracking information for dish items
    for i := range orderReq.DishItems {
        orderReq.DishItems[i].ModificationType = "INITIAL"
        orderReq.DishItems[i].ModificationNumber = 1
        orderReq.DishItems[i].OrderName = orderReq.OrderName
    }

    // Add modification tracking information for set items
    for i := range orderReq.SetItems {
        orderReq.SetItems[i].ModificationType = "INITIAL"
        orderReq.SetItems[i].ModificationNumber = 1
        orderReq.SetItems[i].OrderName = orderReq.OrderName
    }

    // Convert request to protobuf format
    pbReq := ToPBCreateOrderRequest(orderReq)

    // Call the service to create the order
    h.logger.Info(fmt.Sprintf("Creating order with name: %s", orderReq.OrderName))
    createdOrderResponse, err := h.client.CreateOrder(h.ctx, pbReq)
    if err != nil {
        h.logger.Error("Error creating order: " + err.Error())
        http.Error(w, "error creating order: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Convert the response and send it back
    res := ToOrderResFromPbOrderResponse(createdOrderResponse)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(res); err != nil {
        h.logger.Error("Error encoding response: " + err.Error())
        http.Error(w, "error encoding response", http.StatusInternalServerError)
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

func (h *OrderHandlerController) GetOrders(w http.ResponseWriter, r *http.Request) {
    // Only accept GET requests
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse query parameters
    query := r.URL.Query()
    
    // Get page parameter with default value 1
    page := int32(1)
    if pageStr := query.Get("page"); pageStr != "" {
        if pageInt, err := strconv.ParseInt(pageStr, 10, 32); err == nil {
            page = int32(pageInt)
        }
    }

    // Get page_size parameter with default value 10
    pageSize := int32(10)
    if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
        if pageSizeInt, err := strconv.ParseInt(pageSizeStr, 10, 32); err == nil {
            pageSize = int32(pageSizeInt)
        }
    }

    // Validate pagination parameters
    if page < 1 {
        page = 1
    }
    if pageSize < 1 {
        pageSize = 10
    }

    h.logger.Info("Fetching orders list")
    ordersResponse, err := h.client.GetOrders(h.ctx, &order.GetOrdersRequest{
        Page:     page,
        PageSize: pageSize,
    })


    fmt.Printf("golang/quanqr/order/order_handler.go ordersResponse %v\n", ordersResponse)
    if err != nil {
        h.logger.Error("Error fetching orders list: " + err.Error())
        http.Error(w, "failed to fetch orders: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Convert protobuf response to HTTP response
    res := ToOrderListResFromPbOrderListResponse(ordersResponse)

    // Set response headers and encode response
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(res); err != nil {
        h.logger.Error("Error encoding response: " + err.Error())
        http.Error(w, "error encoding response", http.StatusInternalServerError)
        return
    }
}
func (h *OrderHandlerController) UpdateOrder(w http.ResponseWriter, r *http.Request) {
    var orderReq UpdateOrderRequestType
    if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
        http.Error(w, "error decoding request body", http.StatusBadRequest)
        return
    }

    updatedOrderResponse, err := h.client.UpdateOrder(h.ctx, ToPBUpdateOrderRequest(orderReq))
    if err != nil {
        h.logger.Error("Error updating order: " + err.Error())
        http.Error(w, "error updating order", http.StatusInternalServerError)
        return
    }

    res := ToOrderDetailedListResponseFromPbResponse(updatedOrderResponse)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(res)
}

// func (h *OrderHandlerController) AddingSetsDishesOrder(w http.ResponseWriter, r *http.Request) {
//     var orderReq UpdateOrderRequestType
//     if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
//         http.Error(w, "error decoding request body", http.StatusBadRequest)
//         return
//     }

//     // Fixed: Remove extra parentheses and pass the correct parameters
//     updatedOrderResponse, err := h.client.AddingSetsDishesOrder(h.ctx, ToPBUpdateOrderRequest(orderReq))
//     if err != nil {
//         h.logger.Error("Error updating order: " + err.Error())
//         http.Error(w, "error updating order", http.StatusInternalServerError)
//         return
//     }

//     res := ToOrderDetailedListResponseFromPbResponse(updatedOrderResponse)
//     w.Header().Set("Content-Type", "application/json")
//     w.WriteHeader(http.StatusOK)
//     json.NewEncoder(w).Encode(res)
// }
func (h *OrderHandlerController) PayOrders(w http.ResponseWriter, r *http.Request) {
    var req PayOrdersRequestType
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "error decoding request body", http.StatusBadRequest)
        return
    }

    h.logger.Info("Processing order payment")
    paymentResponse, err := h.client.PayOrders(h.ctx, ToPBPayOrdersRequest(req))
    if err != nil {
        h.logger.Error("Error processing payment: " + err.Error())
        http.Error(w, "error processing payment", http.StatusInternalServerError)
        return
    }

    res := ToOrderListResFromPbOrderListResponse(paymentResponse)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(res)
}

func (h *OrderHandlerController) GetOrderProtoListDetail(w http.ResponseWriter, r *http.Request) {
    h.logger.Info("Fetching detailed order list golang/quanqr/order/order_handler.go 1")

    // Parse query parameters for pagination
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1 // Default to first page if invalid
    }
    h.logger.Info("Fetching detailed order list golang/quanqr/order/order_handler.go 2")
    pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
    if err != nil || pageSize < 1 {
        pageSize = 10 // Default page size if invalid
    }
    h.logger.Info("Fetching detailed order list golang/quanqr/order/order_handler.go 3")
    // Create the request with pagination parameters
    req := &order.GetOrdersRequest{
        Page:     int32(page),
        PageSize: int32(pageSize),
    }
    h.logger.Info("Fetching detailed order list golang/quanqr/order/order_handler.go 4")
    // Call the service
    ordersResponse, err := h.client.GetOrderProtoListDetail(h.ctx, req)
    if err != nil {
        h.logger.Error("Error fetching detailed order list: " + err.Error())
        http.Error(w, "failed to fetch detailed orders: "+err.Error(), http.StatusInternalServerError)
        return
    }
    h.logger.Info("Fetching detailed order list golang/quanqr/order/order_handler.go 5")
    fmt.Printf("golang/quanqr/order/order_handler.go GetOrderProtoListDetail res %v\n", ordersResponse)
    // Convert the response
    res := ToOrderDetailedListResponseFromProto(ordersResponse)
 
    // Send response

    h.logger.Info("Fetching detailed order list golang/quanqr/order/order_handler.go 6")
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(res); err != nil {
        h.logger.Error("Error encoding response: " + err.Error())
        http.Error(w, "error encoding response", http.StatusInternalServerError)
        return
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






func ToOrderResFromPbOrderResponse(pbRes *order.OrderResponse) OrderResponse {
    return OrderResponse{
        Data: ToOrderFromPbOrder(pbRes.Data),
    }
}

func ToOrderListResFromPbOrderListResponse(pbRes *order.OrderListResponse) *OrderListResponse {
    if pbRes == nil {
        return nil
    }

    // Initialize response with proper capacity
    orders := make([]OrderType, 0, len(pbRes.Data))
    
    // Convert each order
    for _, pbOrder := range pbRes.Data {
        if pbOrder != nil {
            orders = append(orders, ToOrderFromPbOrder(pbOrder))
        }
    }

    return &OrderListResponse{
        Data: orders,
        Pagination: PaginationInfo{
            CurrentPage: pbRes.GetPagination().GetCurrentPage(),
            TotalPages: pbRes.GetPagination().GetTotalPages(),
            TotalItems: pbRes.GetPagination().GetTotalItems(),
            PageSize:   pbRes.GetPagination().GetPageSize(),
        },
    }
}


func ToOrderDishesFromPbDishOrderItems(pbItems []*order.DishOrderItem) []OrderDish {
    if pbItems == nil {
        return nil
    }

    items := make([]OrderDish, 0, len(pbItems))
    for _, pbItem := range pbItems {
        if pbItem != nil {
            items = append(items, OrderDish{
                DishID:   pbItem.GetDishId(),
                Quantity: pbItem.GetQuantity(),
            })
        }
    }
    return items
}

func ToOrderSetsFromPbSetOrderItems(pbItems []*order.SetOrderItem) []OrderSet {
    if pbItems == nil {
        return nil
    }

    items := make([]OrderSet, 0, len(pbItems))
    for _, pbItem := range pbItems {
        if pbItem != nil {
            items = append(items, OrderSet{
                SetID:    pbItem.GetSetId(),
                Quantity: pbItem.GetQuantity(),
            })
        }
    }
    return items
}



func ToOrderDetailedDishesFromPbOrderDetailedDishes(pbDishes []*order.OrderDetailedDish) []OrderDetailedDish {
    dishes := make([]OrderDetailedDish, len(pbDishes))
    for i, pbDish := range pbDishes {
        dishes[i] = OrderDetailedDish{
            DishID:      pbDish.DishId,
            Quantity:    pbDish.Quantity,
            Name:        pbDish.Name,
            Price:       pbDish.Price,
            Description: pbDish.Description,
            Image:       pbDish.Image,
            Status:      pbDish.Status,
        }
    }
    return dishes
}



func ToOrderDetailedDishFromPbOrderDetailedDish(pbDish *order.OrderDetailedDish) OrderDetailedDish {
    if pbDish == nil {
        return OrderDetailedDish{}
    }

    return OrderDetailedDish{
        DishID:      pbDish.DishId,
        Quantity:    pbDish.Quantity,
        Name:        pbDish.Name,
        Price:       pbDish.Price,
        Description: pbDish.Description,
        Image:       pbDish.Image,
        Status:      pbDish.Status,
    }
}


func ToOrderSetDetailedFromPbOrderSetDetailed(pbSet *order.OrderSetDetailed) OrderSetDetailed {
    if pbSet == nil {
        return OrderSetDetailed{}
    }

    dishes := make([]OrderDetailedDish, len(pbSet.Dishes))
    for i, pbDish := range pbSet.Dishes {
        dishes[i] = ToOrderDetailedDishFromPbOrderDetailedDish(pbDish)
    }

    var createdAt, updatedAt time.Time
    if pbSet.CreatedAt != nil {
        createdAt = pbSet.CreatedAt.AsTime()
    }
    if pbSet.UpdatedAt != nil {
        updatedAt = pbSet.UpdatedAt.AsTime()
    }

    return OrderSetDetailed{
        ID:          pbSet.Id,
        Name:        pbSet.Name,
        Description: pbSet.Description,
        Dishes:      dishes,
        UserID:      pbSet.UserId,
        CreatedAt:   createdAt,
        UpdatedAt:   updatedAt,
        IsFavourite: pbSet.IsFavourite,
        LikeBy:      pbSet.LikeBy,
        IsPublic:    pbSet.IsPublic,
        Image:       pbSet.Image,
        Price:       pbSet.Price,
        Quantity:    pbSet.Quantity,
    }
}


// -------------




// Helper functions for conversion
func ToOrderSetsDetailedFromProto(pbSets []*order.OrderSetDetailed) []OrderSetDetailed {
    if pbSets == nil {
        return nil
    }

    sets := make([]OrderSetDetailed, len(pbSets))
    for i, pbSet := range pbSets {
        sets[i] = OrderSetDetailed{
            ID:          pbSet.Id,
            Name:        pbSet.Name,
            Description: pbSet.Description,
            Dishes:      ToOrderDetailedDishesFromProto(pbSet.Dishes),
            UserID:      pbSet.UserId,
            CreatedAt:   pbSet.CreatedAt.AsTime(),
            UpdatedAt:   pbSet.UpdatedAt.AsTime(),
            IsFavourite: pbSet.IsFavourite,
            LikeBy:      pbSet.LikeBy,
            IsPublic:    pbSet.IsPublic,
            Image:       pbSet.Image,
            Price:       pbSet.Price,
            Quantity:    pbSet.Quantity,
        }
    }
    return sets
}

func ToOrderDetailedDishesFromProto(pbDishes []*order.OrderDetailedDish) []OrderDetailedDish {
    if pbDishes == nil {
        return nil
    }

    dishes := make([]OrderDetailedDish, len(pbDishes))
    for i, pbDish := range pbDishes {
        dishes[i] = OrderDetailedDish{
            DishID:      pbDish.DishId,
            Quantity:    pbDish.Quantity,
            Name:        pbDish.Name,
            Price:       pbDish.Price,
            Description: pbDish.Description,
            Image:       pbDish.Image,
            Status:      pbDish.Status,
        }
    }
    return dishes
}





// -------------------

// Conversion functions


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
        Topping:       req.Topping,
        TrackingOrder:     req.TrackingOrder,
        TakeAway:       req.TakeAway,
        ChiliNumber:    req.ChiliNumber,
        TableToken:     req.TableToken,
        OrderName:      req.OrderName,  // Added new field
    }
}

func ToOrderFromPbOrder(pbOrder *order.Order) OrderType {
    if pbOrder == nil {
        return OrderType{}
    }

    var createdAt, updatedAt time.Time
    if pbOrder.CreatedAt != nil {
        createdAt = pbOrder.CreatedAt.AsTime()
    }
    if pbOrder.UpdatedAt != nil {
        updatedAt = pbOrder.UpdatedAt.AsTime()
    }

    return OrderType{
        ID:             pbOrder.GetId(),
        GuestID:        pbOrder.GetGuestId(),
        UserID:         pbOrder.GetUserId(),
        IsGuest:        pbOrder.GetIsGuest(),
        TableNumber:    pbOrder.GetTableNumber(),
        OrderHandlerID: pbOrder.GetOrderHandlerId(),
        Status:         pbOrder.GetStatus(),
        CreatedAt:      createdAt,
        UpdatedAt:      updatedAt,
        TotalPrice:     pbOrder.GetTotalPrice(),
        DishItems:      ToOrderDishesFromPbDishOrderItems(pbOrder.DishItems),
        SetItems:       ToOrderSetsFromPbSetOrderItems(pbOrder.SetItems),
        Topping:       pbOrder.GetTopping(),
        TrackingOrder:     pbOrder.GetTrackingOrder(),
        TakeAway:       pbOrder.GetTakeAway(),
        ChiliNumber:    pbOrder.GetChiliNumber(),
        TableToken:     pbOrder.GetTableToken(),
        OrderName:      pbOrder.GetOrderName(),  // Added new field
    }
}



func (h *OrderHandlerController) CreateOrder2(w http.ResponseWriter, r *http.Request) {
    var orderReq CreateOrderRequestType
    
    // Read the entire body first
    body, err := io.ReadAll(r.Body)
    if err != nil {
        h.logger.Error("Error reading request body: " + err.Error())
        http.Error(w, "error reading request body", http.StatusBadRequest)
        return
    }
    
    // Log the raw body for debugging
    h.logger.Info(fmt.Sprintf("Raw request body: %s", string(body)))

    // Decode the JSON
    if err := json.Unmarshal(body, &orderReq); err != nil {
        h.logger.Error("Error decoding request body: " + err.Error())
        http.Error(w, "error decoding request body", http.StatusBadRequest)
        return
    }

    h.logger.Info(fmt.Sprintf("golang/quanqr/order/order_handler.go Decoded order request: %+v", orderReq))
    
    pbReq := ToPBCreateOrderRequest(orderReq)
    createdOrderResponse, err := h.client.CreateOrder(h.ctx, pbReq)
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



func (h *OrderHandlerController) FetchOrdersByCriteria(w http.ResponseWriter, r *http.Request) {
    h.logger.Info("Fetching orders by criteria")

    // Parse query parameters
    query := r.URL.Query()
    
    // Get page parameter with default value 1
    page := int32(1)
    if pageStr := query.Get("page"); pageStr != "" {
        if pageInt, err := strconv.ParseInt(pageStr, 10, 32); err == nil {
            page = int32(pageInt)
        }
    }

    // Get page_size parameter with default value 10
    pageSize := int32(10)
    if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
        if pageSizeInt, err := strconv.ParseInt(pageSizeStr, 10, 32); err == nil {
            pageSize = int32(pageSizeInt)
        }
    }

    var req FetchOrdersByCriteriaRequestType
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        if err != io.EOF {
            h.logger.Error("Error decoding request body: " + err.Error())
            http.Error(w, "error decoding request body", http.StatusBadRequest)
            return
        }
    }

    // Override page and pageSize from query parameters if they exist
    req.Page = page
    req.PageSize = pageSize

    // Convert request to protobuf
    pbReq := &order.FetchOrdersByCriteriaRequest{
        OrderIds:  req.OrderIds,
        OrderName: req.OrderName,
        Page:      req.Page,
        PageSize:  req.PageSize,
    }

    // Add date filters if present
    if req.StartDate != nil {
        pbReq.StartDate = timestamppb.New(*req.StartDate)
    }
    if req.EndDate != nil {
        pbReq.EndDate = timestamppb.New(*req.EndDate)
    }

    // Call the service
    response, err := h.client.FetchOrdersByCriteria(h.ctx, pbReq)
    if err != nil {
        h.logger.Error("Error fetching orders by criteria: " + err.Error())
        http.Error(w, "failed to fetch orders: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Convert and send response
    res := ToOrderDetailedListResponseFromProto(response)
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(res); err != nil {
        h.logger.Error("Error encoding response: " + err.Error())
        http.Error(w, "error encoding response", http.StatusInternalServerError)
        return
    }
}

// validateCreateOrderRequest validates the create order request
func validateCreateOrderRequest(req CreateOrderRequestType) error {
    if req.OrderName == "" {
        return fmt.Errorf("order_name is required")
    }

    if req.IsGuest && req.GuestID == 0 {
        return fmt.Errorf("guest_id is required for guest orders")
    }

    if !req.IsGuest && req.UserID == 0 {
        return fmt.Errorf("user_id is required for user orders")
    }

    if req.TableNumber == 0 {
        return fmt.Errorf("table_number is required")
    }

    if len(req.DishItems) == 0 && len(req.SetItems) == 0 {
        return fmt.Errorf("at least one dish or set item is required")
    }

    return nil
}



func ToOrderDetailedResponsesFromPbResponses(pbResponses []*order.OrderDetailedResponse) []OrderDetailedResponse {
    // Handle nil input case to avoid panic
    if pbResponses == nil {
        return nil
    }

    // Create a slice with the same length as input
    responses := make([]OrderDetailedResponse, len(pbResponses))

    // Convert each protobuf response to domain model
    for i, pbRes := range pbResponses {
        // Skip nil entries to avoid panic
        if pbRes == nil {
            continue
        }

        // Convert each protobuf response to our domain model
        responses[i] = OrderDetailedResponse{
            ID:             pbRes.Id,
            GuestID:        pbRes.GuestId,
            UserID:         pbRes.UserId,
            TableNumber:    pbRes.TableNumber,
            OrderHandlerID: pbRes.OrderHandlerId,
            Status:         pbRes.Status,
            TotalPrice:     pbRes.TotalPrice,
            IsGuest:        pbRes.IsGuest,
            Topping:       pbRes.Topping,
            TrackingOrder: pbRes.TrackingOrder,
            TakeAway:      pbRes.TakeAway,
            ChiliNumber:   pbRes.ChiliNumber,
            TableToken:    pbRes.TableToken,
            OrderName:     pbRes.OrderName,
            // Convert nested data structures using existing helper functions
            DataSet:       ToOrderSetsDetailedFromProto(pbRes.DataSet),
            DataDish:      ToOrderDetailedDishesFromProto(pbRes.DataDish),
        }
    }

    return responses
}

// new --------------------------





func (h *OrderHandlerController) CreateOrder3(w http.ResponseWriter, r *http.Request) {
    // Parse and validate the request body
    var orderReq CreateOrderRequestType
    if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
        h.logger.Error("Error decoding request body: " + err.Error())
        http.Error(w, "error decoding request body", http.StatusBadRequest)
        return
    }

    // Set default values for new fields
    if orderReq.Version == 0 {
        orderReq.Version = 1 // Initial version for new orders
    }

    // Validate the request
    if err := validateCreateOrderRequest(orderReq); err != nil {
        h.logger.Error("Validation error: " + err.Error())
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Add modification tracking information for dish items
    for i := range orderReq.DishItems {
        orderReq.DishItems[i].ModificationType = "INITIAL"
        orderReq.DishItems[i].ModificationNumber = 1
        orderReq.DishItems[i].OrderName = orderReq.OrderName
        // Set timestamps for modification tracking
        orderReq.DishItems[i].CreatedAt = time.Now()
        orderReq.DishItems[i].UpdatedAt = time.Now()
    }

    // Add modification tracking information for set items
    for i := range orderReq.SetItems {
        orderReq.SetItems[i].ModificationType = "INITIAL"
        orderReq.SetItems[i].ModificationNumber = 1
        orderReq.SetItems[i].OrderName = orderReq.OrderName
        // Set timestamps for modification tracking
        orderReq.SetItems[i].CreatedAt = time.Now()
        orderReq.SetItems[i].UpdatedAt = time.Now()
    }

    // Convert request to protobuf format
    pbReq := ToPBCreateOrderRequest(orderReq)

    // Call the service to create the order
    h.logger.Info(fmt.Sprintf("Creating order with name: %s", orderReq.OrderName))
    createdOrderResponse, err := h.client.CreateOrder(h.ctx, pbReq)
    if err != nil {
        h.logger.Error("Error creating order: " + err.Error())
        http.Error(w, "error creating order: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Convert the response and send it back
    res := ToOrderResFromPbOrderResponse(createdOrderResponse)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(res); err != nil {
        h.logger.Error("Error encoding response: " + err.Error())
        http.Error(w, "error encoding response", http.StatusInternalServerError)
        return
    }
}

// Updated conversion functions to handle new fields and structures
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
        Topping:        req.Topping,
        TrackingOrder:  req.TrackingOrder,
        TakeAway:       req.TakeAway,
        ChiliNumber:    req.ChiliNumber,
        TableToken:     req.TableToken,
        OrderName:      req.OrderName,
        Version:        req.Version,
        ParentOrderId:  req.ParentOrderID,
    }
}

func ToPBDishOrderItems(items []OrderDish) []*order.DishOrderItem {
    pbItems := make([]*order.DishOrderItem, len(items))
    for i, item := range items {
        pbItems[i] = &order.DishOrderItem{
            Id:                item.ID,
            DishId:           item.DishID,
            Quantity:         item.Quantity,
            CreatedAt:        timestamppb.New(item.CreatedAt),
            UpdatedAt:        timestamppb.New(item.UpdatedAt),
            OrderName:        item.OrderName,
            ModificationType: item.ModificationType,
            ModificationNumber: item.ModificationNumber,
        }
    }
    return pbItems
}

func ToPBSetOrderItems(items []OrderSet) []*order.SetOrderItem {
    pbItems := make([]*order.SetOrderItem, len(items))
    for i, item := range items {
        pbItems[i] = &order.SetOrderItem{
            Id:                item.ID,
            SetId:            item.SetID,
            Quantity:         item.Quantity,
            CreatedAt:        timestamppb.New(item.CreatedAt),
            UpdatedAt:        timestamppb.New(item.UpdatedAt),
            OrderName:        item.OrderName,
            ModificationType: item.ModificationType,
            ModificationNumber: item.ModificationNumber,
        }
    }
    return pbItems
}



// The rest of the file remains largely unchanged...
func (h *OrderHandlerController) AddingSetsDishesOrder(w http.ResponseWriter, r *http.Request) {
    h.logger.Info("[Order Handler.AddingSetsDishesOrder] Starting to process adding sets/dishes request")

    var orderReq UpdateOrderRequestType
    if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
        h.logger.Error(fmt.Sprintf("[Order Handler.AddingSetsDishesOrder] Error decoding request body: %v", err))
        http.Error(w, "error decoding request body", http.StatusBadRequest)
        return
    }

    h.logger.Info(fmt.Sprintf("[Order Handler.AddingSetsDishesOrder] Processing request for Order ID: %d, Table: %d, Order Name: %s",
        orderReq.ID, orderReq.TableNumber, orderReq.OrderName))

    h.logger.Info(fmt.Sprintf("[Order Handler.AddingSetsDishesOrder] Request contains %d dish items and %d set items",
        len(orderReq.DishItems), len(orderReq.SetItems)))

    for _, dish := range orderReq.DishItems {
        h.logger.Info(fmt.Sprintf("[Order Handler.AddingSetsDishesOrder] Dish details - ID: %d, Quantity: %d, Order Name: %s",
            dish.DishID, dish.Quantity, dish.OrderName))
    }

    for _, set := range orderReq.SetItems {
        h.logger.Info(fmt.Sprintf("[Order Handler.AddingSetsDishesOrder] Set details - ID: %d, Quantity: %d, Order Name: %s",
            set.SetID, set.Quantity, set.OrderName))
    }

    pbReq := ToPBUpdateOrderRequest(orderReq)
    h.logger.Info("[Order Handler.AddingSetsDishesOrder] Converting request to protobuf format completed")

    updatedOrderResponse, err := h.client.AddingSetsDishesOrder(h.ctx, pbReq)
    if err != nil {
        h.logger.Error(fmt.Sprintf("[Order Handler.AddingSetsDishesOrder] Error updating order: %v", err))
        http.Error(w, "error updating order", http.StatusInternalServerError)
        return
    }

    h.logger.Info(fmt.Sprintf("[Order Handler.AddingSetsDishesOrder] Proto Response: \nPagination: %+v", 
        updatedOrderResponse.GetPagination()))

    for i, order := range updatedOrderResponse.GetData() {
        h.logger.Info(fmt.Sprintf("[Order Handler.AddingSetsDishesOrder] Proto Order %d: \n"+
            "ID: %d\n"+
            "OrderName: %s\n"+
            "TableNumber: %d\n"+
            "Status: %s\n"+
            "TotalPrice: %d\n"+
            "GuestID: %d\n"+
            "UserID: %d\n"+
            "OrderHandlerID: %d\n"+
            "IsGuest: %v\n"+
            "Topping: %s\n"+
            "TrackingOrder: %s\n"+
            "TakeAway: %v\n"+
            "ChiliNumber: %d\n"+
            "TableToken: %s\n"+
            "CurrentVersion: %d\n"+
            "ParentOrderID: %d\n"+
            "Dish Items Count: %d\n"+
            "Set Items Count: %d",
            i+1,
            order.GetId(),
            order.GetOrderName(),
            order.GetTableNumber(),
            order.GetStatus(),
            order.GetTotalPrice(),
            order.GetGuestId(),
            order.GetUserId(),
            order.GetOrderHandlerId(),
            order.GetIsGuest(),
            order.GetTopping(),
            order.GetTrackingOrder(),
            order.GetTakeAway(),
            order.GetChiliNumber(),
            order.GetTableToken(),
            order.GetCurrentVersion(),
            order.GetParentOrderId(),
            len(order.GetDataDish()),
            len(order.GetDataSet())))
    }

    res := ToOrderDetailedListResponseFromPbResponse(updatedOrderResponse)

    h.logger.Info(fmt.Sprintf("[Order Handler.AddingSetsDishesOrder] Converted Response: \nPagination: %+v", 
        res.Pagination))

    for i, order := range res.Data {
        h.logger.Info(fmt.Sprintf("[Order Handler.AddingSetsDishesOrder] Converted Order %d: \n"+
            "ID: %d\n"+
            "OrderName: %s\n"+
            "TableNumber: %d\n"+
            "Status: %s\n"+
            "TotalPrice: %d\n"+
            "GuestID: %d\n"+
            "UserID: %d\n"+
            "OrderHandlerID: %d\n"+
            "IsGuest: %v\n"+
            "Topping: %s\n"+
            "TrackingOrder: %s\n"+
            "TakeAway: %v\n"+
            "ChiliNumber: %d\n"+
            "TableToken: %s\n"+
            "CurrentVersion: %d\n"+
            "ParentOrderID: %d\n"+
            "Version History Count: %d\n"+
            "Dish Items Count: %d\n"+
            "Set Items Count: %d",
            i+1,
            order.ID,
            order.OrderName,
            order.TableNumber,
            order.Status,
            order.TotalPrice,
            order.GuestID,
            order.UserID,
            order.OrderHandlerID,
            order.IsGuest,
            order.Topping,
            order.TrackingOrder,
            order.TakeAway,
            order.ChiliNumber,
            order.TableToken,
            order.CurrentVersion,
            order.ParentOrderID,
            len(order.VersionHistory),
            len(order.DataDish),
            len(order.DataSet)))

        h.logger.Info(fmt.Sprintf("[Order Handler.AddingSetsDishesOrder] Order %d Total Summary: %+v",
            i+1, order.TotalSummary))
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(res); err != nil {
        h.logger.Error(fmt.Sprintf("[Order Handler.AddingSetsDishesOrder] Error encoding response: %v", err))
        http.Error(w, "error encoding response", http.StatusInternalServerError)
        return
    }

    h.logger.Info("[Order Handler.AddingSetsDishesOrder] Request completed successfully")
}


func ToOrderDetailedListResponseFromPbResponse(pbRes *order.OrderDetailedListResponse) *OrderDetailedListResponse {
    if pbRes == nil {
        return nil
    }
    
    responses := make([]OrderDetailedResponse, len(pbRes.Data))
    for i, pbDetailedRes := range pbRes.Data {
        if pbDetailedRes == nil {
            continue
        }

        // Convert version history
        versionHistory := make([]OrderVersionSummary, len(pbDetailedRes.VersionHistory))
        for j, pbVersion := range pbDetailedRes.VersionHistory {
            changes := make([]OrderItemChange, len(pbVersion.Changes))
            for k, pbChange := range pbVersion.Changes {
                changes[k] = OrderItemChange{
                    ItemType:        pbChange.ItemType,
                    ItemID:          pbChange.ItemId,
                    ItemName:        pbChange.ItemName,
                    QuantityChanged: pbChange.QuantityChanged,
                    Price:           pbChange.Price,
                }
            }
            
            versionHistory[j] = OrderVersionSummary{
                VersionNumber:     pbVersion.VersionNumber,
                TotalDishesCount: pbVersion.TotalDishesCount,
                TotalSetsCount:   pbVersion.TotalSetsCount,
                VersionTotalPrice: pbVersion.VersionTotalPrice,
                ModificationType: pbVersion.ModificationType,
                ModifiedAt:      pbVersion.ModifiedAt.AsTime(),
                Changes:         changes,
            }
        }

        // Convert most ordered items
        var mostOrderedItems []OrderItemCount
        if pbDetailedRes.TotalSummary != nil {
            mostOrderedItems = make([]OrderItemCount, len(pbDetailedRes.TotalSummary.MostOrderedItems))
            for j, pbItem := range pbDetailedRes.TotalSummary.MostOrderedItems {
                mostOrderedItems[j] = OrderItemCount{
                    ItemType:      pbItem.ItemType,
                    ItemID:        pbItem.ItemId,
                    ItemName:      pbItem.ItemName,
                    TotalQuantity: pbItem.TotalQuantity,
                }
            }
        }

        // Convert total summary
        var totalSummary OrderTotalSummary
        if pbDetailedRes.TotalSummary != nil {
            totalSummary = OrderTotalSummary{
                TotalVersions:        pbDetailedRes.TotalSummary.TotalVersions,
                TotalDishesOrdered:  pbDetailedRes.TotalSummary.TotalDishesOrdered,
                TotalSetsOrdered:    pbDetailedRes.TotalSummary.TotalSetsOrdered,
                CumulativeTotalPrice: pbDetailedRes.TotalSummary.CumulativeTotalPrice,
                MostOrderedItems:     mostOrderedItems,
            }
        }

        responses[i] = OrderDetailedResponse{
            ID:             pbDetailedRes.Id,
            GuestID:        pbDetailedRes.GuestId,
            UserID:         pbDetailedRes.UserId,
            TableNumber:    pbDetailedRes.TableNumber,
            OrderHandlerID: pbDetailedRes.OrderHandlerId,
            Status:         pbDetailedRes.Status,
            TotalPrice:     pbDetailedRes.TotalPrice,
            IsGuest:        pbDetailedRes.IsGuest,
            Topping:        pbDetailedRes.Topping,
            TrackingOrder:  pbDetailedRes.TrackingOrder,
            TakeAway:       pbDetailedRes.TakeAway,
            ChiliNumber:    pbDetailedRes.ChiliNumber,
            TableToken:     pbDetailedRes.TableToken,
            OrderName:      pbDetailedRes.OrderName,
            CurrentVersion: pbDetailedRes.CurrentVersion,
            ParentOrderID:  pbDetailedRes.ParentOrderId,
            DataSet:        ToOrderSetsDetailedFromProto(pbDetailedRes.DataSet),
            DataDish:       ToOrderDetailedDishesFromProto(pbDetailedRes.DataDish),
            VersionHistory: versionHistory,
            TotalSummary:   totalSummary,
        }
    }

    return &OrderDetailedListResponse{
        Data: responses,
        Pagination: PaginationInfo{
            CurrentPage: pbRes.GetPagination().GetCurrentPage(),
            TotalPages: pbRes.GetPagination().GetTotalPages(),
            TotalItems: pbRes.GetPagination().GetTotalItems(),
            PageSize:   pbRes.GetPagination().GetPageSize(),
        },
    }
}

// new removing 



func (h *OrderHandlerController) RemovingSetsDishesOrder(w http.ResponseWriter, r *http.Request) {
    h.logger.Info("[Order Handler.RemovingSetsDishesOrder] Starting to process removal request")

    var orderReq UpdateOrderRequestType
    if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
        h.logger.Error(fmt.Sprintf("[Order Handler.RemovingSetsDishesOrder] Error decoding request body: %v", err))
        http.Error(w, "error decoding request body", http.StatusBadRequest)
        return
    }

    h.logger.Info(fmt.Sprintf("[Order Handler.RemovingSetsDishesOrder] Processing removal for Order ID: %d, Table: %d, Order Name: %s",
        orderReq.ID, orderReq.TableNumber, orderReq.OrderName))

    h.logger.Info(fmt.Sprintf("[Order Handler.RemovingSetsDishesOrder] Request contains %d dish removals and %d set removals",
        len(orderReq.DishItems), len(orderReq.SetItems)))

    for _, dish := range orderReq.DishItems {
        h.logger.Info(fmt.Sprintf("[Order Handler.RemovingSetsDishesOrder] Dish removal details - ID: %d, Quantity: %d, Order Name: %s",
            dish.DishID, dish.Quantity, dish.OrderName))
    }

    for _, set := range orderReq.SetItems {
        h.logger.Info(fmt.Sprintf("[Order Handler.RemovingSetsDishesOrder] Set removal details - ID: %d, Quantity: %d, Order Name: %s",
            set.SetID, set.Quantity, set.OrderName))
    }

    pbReq := ToPBUpdateOrderRequest(orderReq)
    h.logger.Info("[Order Handler.RemovingSetsDishesOrder] Converting request to protobuf format completed")

    updatedOrderResponse, err := h.client.RemovingSetsDishesOrder(h.ctx, pbReq)
    if err != nil {
        h.logger.Error(fmt.Sprintf("[Order Handler.RemovingSetsDishesOrder] Error removing items: %v", err))
        http.Error(w, "error removing items from order", http.StatusInternalServerError)
        return
    }

    h.logger.Info(fmt.Sprintf("[Order Handler.RemovingSetsDishesOrder] Proto Response: \nPagination: %+v", 
        updatedOrderResponse.GetPagination()))

    for i, order := range updatedOrderResponse.GetData() {
        h.logger.Info(fmt.Sprintf("[Order Handler.RemovingSetsDishesOrder] Proto Order %d: \n"+
            "ID: %d\n"+
            "Remaining Dish Items: %d\n"+
            "Remaining Set Items: %d\n"+
            "Current Version: %d\n"+
            "Total Price: %d",
            i+1,
            order.GetId(),
            len(order.GetDataDish()),
            len(order.GetDataSet()),
            order.GetCurrentVersion(),
            order.GetTotalPrice()))
    }

    res := ToOrderDetailedListResponseFromPbResponse(updatedOrderResponse)

    h.logger.Info(fmt.Sprintf("[Order Handler.RemovingSetsDishesOrder] Converted Response: \nPagination: %+v", 
        res.Pagination))

    for i, order := range res.Data {
        h.logger.Info(fmt.Sprintf("[Order Handler.RemovingSetsDishesOrder] Converted Order %d: \n"+
            "ID: %d\n"+
            "Final Version: %d\n"+
            "Remaining Dishes: %d\n"+
            "Remaining Sets: %d\n"+
            "Adjusted Total Price: %d",
            i+1,
            order.ID,
            order.CurrentVersion,
            len(order.DataDish),
            len(order.DataSet),
            order.TotalPrice))

   
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(res); err != nil {
        h.logger.Error(fmt.Sprintf("[Order Handler.RemovingSetsDishesOrder] Error encoding response: %v", err))
        http.Error(w, "error encoding response", http.StatusInternalServerError)
        return
    }

    h.logger.Info("[Order Handler.RemovingSetsDishesOrder] Item removal completed successfully")
}

func (h *OrderHandlerController) MarkDishesDelivered(w http.ResponseWriter, r *http.Request) {
    // Initiation logging with core identifiers
    h.logger.Info("[Order Handler.MarkDishesDelivered] Starting delivery processing")

    var dishDeliveryReq CreateDishDeliveryRequestType
    if err := json.NewDecoder(r.Body).Decode(&dishDeliveryReq); err != nil {
        h.logger.Error(fmt.Sprintf("[Order Handler.MarkDishesDelivered] Decode failure: %v", err))
        http.Error(w, "invalid request format", http.StatusBadRequest)
        return
    }

    // Structured validation logging
    h.logger.Info(fmt.Sprintf("[Order Handler.MarkDishesDelivered] Processing OrderID: %d, DeliveredByUser: %d, Items: %d",
        dishDeliveryReq.OrderID, dishDeliveryReq.DeliveredByUserID, len(dishDeliveryReq.DishItems)))

    // Validation pipeline
    if dishDeliveryReq.OrderID == 0 {
        h.logger.Error("[Order Handler.MarkDishesDelivered] Missing order ID")
        http.Error(w, "order ID required", http.StatusBadRequest)
        return
    }

    if dishDeliveryReq.DeliveredByUserID == 0 {
        h.logger.Error("[Order Handler.MarkDishesDelivered] Missing delivery user ID")
        http.Error(w, "delivery user ID required", http.StatusBadRequest)
        return
    }

    if len(dishDeliveryReq.DishItems) == 0 {
        h.logger.Warning(fmt.Sprintf("[Order Handler.MarkDishesDelivered] Empty dishes OrderID: %d", dishDeliveryReq.OrderID))
        http.Error(w, "at least one dish required", http.StatusBadRequest)
        return
    }

    var totalItems int64
    for _, dish := range dishDeliveryReq.DishItems {
        if dish.Quantity <= 0 {
            h.logger.Error(fmt.Sprintf("[Order Handler.MarkDishesDelivered] Invalid quantity %d DishID: %d",
                dish.Quantity, dish.DishID))
            http.Error(w, fmt.Sprintf("invalid quantity for dish %d", dish.DishID), http.StatusBadRequest)
            return
        }
        totalItems += dish.Quantity
    }

    // Set default values if not provided
    if dishDeliveryReq.DeliveryStatus == "" {
        dishDeliveryReq.DeliveryStatus = "DELIVERED"
    }
    if dishDeliveryReq.CreatedAt.IsZero() {
        dishDeliveryReq.CreatedAt = time.Now()
    }
    if dishDeliveryReq.UpdatedAt.IsZero() {
        dishDeliveryReq.UpdatedAt = time.Now()
    }
    if dishDeliveryReq.DeliveredAt.IsZero() {
        dishDeliveryReq.DeliveredAt = time.Now()
    }
    
    // Convert to protobuf after validation
    pbReq := ToPBCreateDishDeliveryRequest(&dishDeliveryReq)
    // h.logger.Info(fmt.Sprintf("[Order Handler.MarkDishesDelivered] Protobuf conversion complete OrderID: %d", dishDeliveryReq.OrderID))

    // Service call with error translation
    deliveryResponse, err := h.client.MarkDishesDelivered(h.ctx, pbReq)
 
    if err != nil {
        h.logger.Error(fmt.Sprintf("[Order Handler.MarkDishesDelivered] Service failure OrderID: %d - %v", dishDeliveryReq.OrderID, err))
        
        // Handle gRPC status errors
        if st, ok := status.FromError(err); ok {
            switch st.Code() {
            case codes.InvalidArgument:
                http.Error(w, st.Message(), http.StatusBadRequest)
            case codes.NotFound:
                http.Error(w, st.Message(), http.StatusNotFound)
            case codes.Internal:
                http.Error(w, st.Message(), http.StatusInternalServerError)
            default:
                http.Error(w, "processing error", http.StatusInternalServerError)
            }
        } else {
            http.Error(w, "unexpected error", http.StatusInternalServerError)
        }
        return
    }

    // Response validation and logging
    if deliveryResponse == nil {
        h.logger.Error(fmt.Sprintf("[Order Handler.MarkDishesDelivered] Empty response OrderID: %d", dishDeliveryReq.OrderID))
        http.Error(w, "empty service response", http.StatusInternalServerError)
        return
    }

    // h.logger.Info(fmt.Sprintf("[Order Handler.MarkDishesDelivered] Success OrderID: %d, NewVersion: %d",
    //     dishDeliveryReq.OrderID, deliveryResponse.GetCurrentVersion()))

    // Convert and serialize response
    res := ToOrderDetailedListResponseFromDeliveryResponse(deliveryResponse)
    // h.logger.Info(fmt.Sprintf("[MarkDishesDelivered] Delivery History handler 12121212: %+v", res.Data))
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    if err := json.NewEncoder(w).Encode(res); err != nil {
        h.logger.Error(fmt.Sprintf("[Order Handler.MarkDishesDelivered] Encoding failure OrderID: %d - %v",
            dishDeliveryReq.OrderID, err))
        http.Error(w, "response serialization error", http.StatusInternalServerError)
        return
    }

    h.logger.Info(fmt.Sprintf("[Order Handler.MarkDishesDelivered] Completed OrderID: %d", dishDeliveryReq.OrderID))
}

func ToPBCreateDishDeliveryRequest(req *CreateDishDeliveryRequestType) *order.CreateDishDeliveryRequest {
    if req == nil {
        return nil
    }

    // Calculate total quantity delivered if not provided
    quantityDelivered := req.QuantityDelivered
    if quantityDelivered == 0 {
        for _, item := range req.DishItems {
            quantityDelivered += int32(item.Quantity)
        }
    }

    // Convert DishItems to []*order.DishOrderItem
    pbDishItems := make([]*order.CreateDishOrderItem, len(req.DishItems))
    for i, item := range req.DishItems {
        pbDishItems[i] = &order.CreateDishOrderItem{
        
            DishId:           item.DishID,
            Quantity:         item.Quantity,

        }
    }

    // Convert to CreateDishDeliveryRequest format
    return &order.CreateDishDeliveryRequest{
        OrderId:           req.OrderID,
        OrderName:         req.OrderName,
        GuestId:          req.GuestID,
        UserId:           req.UserID,
        TableNumber:      req.TableNumber,
        DishItems:        pbDishItems,
        QuantityDelivered: quantityDelivered,
        DeliveryStatus:    req.DeliveryStatus,
        DeliveredAt:       timestamppb.New(req.DeliveredAt),
        DeliveredByUserId: req.DeliveredByUserID,
        CreatedAt:         timestamppb.New(req.CreatedAt),
        UpdatedAt:         timestamppb.New(req.UpdatedAt),
        IsGuest:          req.IsGuest,
    }
}


// new -----------
func ToOrderDetailedListResponseFromDeliveryResponse(pbRes *order.OrderDetailedResponseWithDelivery) *OrderDetailedListResponse {
 

    if pbRes == nil {
        fmt.Printf("[ToOrderDetailedListResponseFromDeliveryResponse] Warning: received nil protobuf response\n")
        return nil
    }



    // Convert DataSet

    dataSet := make([]OrderSetDetailed, len(pbRes.DataSet))
    for i, set := range pbRes.DataSet {
        dishes := make([]OrderDetailedDish, len(set.Dishes))
        for j, dish := range set.Dishes {
            dishes[j] = OrderDetailedDish{
                DishID:      dish.DishId,
                Quantity:    dish.Quantity,
                Name:        dish.Name,
                Price:       dish.Price,
                Description: dish.Description,
                Image:       dish.Image,
                Status:      dish.Status,
            }
        }



        dataSet[i] = OrderSetDetailed{
            ID:          set.Id,
            Name:        set.Name,
            Description: set.Description,
            Dishes:      dishes,
            UserID:      set.UserId,
            CreatedAt:   set.CreatedAt.AsTime(),
            UpdatedAt:   set.UpdatedAt.AsTime(),
            IsFavourite: set.IsFavourite,
            LikeBy:      set.LikeBy,
            IsPublic:    set.IsPublic,
            Image:       set.Image,
            Price:       set.Price,
            Quantity:    set.Quantity,
        }
    }

    // Convert DataDish

    dataDish := make([]OrderDetailedDish, len(pbRes.DataDish))
    for i, dish := range pbRes.DataDish {
        dataDish[i] = OrderDetailedDish{
            DishID:      dish.DishId,
            Quantity:    dish.Quantity,
            Name:        dish.Name,
            Price:       dish.Price,
            Description: dish.Description,
            Image:       dish.Image,
            Status:      dish.Status,
        }
    }

    // Convert VersionHistory

    versionHistory := make([]OrderVersionSummary, len(pbRes.VersionHistory))
    for i, vh := range pbRes.VersionHistory {
        changes := make([]OrderItemChange, len(vh.Changes))
        for j, change := range vh.Changes {
            changes[j] = OrderItemChange{
                ItemType:        change.ItemType,
                ItemID:          change.ItemId,
                ItemName:        change.ItemName,
                QuantityChanged: change.QuantityChanged,
                Price:          change.Price,
            }
        }


        versionHistory[i] = OrderVersionSummary{
            VersionNumber:     vh.VersionNumber,
            TotalDishesCount: vh.TotalDishesCount,
            TotalSetsCount:   vh.TotalSetsCount,
            VersionTotalPrice: vh.VersionTotalPrice,
            ModificationType: vh.ModificationType,
            ModifiedAt:      vh.ModifiedAt.AsTime(),
            Changes:         changes,
        }
    }

    // Convert MostOrderedItems

    mostOrderedItems := make([]OrderItemCount, len(pbRes.TotalSummary.MostOrderedItems))
    for i, item := range pbRes.TotalSummary.MostOrderedItems {
        mostOrderedItems[i] = OrderItemCount{
            ItemType:     item.ItemType,
            ItemID:       item.ItemId,
            ItemName:     item.ItemName,
            TotalQuantity: item.TotalQuantity,
        }
    }

    // Convert DeliveryHistory

    deliveryHistory := make([]DishDelivery, len(pbRes.DeliveryHistory))
    for i, dh := range pbRes.DeliveryHistory {
        dishItems := make([]OrderDish, len(dh.DishItems))
        for j, item := range dh.DishItems {
            dishItems[j] = OrderDish{
                ID:                item.Id,
                DishID:           item.DishId,
                Quantity:         item.Quantity,
                CreatedAt:        item.CreatedAt.AsTime(),
                UpdatedAt:        item.UpdatedAt.AsTime(),
                OrderName:        item.OrderName,
                ModificationType: item.ModificationType,
                ModificationNumber: item.ModificationNumber,
            }
        }



        deliveryHistory[i] = DishDelivery{
            ID:                dh.Id,
            OrderID:           dh.OrderId,
            OrderName:         dh.OrderName,
            GuestID:           dh.GuestId,
            UserID:            dh.UserId,
            TableNumber:       dh.TableNumber,
            DishItems:         dishItems,
            QuantityDelivered: dh.QuantityDelivered,
            DeliveryStatus:    dh.DeliveryStatus,
            DeliveredAt:       dh.DeliveredAt.AsTime(),
            DeliveredByUserID: dh.DeliveredByUserId,
            CreatedAt:         dh.CreatedAt.AsTime(),
            UpdatedAt:         dh.UpdatedAt.AsTime(),
            IsGuest:           dh.IsGuest,
            ModificationNumber: dh.ModificationNumber, 
        }

        fmt.Printf("[DeliveryHistory] 12121212121  Details: %+v\n", deliveryHistory)
    }

    response := &OrderDetailedListResponse{
        Data: []OrderDetailedResponse{
            {
                DataSet:              dataSet,
                DataDish:            dataDish,
                ID:                  pbRes.Id,
                GuestID:             pbRes.GuestId,
                UserID:              pbRes.UserId,
                TableNumber:         pbRes.TableNumber,
                OrderHandlerID:      pbRes.OrderHandlerId,
                Status:              pbRes.Status,
                TotalPrice:          pbRes.TotalPrice,
                IsGuest:             pbRes.IsGuest,
                Topping:             pbRes.Topping,
                TrackingOrder:       pbRes.TrackingOrder,
                TakeAway:            pbRes.TakeAway,
                ChiliNumber:         pbRes.ChiliNumber,
                TableToken:          pbRes.TableToken,
                OrderName:           pbRes.OrderName,
                CurrentVersion:      pbRes.CurrentVersion,
                ParentOrderID:       pbRes.ParentOrderId,
                VersionHistory:     versionHistory,
                TotalSummary: OrderTotalSummary{
                    TotalVersions:        pbRes.TotalSummary.TotalVersions,
                    TotalDishesOrdered:   pbRes.TotalSummary.TotalDishesOrdered,
                    TotalSetsOrdered:     pbRes.TotalSummary.TotalSetsOrdered,
                    CumulativeTotalPrice: pbRes.TotalSummary.CumulativeTotalPrice,
                    MostOrderedItems:     mostOrderedItems,
                },
                DeliveryHistory:      deliveryHistory,
                CurrentDeliveryStatus: DeliveryStatus(pbRes.CurrentDeliveryStatus.String()),
                TotalItemsDelivered:   pbRes.TotalItemsDelivered,
                LastDeliveryAt:        pbRes.LastDeliveryAt.AsTime(),
            },
        },
        Pagination: PaginationInfo{},
    }


    fmt.Printf("\n=== DeliveryHistory Details golang/quanqr/order/order_handler.go ===\n")
    for i, delivery := range deliveryHistory {
        fmt.Printf("\nDelivery #%d:\n", i+1)
        // fmt.Printf("------------------------\n")
        // fmt.Printf("ID: %d\n", delivery.ID)
        // fmt.Printf("OrderID: %d\n", delivery.OrderID)
        // fmt.Printf("OrderName: %s\n", delivery.OrderName)
        // fmt.Printf("GuestID: %d\n", delivery.GuestID)
        // fmt.Printf("UserID: %d\n", delivery.UserID)
        // fmt.Printf("TableNumber: %d\n", delivery.TableNumber)
        // fmt.Printf("QuantityDelivered: %d\n", delivery.QuantityDelivered)
        // fmt.Printf("DeliveryStatus: %s\n", delivery.DeliveryStatus)
        // fmt.Printf("DeliveredAt: %v\n", delivery.DeliveredAt)
        // fmt.Printf("DeliveredByUserID: %d\n", delivery.DeliveredByUserID)
        // fmt.Printf("CreatedAt: %v\n", delivery.CreatedAt)
        // fmt.Printf("UpdatedAt: %v\n", delivery.UpdatedAt)
        // fmt.Printf("IsGuest: %v\n", delivery.IsGuest)
        fmt.Printf("ModificationNumber: %d\n", delivery.ModificationNumber)
        
        fmt.Printf("DishItems:\n")
        for j, item := range delivery.DishItems {
            fmt.Printf("  Item #%d:\n", j+1)
            fmt.Printf("    - ID: %d\n", item.ID)
            fmt.Printf("    - Quantity: %d\n", item.Quantity)
            // Add any other DishItems fields you want to see
        }
        fmt.Printf("\n")
    }
    fmt.Printf("=== End of DeliveryHistory Details golang/quanqr/order/order_handler.go ===\n\n")

    return response
}

// 
func ToOrderDetailedListResponseFromProto(pbRes *order.OrderDetailedListResponse) OrderDetailedListResponse {
    if pbRes == nil {
        return OrderDetailedListResponse{}
    }

    detailedResponses := make([]OrderDetailedResponse, len(pbRes.Data))
    for i, pbDetailedRes := range pbRes.Data {
        if pbDetailedRes == nil {
            continue
        }

        // Convert version history with nil checks
        var versionHistory []OrderVersionSummary
        if pbDetailedRes.VersionHistory != nil {
            versionHistory = make([]OrderVersionSummary, len(pbDetailedRes.VersionHistory))
            for j, pbVersion := range pbDetailedRes.VersionHistory {
                if pbVersion == nil {
                    continue
                }

                var changes []OrderItemChange
                if pbVersion.Changes != nil {
                    changes = make([]OrderItemChange, len(pbVersion.Changes))
                    for k, pbChange := range pbVersion.Changes {
                        if pbChange == nil {
                            continue
                        }
                        changes[k] = OrderItemChange{
                            ItemType:        pbChange.ItemType,
                            ItemID:          pbChange.ItemId,
                            ItemName:        pbChange.ItemName,
                            QuantityChanged: pbChange.QuantityChanged,
                            Price:           pbChange.Price,
                        }
                    }
                }

                var modifiedAt time.Time
                if pbVersion.ModifiedAt != nil {
                    modifiedAt = pbVersion.ModifiedAt.AsTime()
                }
                
                versionHistory[j] = OrderVersionSummary{
                    VersionNumber:     pbVersion.VersionNumber,
                    TotalDishesCount: pbVersion.TotalDishesCount,
                    TotalSetsCount:   pbVersion.TotalSetsCount,
                    VersionTotalPrice: pbVersion.VersionTotalPrice,
                    ModificationType: pbVersion.ModificationType,
                    ModifiedAt:      modifiedAt,
                    Changes:         changes,
                }
            }
        }

        // Convert total summary with nil checks
        var totalSummary OrderTotalSummary
        if pbDetailedRes.TotalSummary != nil {
            var mostOrderedItems []OrderItemCount
            if pbDetailedRes.TotalSummary.MostOrderedItems != nil {
                mostOrderedItems = make([]OrderItemCount, len(pbDetailedRes.TotalSummary.MostOrderedItems))
                for j, pbItem := range pbDetailedRes.TotalSummary.MostOrderedItems {
                    if pbItem == nil {
                        continue
                    }
                    mostOrderedItems[j] = OrderItemCount{
                        ItemType:      pbItem.ItemType,
                        ItemID:        pbItem.ItemId,
                        ItemName:      pbItem.ItemName,
                        TotalQuantity: pbItem.TotalQuantity,
                    }
                }
            }

            totalSummary = OrderTotalSummary{
                TotalVersions:        pbDetailedRes.TotalSummary.TotalVersions,
                TotalDishesOrdered:  pbDetailedRes.TotalSummary.TotalDishesOrdered,
                TotalSetsOrdered:    pbDetailedRes.TotalSummary.TotalSetsOrdered,
                CumulativeTotalPrice: pbDetailedRes.TotalSummary.CumulativeTotalPrice,
                MostOrderedItems:     mostOrderedItems,
            }
        }

        var dataSet []OrderSetDetailed
        if pbDetailedRes.DataSet != nil {
            dataSet = ToOrderSetsDetailedFromProto(pbDetailedRes.DataSet)
        }

        var dataDish []OrderDetailedDish
        if pbDetailedRes.DataDish != nil {
            dataDish = ToOrderDetailedDishesFromProto(pbDetailedRes.DataDish)
        }

        detailedResponses[i] = OrderDetailedResponse{
            ID:             pbDetailedRes.Id,
            GuestID:        pbDetailedRes.GuestId,
            UserID:         pbDetailedRes.UserId,
            TableNumber:    pbDetailedRes.TableNumber,
            OrderHandlerID: pbDetailedRes.OrderHandlerId,
            Status:         pbDetailedRes.Status,
            TotalPrice:     pbDetailedRes.TotalPrice,
            IsGuest:        pbDetailedRes.IsGuest,
            Topping:        pbDetailedRes.Topping,
            TrackingOrder:  pbDetailedRes.TrackingOrder,
            TakeAway:       pbDetailedRes.TakeAway,
            ChiliNumber:    pbDetailedRes.ChiliNumber,
            TableToken:     pbDetailedRes.TableToken,
            OrderName:      pbDetailedRes.OrderName,
            CurrentVersion: pbDetailedRes.CurrentVersion,
            ParentOrderID:  pbDetailedRes.ParentOrderId,
            DataSet:        dataSet,
            DataDish:       dataDish,
            VersionHistory: versionHistory,
            TotalSummary:   totalSummary,
        }
    }

    var pagination PaginationInfo
    if pbRes.Pagination != nil {
        pagination = PaginationInfo{
            CurrentPage: pbRes.Pagination.CurrentPage,
            TotalPages: pbRes.Pagination.TotalPages,
            TotalItems: pbRes.Pagination.TotalItems,
            PageSize:   pbRes.Pagination.PageSize,
        }
    }

    return OrderDetailedListResponse{
        Data:       detailedResponses,
        Pagination: pagination,
    }
}

// new 

func (h *OrderHandlerController) AddingDishesToOrder(w http.ResponseWriter, r *http.Request) {
    // Log the incoming request
    h.logger.Info("golang/quanqr/order/order_handler.go Received AddingDishesToOrder request")

    // Parse and validate the request body
    var req CreateDishOrderItemWithOrderID
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Error("golang/quanqr/order/order_handler.go Error decoding request body: " + err.Error(),
            "raw_body", r.Body)
        http.Error(w, "error decoding request body", http.StatusBadRequest)
        return
    }

    // Log all input fields
    h.logger.Info("golang/quanqr/order/order_handler.go Request details",
        "order_id", req.OrderID,
        "dish_id", req.DishID,
        "quantity", req.Quantity,
        "order_name", req.OrderName)

    // Basic validation with detailed logging
    if req.OrderID <= 0 {
        h.logger.Error("golang/quanqr/order/order_handler.go Invalid order ID",
            "provided_order_id", req.OrderID)
        http.Error(w, "invalid order ID", http.StatusBadRequest)
        return
    }
    if req.DishID <= 0 {
        h.logger.Error("golang/quanqr/order/order_handler.go Invalid dish ID",
            "provided_dish_id", req.DishID)
        http.Error(w, "invalid dish ID", http.StatusBadRequest)
        return
    }
    if req.Quantity <= 0 {
        h.logger.Error("golang/quanqr/order/order_handler.go Invalid quantity",
            "provided_quantity", req.Quantity)
        http.Error(w, "quantity must be positive", http.StatusBadRequest)
        return
    }
    if req.OrderName == "" {
        h.logger.Error("golang/quanqr/order/order_handler.go Order name is required",
            "provided_order_name", req.OrderName)
        http.Error(w, "order name is required", http.StatusBadRequest)
        return
    }

    // Log successful validation
    h.logger.Info("golang/quanqr/order/order_handler.go Request validation successful")

    // Convert request to protobuf format
    pbReq := &order.CreateDishOrderItemWithOrderID{
        OrderId:   req.OrderID,
        DishId:    req.DishID,
        Quantity:  req.Quantity,
        OrderName: req.OrderName,
    }

    // Log before service call
    h.logger.Info("golang/quanqr/order/order_handler.go Calling gRPC service AddingDishesToOrder",
        "order_id", pbReq.OrderId,
        "dish_id", pbReq.DishId,
        "quantity", pbReq.Quantity,
        "order_name", pbReq.OrderName)

    response, err := h.client.AddingDishesToOrder(h.ctx, pbReq)
    if err != nil {
        if s, ok := status.FromError(err); ok {
            h.logger.Error("golang/quanqr/order/order_handler.go gRPC service returned error",
                "error_code", s.Code(),
                "error_message", s.Message(),
                "request_data", pbReq)
            switch s.Code() {
            case codes.InvalidArgument:
                http.Error(w, s.Message(), http.StatusBadRequest)
            default:
                http.Error(w, "internal server error", http.StatusInternalServerError)
            }
        } else {
            h.logger.Error("golang/quanqr/order/order_handler.go Non-gRPC error occurred",
                "error", err.Error(),
                "request_data", pbReq)
            http.Error(w, "internal server error", http.StatusInternalServerError)
        }
        return
    }

    // Log successful response
    h.logger.Info("golang/quanqr/order/order_handler.go Successfully received response from gRPC service",
        "response_id", response.Id,
        "dish_id", response.DishId,
        "quantity", response.Quantity,
        "order_name", response.OrderName,
        "modification_type", response.ModificationType,
        "modification_number", response.ModificationNumber)

    // Convert response back to JSON format
    result := OrderDish{
        ID:                response.Id,
        DishID:           response.DishId,
        Quantity:         response.Quantity,
        CreatedAt:        time.Unix(response.CreatedAt.GetSeconds(), int64(response.CreatedAt.GetNanos())),
        UpdatedAt:        time.Unix(response.UpdatedAt.GetSeconds(), int64(response.UpdatedAt.GetNanos())),
        OrderName:        response.OrderName,
        ModificationType: response.ModificationType,
        ModificationNumber: response.ModificationNumber,
    }

    // Send response
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(result); err != nil {
        h.logger.Error("golang/quanqr/order/order_handler.go Error encoding response",
            "error", err.Error(),
            "result", result)
        http.Error(w, "error encoding response", http.StatusInternalServerError)
        return
    }

    // Log successful completion
    h.logger.Info("golang/quanqr/order/order_handler.go Successfully completed AddingDishesToOrder request",
        "order_id", req.OrderID,
        "dish_id", req.DishID,
        "response_id", result.ID)
}

func (h *OrderHandlerController) AddingSetToOrder(w http.ResponseWriter, r *http.Request) {
    // Parse and validate the request body
    var req CreateSetOrderItemWithOrderID
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Error("Error decoding request body: " + err.Error())
        http.Error(w, "error decoding request body", http.StatusBadRequest)
        return
    }

    // Basic validation
    if req.OrderID <= 0 {
        h.logger.Error("Invalid order ID")
        http.Error(w, "invalid order ID", http.StatusBadRequest)
        return
    }
    if req.SetID <= 0 {
        h.logger.Error("Invalid set ID")
        http.Error(w, "invalid set ID", http.StatusBadRequest)
        return
    }
    if req.Quantity <= 0 {
        h.logger.Error("Invalid quantity")
        http.Error(w, "quantity must be positive", http.StatusBadRequest)
        return
    }
    if req.OrderName == "" {
        h.logger.Error("Order name is required")
        http.Error(w, "order name is required", http.StatusBadRequest)
        return
    }

    // Convert request to protobuf format
    pbReq := &order.CreateSetOrderItemWithOrderID{
        OrderId:   req.OrderID,
        SetId:     req.SetID,
        Quantity:  req.Quantity,
        OrderName: req.OrderName,
    }

    // Call the service
    h.logger.Info(fmt.Sprintf("Adding set to order ID %d: set_id=%d, quantity=%d", 
        req.OrderID, req.SetID, req.Quantity))
    response, err := h.client.AddingSetToOrder(h.ctx, pbReq)
    if err != nil {
        h.logger.Error("Error adding set to order: " + err.Error())
        if s, ok := status.FromError(err); ok {
            switch s.Code() {
            case codes.InvalidArgument:
                http.Error(w, s.Message(), http.StatusBadRequest)
            default:
                http.Error(w, "internal server error", http.StatusInternalServerError)
            }
        } else {
            http.Error(w, "internal server error", http.StatusInternalServerError)
        }
        return
    }

    // Convert response back to JSON format
    result := ResponseSetOrderItemWithOrderID{
        Set: OrderSet{
            ID:                response.Set.Id,
            SetID:            response.Set.SetId,
            Quantity:         response.Set.Quantity,
            CreatedAt:        time.Unix(response.Set.CreatedAt.GetSeconds(), int64(response.Set.CreatedAt.GetNanos())),
            UpdatedAt:        time.Unix(response.Set.UpdatedAt.GetSeconds(), int64(response.Set.UpdatedAt.GetNanos())),
            OrderName:        response.Set.OrderName,
            ModificationType: response.Set.ModificationType,
            ModificationNumber: response.Set.ModificationNumber,
        },
        Dishes: make([]OrderDetailedDish, len(response.Dishes)),
    }

    // Convert dishes
    for i, dish := range response.Dishes {
        result.Dishes[i] = OrderDetailedDish{
            DishID:      dish.DishId,
            Quantity:    dish.Quantity,
            Name:        dish.Name,
            Price:       dish.Price,
            Description: dish.Description,
            Image:       dish.Image,
            Status:      dish.Status,
        }
    }

    // Send response
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(result); err != nil {
        h.logger.Error("Error encoding response: " + err.Error())
        http.Error(w, "error encoding response", http.StatusInternalServerError)
        return
    }
}