package service

import (
	"context"
	"english-ai-full/ecomm-grpc/proto"
	"log"
	"time"

	repository "english-ai-full/ecomm-grpc/repository/user_repository"

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
    log.Println("Creating new user:",
    "Name:", req.Name,
    "Email:", req.Email,
    "Password:", req.Password,
    "Role:", req.Role,
    "Phone:", req.Phone,
    "Image:", req.Image,
    "Address:", req.Address,
    "CreatedAt:", time.Now(),
    "UpdatedAt:", time.Now(),
)

newUser := &types.UserReqModel{
    Name:      req.Name,
    Email:     req.Email,
    Password:  req.Password,
    Role:      req.Role,
    Phone:     &req.Phone,
    Image:     req.Image,
    Address:   req.Address,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}

	createdUser, err := us.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		log.Println("Error creating user:", err)
		return nil, err
	}

	req.Id = createdUser.ID

	log.Println("User created successfully. ID:", req.Id)
	return req, nil
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
        Role:      req.Role,
        Phone:     &req.Phone,
        Image:     req.Image,
        Address:   req.Address,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    _, err := us.userRepo.Update(updatedUser)
    if err != nil {
        return nil, err
    }
    return req, nil
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
    if len(users) > 0 {
        lastUser := users[len(users)-1]
        log.Println("User find all service", lastUser.ID, 
            "Name:", lastUser.Name, 
            "Email:", lastUser.Email, 
            "Password:", lastUser.Password, 
            "Role:", lastUser.Role, 
            "Phone:", lastUser.Phone, 
            "Image:", lastUser.Image, 
            "Address:", lastUser.Address, 
            "CreatedAt:", lastUser.CreatedAt, 
            "UpdatedAt:", lastUser.UpdatedAt)
    }
    return &proto.UserList{Users: convertModelUsersToProtoUsers(users)}, nil
}

func (us *UserServericeStruct) FindByEmail(ctx context.Context, req *proto.UserReq) (*proto.UserRes, error) {
    user, err := us.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    return convertModelUserToProtoUserRes(user), nil
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
    return convertModelUserToProtoUserReq(user), nil
}
func convertModelUserToProtoUserReq(user *types.UserReqModel) *proto.UserReq {
    return &proto.UserReq{
        Id:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Password:  user.Password,
        Role:      user.Role,
        Phone:     getPhoneValue(user.Phone),
        Image:     user.Image,
        Address:   user.Address,
        CreatedAt: timestamppb.New(user.CreatedAt),
        UpdatedAt: timestamppb.New(user.UpdatedAt),
    }
}


func convertModelUserToProtoUserRes(user *types.UserReqModel) *proto.UserRes {
    return &proto.UserRes{
        Id:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Password:  user.Password,
        Role:      user.Role,
        Phone:     getPhoneValue(user.Phone),
        Image:     user.Image,
        Address:   user.Address,
        CreatedAt: timestamppb.New(user.CreatedAt),
        UpdatedAt: timestamppb.New(user.UpdatedAt),
    }
}

func convertModelUsersToProtoUsers(users []types.UserReqModel) []*proto.UserRes {
    protoUsers := make([]*proto.UserRes, len(users))
    for i, user := range users {
        protoUsers[i] = convertModelUserToProtoUserRes(&user) // Pass the pointer to the user variable
    }
    return protoUsers
}


func ConvertProtoToUserReqModel(req *proto.UserReq) *types.UserReqModel {
    return &types.UserReqModel{
        ID:       req.Id,
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
        Role:     req.Role,
        Phone:    &req.Phone,
        Image:    req.Image,
        Address:  req.Address,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}
func convertProtoUserReqToModelUser(user *proto.UserReq) types.UserReqModel {
    return types.UserReqModel{
        ID:        user.Id,
        Name:      user.Name,
        Email:     user.Email,
        Password:  user.Password,
        Role:      user.Role,
        Phone:     &user.Phone,
        Image:     user.Image,
        Address:   user.Address,
    }
}
func getPhoneValue(phone *int64) int64 {
    if phone != nil {
        return *phone
    }
    return 0
}



func (us *UserServericeStruct) CreateSession(ctx context.Context, sr *proto.SessionReq) (*proto.SessionRes, error) {
	sess, err := us.userRepo.CreateSession(ctx, &types.Session{
		ID:           sr.GetId(),
		UserEmail:    sr.GetUserEmail(),
		RefreshToken: sr.GetRefreshToken(),
		IsRevoked:    sr.GetIsRevoked(),
		ExpiresAt:    sr.GetExpiresAt().AsTime(),
	})
	if err != nil {
		return nil, err
	}

	return &proto.SessionRes{
		Id:           sess.ID,
		UserEmail:    sess.UserEmail,
		RefreshToken: sess.RefreshToken,
		IsRevoked:    sess.IsRevoked,
		ExpiresAt:    timestamppb.New(sess.ExpiresAt),
	}, nil
}

func (us *UserServericeStruct) GetSession(ctx context.Context, sr *proto.SessionReq) (*proto.SessionRes, error) {
	sess, err := us.userRepo.GetSession(ctx, sr.GetId())
	if err != nil {
		return nil, err
	}

	return &proto.SessionRes{
		Id:           sess.ID,
		UserEmail:    sess.UserEmail,
		RefreshToken: sess.RefreshToken,
		IsRevoked:    sess.IsRevoked,
		ExpiresAt:    timestamppb.New(sess.ExpiresAt),
	}, nil
}

func (us *UserServericeStruct) RevokeSession(ctx context.Context, sr *proto.SessionReq) (*proto.SessionRes, error) {
	err := us.userRepo.RevokeSession(ctx, sr.GetId())
	if err != nil {
		return nil, err
	}

	return &proto.SessionRes{}, nil
}

func (us *UserServericeStruct) DeleteSession(ctx context.Context, sr *proto.SessionReq) (*proto.SessionRes, error) {
	err := us.userRepo.DeleteSession(ctx, sr.GetId())
	if err != nil {
		return nil, err
	}

	return &proto.SessionRes{}, nil
}