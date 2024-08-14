package service
import (
    "context"
    "english-ai-full/ecomm-grpc/proto"
    "english-ai-full/ecomm-grpc/models"
    "english-ai-full/ecomm-grpc/repository/user_repository"
    "google.golang.org/protobuf/types/known/emptypb"
    "google.golang.org/protobuf/types/known/timestamppb"
)

type UserServericeStruct struct {
    userRepo repository.UserRepositoryInterface
    proto.UnimplementedEcommUserServer
}

func NewUserServer(userRepo repository.UserRepositoryInterface) proto.EcommUserServer {
    return &UserServericeStruct{
        userRepo: userRepo,
    }
}

func (us *UserServericeStruct) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.User, error) {
    newUser := &models.User{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    }
    err := us.userRepo.CreateUser(ctx, newUser)
    if err != nil {
        return nil, err
    }

    return convertModelUserToProtoUser(*newUser), nil

}

func (us *UserServericeStruct) SaveUser(ctx context.Context, req *proto.User) (*emptypb.Empty, error) {
    user := convertProtoUserToModelUser(req)
    err := us.userRepo.Save(user)
    return &emptypb.Empty{}, err
}

func (us *UserServericeStruct) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.User, error) {
    updatedUser := models.User{
        ID:       req.Id,
        Username: req.Username,
        Email:    req.Email,
    }
    user, err := us.userRepo.Update(updatedUser)
    if err != nil {
        return nil, err
    }
    return convertModelUserToProtoUser(user), nil

}

func (us *UserServericeStruct) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*emptypb.Empty, error) {
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

func (us *UserServericeStruct) FindByEmail(ctx context.Context, req *proto.FindByEmailRequest) (*proto.User, error) {
    user, err := us.userRepo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    return convertModelUserToProtoUser(*user), nil

}

func (us *UserServericeStruct) FindUsersByPage(ctx context.Context, req *proto.PageRequest) (*proto.UserList, error) {
    users, err := us.userRepo.FindUsersByPage(int(req.PageNumber), int(req.PageSize))
    if err != nil {
        return nil, err
    }
    return &proto.UserList{Users: convertModelUsersToProtoUsers(users)}, nil
}

func (us *UserServericeStruct) Login(ctx context.Context, req *proto.LoginRequest) (*proto.User, error) {
    user, err := us.userRepo.Login(ctx, req.Email, req.Password)
    if err != nil {
        return nil, err
    }
    return convertModelUserToProtoUser(*user), nil

}

func (us *UserServericeStruct) Register(ctx context.Context, req *proto.CreateUserRequest) (*proto.RegisterResponse, error) {
    newUser := models.User{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    }
    registeredUser, err := us.userRepo.Register(newUser)
    if err != nil {
        return &proto.RegisterResponse{Success: false}, err
    }
    
    // Convert the registered user to a proto User
    protoUser := convertModelUserToProtoUser(registeredUser)
    
    // Return a successful response with the registered user
    return &proto.RegisterResponse{
        Success: true,
        User: protoUser,
    }, nil
}

// Helper functions to convert between model and proto user types
func convertModelUserToProtoUser(user models.User) *proto.User {
    return &proto.User{
        Id:        user.ID,
        Username:  user.Username,
        Email:     user.Email,
        Password:  user.Password,
        CreatedAt: timestamppb.New(user.CreatedAt),
        UpdatedAt: timestamppb.New(user.UpdatedAt),
    }
}
func convertProtoUserToModelUser(user *proto.User) models.User {
    return models.User{
        ID:        user.Id,
        Username:  user.Username,
        Email:     user.Email,
        Password:  user.Password,
        CreatedAt: user.CreatedAt.AsTime(),
        UpdatedAt: user.UpdatedAt.AsTime(),
    }
}

func convertModelUsersToProtoUsers(users []models.User) []*proto.User {
    protoUsers := make([]*proto.User, len(users))
    for i, user := range users {
        protoUsers[i] = convertModelUserToProtoUser(user)
    }
    return protoUsers
}