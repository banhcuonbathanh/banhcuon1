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
    os.logger.Info("Fetching orders list")
    
    orders, err := os.orderRepo.GetOrders(ctx, req)
    if err != nil {
        os.logger.Error("Error fetching orders: " + err.Error())
        return nil, err
    }
    
    return &order.OrderListResponse{
        Data: orders,
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

func (os *OrderServiceStruct) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderResponse, error) {
    os.logger.Info("Updating order: " + fmt.Sprint(req.Id))
    
    updatedOrder, err := os.orderRepo.UpdateOrder(ctx, req)
    if err != nil {
        os.logger.Error("Error updating order: " + err.Error())
        return nil, err
    }
    
    return &order.OrderResponse{
        Data: updatedOrder,
    }, nil
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