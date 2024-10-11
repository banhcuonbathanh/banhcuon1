package set_qr

import (
	"context"
	"log"

	"english-ai-full/quanqr/proto_qr/set"

	"google.golang.org/protobuf/types/known/emptypb"
)

type SetServiceStruct struct {
	setRepo *SetRepository
	set.UnimplementedSetServiceServer
}

func NewSetService(setRepo *SetRepository) *SetServiceStruct {
	return &SetServiceStruct{
		setRepo: setRepo,
	}
}

func (ss *SetServiceStruct) GetSetProtoList(ctx context.Context, _ *emptypb.Empty) (*set.SetProtoListResponse, error) {
	log.Println("Fetching set list")
	sets, err := ss.setRepo.GetSetProtoList(ctx)
	if err != nil {
		log.Println("Error fetching set list:", err)
		return nil, err
	}
	return &set.SetProtoListResponse{
		Data: sets,
	}, nil
}

func (ss *SetServiceStruct) GetSetProtoDetail(ctx context.Context, req *set.SetProtoIdParam) (*set.SetProtoResponse, error) {
	log.Println("Fetching set detail for ID:", req.Id)
	s, err := ss.setRepo.GetSetProtoDetail(ctx, req.Id)
	if err != nil {
		log.Println("Error fetching set detail:", err)
		return nil, err
	}
	return &set.SetProtoResponse{
		Data: s,
	}, nil
}

func (ss *SetServiceStruct) CreateSetProto(ctx context.Context, req *set.CreateSetProtoRequest) (*set.SetProtoResponse, error) {
	log.Println("Creating new set:", req.Name)
	createdSet, err := ss.setRepo.CreateSetProto(ctx, req)
	if err != nil {
		log.Println("Error creating set:", err)
		return nil, err
	}
	log.Println("Set created successfully. ID:", createdSet.Id)
	return &set.SetProtoResponse{
		Data: createdSet,
	}, nil
}

func (ss *SetServiceStruct) UpdateSetProto(ctx context.Context, req *set.UpdateSetProtoRequest) (*set.SetProtoResponse, error) {
	log.Println("Updating set:", req.Id)
	updatedSet, err := ss.setRepo.UpdateSetProto(ctx, req)
	if err != nil {
		log.Println("Error updating set:", err)
		return nil, err
	}
	return &set.SetProtoResponse{
		Data: updatedSet,
	}, nil
}

func (ss *SetServiceStruct) DeleteSetProto(ctx context.Context, req *set.SetProtoIdParam) (*set.SetProtoResponse, error) {
	log.Println("Deleting set:", req.Id)
	deletedSet, err := ss.setRepo.DeleteSetProto(ctx, req.Id)
	if err != nil {
		log.Println("Error deleting set:", err)
		return nil, err
	}
	return &set.SetProtoResponse{
		Data: deletedSet,
	}, nil
}

// type SetServiceStruct struct {
//     setRepo *SetRepository
//     set.UnimplementedSetServiceServer
// }

// func NewSetService(setRepo *SetRepository) *SetServiceStruct {
//     return &SetServiceStruct{
//         setRepo: setRepo,
//     }
// }

// func (ss *SetServiceStruct) GetSetProtoList(ctx context.Context, _ *emptypb.Empty) (*set.SetProtoListResponse, error) {
//     log.Println("Fetching set list")
//     sets, err := ss.setRepo.GetSetProtoList(ctx)
//     if err != nil {
//         log.Println("Error fetching set list:", err)
//         return nil, err
//     }
//     return &set.SetProtoListResponse{
//         Data: sets,
//     }, nil
// }

// func (ss *SetServiceStruct) GetSetProtoDetail(ctx context.Context, req *set.SetProtoIdParam) (*set.SetProtoResponse, error) {
//     log.Println("Fetching set detail for ID:", req.Id)
//     s, err := ss.setRepo.GetSetProtoDetail(ctx, req.Id)
//     if err != nil {
//         log.Println("Error fetching set detail:", err)
//         return nil, err
//     }
//     return &set.SetProtoResponse{
//         Data: s,
//     }, nil
// }

// func (ss *SetServiceStruct) CreateSetProto(ctx context.Context, req *set.CreateSetProtoRequest) (*set.SetProtoResponse, error) {
//     log.Println("Creating new set:", req.Name)
//     createdSet, err := ss.setRepo.CreateSetProto(ctx, req)
//     if err != nil {
//         log.Println("Error creating set:", err)
//         return nil, err
//     }
//     log.Println("Set created successfully. ID:", createdSet.Id)
//     return &set.SetProtoResponse{
//         Data: createdSet,
//     }, nil
// }

// func (ss *SetServiceStruct) UpdateSetProto(ctx context.Context, req *set.UpdateSetProtoRequest) (*set.SetProtoResponse, error) {
//     log.Println("Updating set:", req.Id)
//     updatedSet, err := ss.setRepo.UpdateSetProto(ctx, req)
//     if err != nil {
//         log.Println("Error updating set:", err)
//         return nil, err
//     }
//     return &set.SetProtoResponse{
//         Data: updatedSet,
//     }, nil
// }

// func (ss *SetServiceStruct) DeleteSetProto(ctx context.Context, req *set.SetProtoIdParam) (*set.SetProtoResponse, error) {
//     log.Println("Deleting set:", req.Id)
//     deletedSet, err := ss.setRepo.DeleteSetProto(ctx, req.Id)
//     if err != nil {
//         log.Println("Error deleting set:", err)
//         return nil, err
//     }
//     return &set.SetProtoResponse{
//         Data: deletedSet,
//     }, nil
// }