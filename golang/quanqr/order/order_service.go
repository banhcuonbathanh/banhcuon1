package order_grpc

import (
	"context"
	"log"

	"english-ai-full/quanqr/proto_qr/order"
)

type OrderServiceStruct struct {
	orderRepo *OrderRepository
	order.UnimplementedOrderServiceServer
}

func NewOrderService(orderRepo *OrderRepository) *OrderServiceStruct {
	return &OrderServiceStruct{
		orderRepo: orderRepo,
	}
}

func (os *OrderServiceStruct) CreateOrders(ctx context.Context, req *order.CreateOrderRequest) (*order.OrderListResponse, error) {
	log.Println("Creating new order:",
		"GuestId:", req.GuestId,
		"TableNumber:", req.TableNumber,
		"DishSnapshotId:", req.DishSnapshotId,
		"OrderHandlerId:", req.OrderHandlerId,
		"Status:", req.Status,
		"TotalPrice:", req.TotalPrice,
	)

	createdOrders, err := os.orderRepo.CreateOrders(ctx, req)
	if err != nil {
		log.Println("Error creating order:", err)
		return nil, err
	}

	log.Println("Order created successfully. ID:", createdOrders[0].Id)
	return &order.OrderListResponse{
		Data:    createdOrders,
	
	}, nil
}

func (os *OrderServiceStruct) GetOrders(ctx context.Context, req *order.GetOrdersRequest) (*order.OrderListResponse, error) {
	orders, err := os.orderRepo.GetOrders(ctx, req)
	if err != nil {
		log.Println("Error fetching orders:", err)
		return nil, err
	}
	return &order.OrderListResponse{
		Data:    orders,

	}, nil
}

func (os *OrderServiceStruct) GetOrderDetail(ctx context.Context, req *order.OrderDetailIdParam) (*order.OrderResponse, error) {
	o, err := os.orderRepo.GetOrderDetail(ctx, req.Id)
	if err != nil {
		log.Println("Error fetching order detail:", err)
		return nil, err
	}
	return &order.OrderResponse{
		Data:    o,
	
	}, nil
}

func (os *OrderServiceStruct) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderResponse, error) {
	updatedOrder, err := os.orderRepo.UpdateOrder(ctx, req)
	if err != nil {
		log.Println("Error updating order:", err)
		return nil, err
	}
	return &order.OrderResponse{
		Data:    updatedOrder,
	
	}, nil
}

func (os *OrderServiceStruct) PayGuestOrders(ctx context.Context, req *order.PayGuestOrdersRequest) (*order.OrderListResponse, error) {
	paidOrders, err := os.orderRepo.PayGuestOrders(ctx, req.GuestId)
	if err != nil {
		log.Println("Error paying guest orders:", err)
		return nil, err
	}
	return &order.OrderListResponse{
		Data:    paidOrders,
	
	}, nil
}