package order_grpc

import (
	"context"
	"fmt"

	"english-ai-full/logger"
	"english-ai-full/quanqr/proto_qr/order"
)
type OrderServiceStruct struct {
    orderRepo *OrderRepository
    logger    *logger.Logger
    order.UnimplementedOrderServiceServer
}

func NewOrderService(orderRepo *OrderRepository) *OrderServiceStruct {
    return &OrderServiceStruct{
        orderRepo: orderRepo,
        logger:    logger.NewLogger(),
    }
}

func (os *OrderServiceStruct) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.OrderResponse, error) {
    os.logger.Info(fmt.Sprintf("Creating new order: %+v", req))
    
    createdOrder, err := os.orderRepo.CreateOrder(ctx, req)
    if err != nil {
        os.logger.Error("Error creating order: " + err.Error())
        return nil, err
    }

    os.logger.Info("Order created successfully. ID: " + fmt.Sprint(createdOrder.Id))
    return &order.OrderResponse{
        Data: createdOrder,
    }, nil
}

func (os *OrderServiceStruct) GetOrders(ctx context.Context, req *order.GetOrdersRequest) (*order.OrderListResponse, error) {
    os.logger.Info("Fetching orders list with pagination")
    
    // Validate pagination parameters
    if req.Page < 1 {
        req.Page = 1
    }
    if req.PageSize < 1 {
        req.PageSize = 10 // Default page size
    }
    
    orders, totalItems, err := os.orderRepo.GetOrders(ctx, req.Page, req.PageSize)
    if err != nil {
        os.logger.Error("Error fetching orders: " + err.Error())
        return nil, fmt.Errorf("failed to fetch orders: %w", err)
    }
    
    // Calculate total pages
    totalPages := (totalItems + int64(req.PageSize) - 1) / int64(req.PageSize)
    
    return &order.OrderListResponse{
        Data: orders,
        Pagination: &order.PaginationInfo{
            CurrentPage: req.Page,
            TotalPages: int32(totalPages),
            TotalItems: totalItems,
            PageSize:   req.PageSize,
        },
    }, nil
}


func (os *OrderServiceStruct) GetOrderDetail(ctx context.Context, req *order.OrderIdParam) (*order.OrderResponse, error) {
    os.logger.Info("Fetching order detail for ID: " + fmt.Sprint(req.Id))
    
    orderDetail, err := os.orderRepo.GetOrderDetail(ctx, req.Id)
    if err != nil {
        os.logger.Error("Error fetching order detail: " + err.Error())
        return nil, err
    }
    
    return &order.OrderResponse{
        Data: orderDetail,
    }, nil
}

func (os *OrderServiceStruct) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderDetailedListResponse, error) {
    os.logger.Info("Updating order: " + fmt.Sprint(req.Id))
    
    updatedOrder, err := os.orderRepo.UpdateOrder(ctx, req)
    if err != nil {
        os.logger.Error("Error updating order: " + err.Error())
        return nil, err
    }
    
    // Assuming updatedOrder is already an OrderDetailedListResponse
    return updatedOrder, nil
}
func (os *OrderServiceStruct) PayOrders(ctx context.Context, req *order.PayOrdersRequest) (*order.OrderListResponse, error) {
    os.logger.Info("Processing payment for orders")
    
    paidOrders, err := os.orderRepo.PayOrders(ctx, req)
    if err != nil {
        os.logger.Error("Error processing payment for orders: " + err.Error())
        return nil, err
    }
    
    return &order.OrderListResponse{
        Data: paidOrders,
    }, nil
}


// Add GetOrderProtoListDetail to OrderServiceStruct
func (os *OrderServiceStruct) GetOrderProtoListDetail(ctx context.Context, req *order.GetOrdersRequest) (*order.OrderDetailedListResponse, error) {
    // os.logger.Info("Fetching detailed order list with pagination")

    // Validate pagination parameters
    if req.Page < 1 {
        req.Page = 1
    }
    if req.PageSize < 1 {
        req.PageSize = 10 // Default page size
    }

    // Call repository method
    detailedList, err := os.orderRepo.GetOrderProtoListDetail(ctx, req.Page, req.PageSize)
    if err != nil {
        os.logger.Error("Error fetching detailed order list: " + err.Error())
        return nil, fmt.Errorf("failed to fetch detailed order list: %w", err)
    }

    return detailedList, nil
}



func (os *OrderServiceStruct) FetchOrdersByCriteria(ctx context.Context, req *order.FetchOrdersByCriteriaRequest) (*order.OrderDetailedListResponse, error) {
    os.logger.Info("Fetching orders by criteria")
    
    // Validate pagination parameters
    if req.Page < 1 {
        req.Page = 1
    }
    if req.PageSize < 1 {
        req.PageSize = 10 // Default page size
    }
    
    detailedList, err := os.orderRepo.FetchOrdersByCriteria(ctx, req)
    if err != nil {
        os.logger.Error("Error fetching orders by criteria: " + err.Error())
        return nil, fmt.Errorf("failed to fetch orders by criteria: %w", err)
    }
    
    return detailedList, nil
}


