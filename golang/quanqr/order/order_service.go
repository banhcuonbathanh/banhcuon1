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

func (os *OrderServiceStruct) CreateOrders(ctx context.Context, req *order.CreateOrdersRequest) (*order.OrderListResponse, error) {
	log.Println("Creating new orders for guest:", req.GuestId)

	createdOrders, err := os.orderRepo.CreateOrders(ctx, req)
	if err != nil {
		log.Println("Error creating orders:", err)
		return nil, err
	}

	log.Println("Orders created successfully for guest:", req.GuestId)
	return &order.OrderListResponse{
		Data:    createdOrders,
		Message: "Orders created successfully",
	}, nil
}

func (os *OrderServiceStruct) GetOrders(ctx context.Context, req *order.GetOrdersRequest) (*order.OrderListResponse, error) {
	log.Println("Fetching orders from", req.FromDate, "to", req.ToDate)

	orders, err := os.orderRepo.GetOrders(ctx, req)
	if err != nil {
		log.Println("Error fetching orders:", err)
		return nil, err
	}

	return &order.OrderListResponse{
		Data:    orders,
		Message: "Orders fetched successfully",
	}, nil
}

func (os *OrderServiceStruct) GetOrderDetail(ctx context.Context, req *order.OrderDetailIdParam) (*order.OrderResponse, error) {
	log.Println("Fetching order detail for ID:", req.Id)

	orderDetail, err := os.orderRepo.GetOrderDetail(ctx, req.Id)
	if err != nil {
		log.Println("Error fetching order detail:", err)
		return nil, err
	}

	return &order.OrderResponse{
		Data:    orderDetail,
		Message: "Order detail fetched successfully",
	}, nil
}

func (os *OrderServiceStruct) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.OrderResponse, error) {
	log.Println("Updating order:", req.OrderId)

	updatedOrder, err := os.orderRepo.UpdateOrder(ctx, req)
	if err != nil {
		log.Println("Error updating order:", err)
		return nil, err
	}

	return &order.OrderResponse{
		Data:    updatedOrder,
		Message: "Order updated successfully",
	}, nil
}

func (os *OrderServiceStruct) PayGuestOrders(ctx context.Context, req *order.PayGuestOrdersRequest) (*order.OrderListResponse, error) {
	log.Println("Paying orders for guest:", req.GuestId)

	paidOrders, err := os.orderRepo.PayGuestOrders(ctx, req.GuestId)
	if err != nil {
		log.Println("Error paying guest orders:", err)
		return nil, err
	}

	return &order.OrderListResponse{
		Data:    paidOrders,
		Message: "Guest orders paid successfully",
	}, nil
}