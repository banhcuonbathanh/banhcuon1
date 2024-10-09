package service

import (
	"context"
	"english-ai-full/ecomm-grpc/proto"

	"time"

	repository "english-ai-full/ecomm-grpc/repository/user_repository"

	"english-ai-full/ecomm-api/types"
	logg "english-ai-full/logger"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServericeStruct struct {
    userRepo *repository.UserRepository
    logger   *logg.Logger
    proto.UnimplementedEcommUserServer
}

func NewUserServer(userRepo *repository.UserRepository) *UserServericeStruct {
    return &UserServericeStruct{
        userRepo: userRepo,
        logger:   logg.NewLogger(),
    }
}

func (us *UserServericeStruct) CreateUser(ctx context.Context, req *proto.UserReq) (*proto.UserReq, error) {
    us.logger.Info("Creating new user: " +
        "Name: " + req.Name +
        ", Email: " + req.Email +
        ", Role: " + req.Role +
        ", Phone: " + req.Phone +
        ", Image: " + req.Image +
        ", Address: " + req.Address +
        ", CreatedAt: " + time.Now().String() +
        ", UpdatedAt: " + time.Now().String())

    newUser := &types.UserReqModel{
        Name:      req.Name,
        Email:     req.Email,
        Password:  req.Password,
        Role:      req.Role,
        Phone:     req.Phone,
        Image:     req.Image,
        Address:   req.Address,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    createdUser, err := us.userRepo.CreateUser(ctx, newUser)
    if err != nil {
        us.logger.Error("Error creating user: " + err.Error())
        return nil, err
    }

    req.Id = createdUser.ID

    us.logger.Info("User created successfully. ID: " )
    return req, nil
}

func (us *UserServericeStruct) SaveUser(ctx context.Context, req *proto.UserReq) (*emptypb.Empty, error) {
    user := convertProtoUserReqToModelUser(req)
    err := us.userRepo.Save(user)
    if err != nil {
        us.logger.Error("Error saving user: " + err.Error())
    } else {
        us.logger.Info("User saved successfully. ID: " )
    }
    return &emptypb.Empty{}, err
}

func (us *UserServericeStruct) UpdateUser(ctx context.Context, req *proto.UserReq) (*proto.UserReq, error) {
    updatedUser := types.UserReqModel{
        ID:        req.Id,
        Name:      req.Name,
        Email:     req.Email,
        Password:  req.Password,
        Role:      req.Role,
        Phone:     req.Phone,
        Image:     req.Image,
        Address:   req.Address,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    _, err := us.userRepo.Update(updatedUser)
    if err != nil {
        us.logger.Error("Error updating user: " + err.Error())
        return nil, err
    }
    us.logger.Info("User updated successfully. ID: " )
    return req, nil
}

func (us *UserServericeStruct) DeleteUser(ctx context.Context, req *proto.UserReq) (*emptypb.Empty, error) {
    err := us.userRepo.Delete(int(req.Id))
    if err != nil {
        us.logger.Error("Error deleting user: " + err.Error())
    } else {
        us.logger.Info("User deleted successfully. ID: " )
    }
    return &emptypb.Empty{}, err
}

func (us *UserServericeStruct) FindAllUsers(ctx context.Context, _ *emptypb.Empty) (*proto.UserList, error) {
    users, err := us.userRepo.FindAll()
    if err != nil {
        us.logger.Error("Error finding all users: " + err.Error())
        return nil, err
    }
    if len(users) > 0 {
        lastUser := users[len(users)-1]
        us.logger.Info("User find all service: Last user - ID: "  +
            ", Name: " + lastUser.Name +
            ", Email: " + lastUser.Email +
            ", Role: " + lastUser.Role +
            ", Phone: " + lastUser.Phone +
            ", Image: " + lastUser.Image +
            ", Address: " + lastUser.Address +
            ", CreatedAt: " + lastUser.CreatedAt.String() +
            ", UpdatedAt: " + lastUser.UpdatedAt.String())
    }
    return &proto.UserList{Users: convertModelUsersToProtoUsers(users)}, nil
}

func (us *UserServericeStruct) FindByEmail(ctx context.Context, req *proto.UserReq) (*proto.UserRes, error) {
    user, err := us.userRepo.FindByEmail(req.Email)
    if err != nil {
        us.logger.Error("Error finding user by email: " + err.Error())
        return nil, err
    }
    us.logger.Info("User found by email: " + req.Email)
    return convertModelUserToProtoUserRes(user), nil
}

func (us *UserServericeStruct) FindUsersByPage(ctx context.Context, req *proto.PageRequest) (*proto.UserList, error) {
    users, err := us.userRepo.FindUsersByPage(int(req.PageNumber), int(req.PageSize))
    if err != nil {
        us.logger.Error("Error finding users by page: " + err.Error())
        return nil, err
    }
    us.logger.Info("Found users by page: Page " + string(req.PageNumber) + ", Size " + string(req.PageSize))
    return &proto.UserList{Users: convertModelUsersToProtoUsers(users)}, nil
}

func (us *UserServericeStruct) Login(ctx context.Context, req *proto.LoginRequest) (*proto.UserReq, error) {
    user, err := us.userRepo.Login(ctx, req.Email, req.Password)
    if err != nil {
        us.logger.Error("Login failed: " + err.Error())
        return nil, err
    }
    us.logger.Info("User logged in successfully: " + req.Email)
    return convertModelUserToProtoUserReq(user), nil
}

func convertModelUserToProtoUserReq(user *types.UserReqModel) *proto.UserReq {
    return &proto.UserReq{
        Id:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Password:  user.Password,
        Role:      user.Role,
        Phone:     user.Phone,
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
        Phone:     user.Phone,
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
        Phone:    req.Phone,
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
        Phone:     user.Phone,
        Image:     user.Image,
        Address:   user.Address,
    }
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