func (os *OrderServiceStruct) AddingSetsDishesOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderDetailedListResponse, error) {
    // // Log the incoming request ID with context for order modifications
    // os.logger.Info(fmt.Sprintf("[Order Service.AddingSetsDishesOrder] Starting addingSetsDishesOrder with order service 1212 ID: %d, table: %d, order name: %s, version: %d",
    //     req.Id, req.TableNumber, req.OrderName, req.Version))
    
    // // Log the details of dishes and sets being added
    // os.logger.Info(fmt.Sprintf("[Order Service.AddingSetsDishesOrder] Adding order items - Dishes: %d, Sets: %d",
    //     len(req.DishItems), len(req.SetItems)))
    
    // Detailed logging of dish items
    for _, dish := range req.DishItems {
        os.logger.Info(fmt.Sprintf("[Order Service.AddingSetsDishesOrder] Dish item details - ID: %d, Quantity: %d, Order Name: %s, Modification: %s, Mod Number: %d",
            dish.DishId, dish.Quantity, dish.OrderName, dish.ModificationType, dish.ModificationNumber))
    }
    
    // Detailed logging of set items
    for _, set := range req.SetItems {
        os.logger.Info(fmt.Sprintf("[Order Service.AddingSetsDishesOrder] Set item details - ID: %d, Quantity: %d, Order Name: %s, Modification: %s, Mod Number: %d",
            set.SetId, set.Quantity, set.OrderName, set.ModificationType, set.ModificationNumber))
    }
    
    // Attempt to update the order through repository
    updatedOrderResponse, err := os.orderRepo.AddingSetsDishesOrder(ctx, req)
    if err != nil {
        errMsg := fmt.Sprintf("Failed to update order %d (version %d): %s",
            req.Id, req.Version, err.Error())
        os.logger.Info(errMsg)
        return nil, fmt.Errorf("dsfg")
    }
    
    // Log successful update with detailed order information
    if updatedOrderResponse != nil && len(updatedOrderResponse.Data) > 0 {
        latestOrder := updatedOrderResponse.Data[0]
        os.logger.Info(fmt.Sprintf(" [Order Service.AddingSetsDishesOrder] Successfully updated order - ID: %d, Status: %s, Version: %d, Total Price: %d",
            latestOrder.Id, latestOrder.Status, latestOrder.CurrentVersion, latestOrder.TotalPrice))
        
        // Log version history if available
        if len(latestOrder.VersionHistory) > 0 {
            latestVersion := latestOrder.VersionHistory[len(latestOrder.VersionHistory)-1]
            os.logger.Info(fmt.Sprintf(" [Order Service.AddingSetsDishesOrder] Version update details - Number: %d, Dishes: %d, Sets: %d, Price: %d, Type: %s",
                latestVersion.VersionNumber,
                latestVersion.TotalDishesCount,
                latestVersion.TotalSetsCount,
                latestVersion.VersionTotalPrice,
                latestVersion.ModificationType))
        }
        
        // Log total summary for the order
        if latestOrder.TotalSummary != nil {
            os.logger.Info(fmt.Sprintf("[Order Service.AddingSetsDishesOrder] Order total summary - Versions: %d, Total Dishes: %d, Total Sets: %d, Total Price: %d",
                latestOrder.TotalSummary.TotalVersions,
                latestOrder.TotalSummary.TotalDishesOrdered,
                latestOrder.TotalSummary.TotalSetsOrdered,
                latestOrder.TotalSummary.CumulativeTotalPrice))
        }
    }

    return updatedOrderResponse, nil
}


