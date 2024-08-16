package service

import (
	"context"
	"english-ai-full/ecomm-grpc/proto"
	"time"

	"english-ai-full/ecomm-grpc/repository/user_repository"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"english-ai-full/ecomm-api/types"
)

type UserServericeStruct struct {
    userRepo *repository.UserRepository
    proto.UnimplementedEcommUserServer
}

func NewUserServer(userRepo *repository.UserRepository) *UserServericeStruct {
    return &UserServericeStruct{
        userRepo: userRepo,
    }
}

func (us *UserServericeStruct) CreateUser(ctx context.Context, req *proto.UserReq) (*proto.UserReq, error) {
    newUser := &types.UserReqModel{
        ID:        req.Id,
        Name:      req.Name,
        Email:     req.Email,
        Password:  req.Password,
        IsAdmin:   req.IsAdmin,
        Phone:     req.Phone,
        Image:     req.Image,
        Address:   req.Address,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    err := us.userRepo.CreateUser(ctx, newUser)
    if err != nil {
        return nil, err
    }

    return req, nil // Return the input UserReq instead of UserRes
}


func (us *UserServericeStruct) SaveUser(ctx context.Context, req *proto.UserReq) (*emptypb.Empty, error) {
    user := convertProtoUserReqToModelUser(req)
    err := us.userRepo.Save(user)
    return &emptypb.Empty{}, err
}

func (us *UserServericeStruct) UpdateUser(ctx context.Context, req *proto.UserReq) (*proto.UserReq, error) {
    updatedUser := types.UserReqModel{
        ID:        req.Id,
        Name:      req.Name,
        Email:     req.Email,
        Password:  req.Password,
        IsAdmin:   req.IsAdmin,
        Phone:     req.Phone,
        Image:     req.Image,
        Address:   req.Address,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    _, err := us.userRepo.Update(updatedUser)
    if err != nil {
        return nil, err
    }
    return req, nil // Return the input UserReq
}


func (us *UserServericeStruct) DeleteUser(ctx context.Context, req *proto.UserReq) (*emptypb.Empty, error) {
    err := us.userRepo.Delete(int(req.Id))
    return &emptypb.Empty{}, err
}

func (us *UserServericeStruct) FindAllUsers(ctx context.Context, _ *emptypb.Empty) (*proto.UserList, error) {
    users, err := us.userRepo.FindAll()
    if err != nil {
        return nil, err
    }
    return &proto.UserList{Users: convertModelUsersToProtoUsers(users)}, nil
}

func (us *UserServericeStruct) FindByEmail(ctx context.Context, req *proto.UserReq) (*proto.UserReq, error) {
    user, err := us.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    return convertModelUserToProtoUserReq(*user), nil
}

func (us *UserServericeStruct) FindUsersByPage(ctx context.Context, req *proto.PageRequest) (*proto.UserList, error) {
    users, err := us.userRepo.FindUsersByPage(int(req.PageNumber), int(req.PageSize))
    if err != nil {
        return nil, err
    }
    return &proto.UserList{Users: convertModelUsersToProtoUsers(users)}, nil
}

func (us *UserServericeStruct) Login(ctx context.Context, req *proto.LoginRequest) (*proto.UserReq, error) {
    user, err := us.userRepo.Login(ctx, req.Email, req.Password)
    if err != nil {
        return nil, err
    }
    return convertModelUserToProtoUserReq(*user), nil
}

// func (us *UserServericeStruct) Register(ctx context.Context, req *proto.UserReq) (*proto.RegisterResponse, error) {
//     // Implement registration logic here
//     // For example:
//     newUser := convertProtoUserReqToModelUser(req)
//     err := us.userRepo.CreateUser(ctx, &newUser)
//     if err != nil {
//         return nil, err
//     }
//     return &proto.RegisterResponse{Success: true, Message: "User registered successfully"}, nil
// }
// Helper functions to convert between model and proto user types
func convertModelUserToProtoUser(user types.UserReqModel) *proto.UserRes {
    return &proto.UserRes{
        Id:        user.ID,
        Name:      user.Name, // Assuming you meant to map Username to Name
        Email:     user.Email,
        Password:  user.Password,
        CreatedAt: timestamppb.New(user.CreatedAt),
        UpdatedAt: timestamppb.New(user.UpdatedAt),
    }
}

// func convertProtoUserToModelUser(user *proto.UserRes) types.UserReqModel {
//     return types.UserReqModel{
//         ID:        user.Id,
//         Name:  user.Name,
//         Email:     user.Email,
//         Password:  user.Password,
//         CreatedAt: user.CreatedAt.AsTime(),
//         UpdatedAt: user.UpdatedAt.AsTime(),
//     }
// }

func convertModelUsersToProtoUsers(users []types.UserReqModel) []*proto.UserRes {
    protoUsers := make([]*proto.UserRes, len(users))
    for i, user := range users {
        protoUsers[i] = convertModelUserToProtoUser(user)
    }
    return protoUsers
}

func convertModelUserToProtoUserReq(user types.UserReqModel) *proto.UserReq {
    return &proto.UserReq{
        Id:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Password:  user.Password,
        IsAdmin:   user.IsAdmin,
        Phone:     user.Phone,
        Image:     user.Image,
        Address:   user.Address,
    }
}

func convertProtoUserReqToModelUser(user *proto.UserReq) types.UserReqModel {
    return types.UserReqModel{
        ID:        user.Id,
        Name:      user.Name,
        Email:     user.Email,
        Password:  user.Password,
        IsAdmin:   user.IsAdmin,
        Phone:     user.Phone,
        Image:     user.Image,
        Address:   user.Address,
    }
}