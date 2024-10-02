package qr_guests

import (
	"context"
	"log"

	"english-ai-full/quanqr/proto_qr/guest"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GuestServiceStruct struct {
	guestRepo *GuestRepository
	guest.UnimplementedGuestServiceServer
}

func NewGuestService(guestRepo *GuestRepository) *GuestServiceStruct {
	return &GuestServiceStruct{
		guestRepo: guestRepo,
	}
}

func (gs *GuestServiceStruct) GuestLoginGRPC(ctx context.Context, req *guest.GuestLoginRequest) (*guest.GuestLoginResponse, error) {
	log.Println("Guest login attempt:",
		"Name:", req.Name,
		"Table Number:", req.TableNumber,
	)

	response, err := gs.guestRepo.GuestLogin(ctx, req)
	if err != nil {
		log.Println("Error during guest login:", err)
		return nil, err
	}

	log.Println("Guest logged in successfully. Guest ID:", response.Guest.Id)
	return response, nil
}

func (gs *GuestServiceStruct) GuestLogoutGRPC(ctx context.Context, req *guest.LogoutRequest) (*emptypb.Empty, error) {
	log.Println("Guest logout attempt")

	err := gs.guestRepo.GuestLogout(ctx, req)
	if err != nil {
		log.Println("Error during guest logout:", err)
		return nil, err
	}

	log.Println("Guest logged out successfully")
	return &emptypb.Empty{}, nil
}

func (gs *GuestServiceStruct) GuestRefreshTokenGRPC(ctx context.Context, req *guest.RefreshTokenRequest) (*guest.RefreshTokenResponse, error) {
	log.Println("Token refresh attempt")

	response, err := gs.guestRepo.RefreshToken(ctx, req)
	if err != nil {
		log.Println("Error during token refresh:", err)
		return nil, err
	}

	log.Println("Token refreshed successfully")
	return response, nil
}

func (gs *GuestServiceStruct) GuestCreateOrdersGRPC(ctx context.Context, req *guest.CreateOrdersRequest) (*guest.OrdersResponse, error) {
	log.Println("Create orders attempt:",
		"Number of items:", len(req.Items),
	)

	// Validate that all items have the same guest_Id
	if len(req.Items) > 0 {
		guestID := req.Items[0].GuestId
		for _, item := range req.Items[1:] {
			if item.GuestId != guestID {
				return nil, status.Errorf(codes.InvalidArgument, "All items must have the same guest_Id")
			}
		}
	}

	response, err := gs.guestRepo.CreateOrders(ctx, req)
	if err != nil {
		log.Println("Error creating orders:", err)
		return nil, status.Errorf(codes.Internal, "Failed to create orders: %v", err)
	}

	log.Println("Orders created successfully. Number of orders:", len(response.Data))
	return response, nil
}

func (gs *GuestServiceStruct) GuestGetOrdersGRPC(ctx context.Context, req *guest.GuestGetOrdersGRPCRequest) (*guest.ListOrdersResponse, error) {
	log.Println("Get orders attempt for guest ID:", req.GuestId)

	response, err := gs.guestRepo.GetOrders(ctx, req)
	if err != nil {
		log.Println("Error fetching orders:", err)
		return nil, err
	}

	log.Println("Orders fetched successfully. Number of orders:", len(response.Orders))
	return response, nil
}


// func getGuestIDFromContext(ctx context.Context) int64 {
// 	// Implement your logic to extract guest ID from context
// 	// This might involve parsing a JWT token or fetching from a session store
// 	return 0 // Placeholder return
// }