func (os *OrderServiceStruct) RemovingSetsDishesOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderDetailedListResponse, error) {
    // Log initiation of removal process
    os.logger.Info(fmt.Sprintf("[Order Service.RemovingSetsDishesOrder] Starting removal for OrderID: %d, Table: %d, Version: %d",
        req.Id, req.TableNumber, req.Version))

    // Log quantities being removed
    os.logger.Info(fmt.Sprintf("[Order Service.RemovingSetsDishesOrder] Removing items - Dishes: %d, Sets: %d",
        len(req.DishItems), len(req.SetItems)))

    // Detailed logging for dish removals
    for _, dish := range req.DishItems {
        os.logger.Info(fmt.Sprintf("[Order Service.RemovingSetsDishesOrder] Dish removal details - ID: %d, Quantity: %d, Order Name: %s",
            dish.DishId, dish.Quantity, dish.OrderName))
    }

    // Detailed logging for set removals
    for _, set := range req.SetItems {
        os.logger.Info(fmt.Sprintf("[Order Service.RemovingSetsDishesOrder] Set removal details - ID: %d, Quantity: %d, Order Name: %s",
            set.SetId, set.Quantity, set.OrderName))
    }

    // Execute removal through repository
    updatedOrderResponse, err := os.orderRepo.RemovingSetsDishesOrder(ctx, req)
    if err != nil {
        errMsg := fmt.Sprintf("[Order Service.RemovingSetsDishesOrder] Failed to remove items from order %d: %s",
            req.Id, err.Error())
        os.logger.Error(errMsg)
        return nil, fmt.Errorf("failed to remove items from order: %w", err)
    }

    // Log successful removal details
    if updatedOrderResponse != nil && len(updatedOrderResponse.Data) > 0 {
        latestOrder := updatedOrderResponse.Data[0]
        os.logger.Info(fmt.Sprintf("[Order Service.RemovingSetsDishesOrder] Successfully removed items - OrderID: %d, NewVersion: %d, TotalPrice: %d",
            latestOrder.Id, latestOrder.CurrentVersion, latestOrder.TotalPrice))

        // Log version history changes
        if len(latestOrder.VersionHistory) > 0 {
            latestVersion := latestOrder.VersionHistory[len(latestOrder.VersionHistory)-1]
            os.logger.Info(fmt.Sprintf("[Order Service.RemovingSetsDishesOrder] Version update - Number: %d, Type: %s, PriceImpact: %d",
                latestVersion.VersionNumber,
                latestVersion.ModificationType,
                latestVersion.VersionTotalPrice))
        }

        // Log post-removal summary
        if latestOrder.TotalSummary != nil {
            os.logger.Info(fmt.Sprintf("[Order Service.RemovingSetsDishesOrder] Post-removal summary - TotalDishes: %d, TotalSets: %d, TotalPrice: %d",
                latestOrder.TotalSummary.TotalDishesOrdered,
                latestOrder.TotalSummary.TotalSetsOrdered,
                latestOrder.TotalSummary.CumulativeTotalPrice))
        }
    }

    return updatedOrderResponse, nil
}

func (os *OrderServiceStruct) MarkDishesDelivered(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderDetailedListResponse, error) {
    // Log delivery initiation with version check
    os.logger.Info(fmt.Sprintf("[OrderService.MarkDishesDelivered] Initiating delivery for OrderID: %d, Handler: %d, Version: %d",
        req.Id, req.OrderHandlerId, req.Version))

    // Validate mandatory delivery fields
    if req.OrderHandlerId == 0 {
        os.logger.Error("[OrderService.MarkDishesDelivered] Missing order handler ID")
        return nil, fmt.Errorf("delivery requires valid handler identification")
    }

    if len(req.DishItems) == 0 {
        os.logger.Warning(fmt.Sprintf("[OrderService.MarkDishesDelivered] Empty delivery attempt OrderID: %d", req.Id))
        return nil, fmt.Errorf("at least one dish must be specified for delivery")
    }

    // Log delivery details with quantity validation
    totalDelivered := 0
    for _, dish := range req.DishItems {
        if dish.Quantity <= 0 {
            os.logger.Error(fmt.Sprintf("[OrderService.MarkDishesDelivered] Invalid quantity %d for DishID: %d",
                dish.Quantity, dish.DishId))
            return nil, fmt.Errorf("invalid quantity for dish %d: must be positive", dish.DishId)
        }
        
        os.logger.Info(fmt.Sprintf("[OrderService.MarkDishesDelivered] Scheduling delivery - DishID: %d, Qty: %d, Handler: %d",
            dish.DishId, dish.Quantity, req.OrderHandlerId))
        
        totalDelivered += int(dish.Quantity)
    }

    os.logger.Info(fmt.Sprintf("[OrderService.MarkDishesDelivered] Total items to deliver: %d", totalDelivered))

    // Execute delivery through repository
    deliveryResponse, err := os.orderRepo.MarkDishesDelivered(ctx, req)
    if err != nil {
        errMsg := fmt.Sprintf("[OrderService.MarkDishesDelivered] Delivery failed for OrderID: %d - %s",
            req.Id, err.Error())
        os.logger.Error(errMsg)
        return nil, fmt.Errorf("delivery processing failed: %w", err)
    }

    // Log delivery outcome
    if deliveryResponse != nil && len(deliveryResponse.Data) > 0 {
        orderDetails := deliveryResponse.Data[0]
        
        os.logger.Info(fmt.Sprintf("[OrderService.MarkDishesDelivered] Successfully delivered - OrderID: %d, NewVersion: %d, TotalPrice: %d",
            orderDetails.Id, orderDetails.CurrentVersion, orderDetails.TotalPrice))

        // Log delivery audit trail
 

        // Log delivery impact analysis
  

 
    }

    return deliveryResponse, nil
}