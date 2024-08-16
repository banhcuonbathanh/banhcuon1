package mapping

import (


	"english-ai-full/ecomm-api/types"
	pb "english-ai-full/ecomm-grpc/proto"
)

// UserReqModel represents the user request model

// ToPBUserReq converts a UserReqModel to a pb.UserReq
func ToPBUserReq(u types.UserReqModel) *pb.UserReq {
	return &pb.UserReq{
		Id:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
		Phone:    u.Phone,
		Image:    u.Image,
		Address:  u.Address,
	}
}

// UserRes is the local user response struct


// ToUserRes converts a pb.UserReq to a local UserRes struct
func ToUserRes(u *pb.UserRes) types.UserResModel {
	return types.UserResModel{
		ID:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		Phone:     u.Phone,
		Image:     u.Image,
		Address:   u.Address,
		CreatedAt: u.CreatedAt.AsTime(),
		UpdatedAt: u.UpdatedAt.AsTime(),
	}
}

// ToUserResFromPbUserRes converts a pb.UserRes to a local UserRes struct
func ToUserResFromPbUserRes(u *pb.UserRes) types.UserResModel {
	return types.UserResModel{
		ID:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		Phone:     u.Phone,
		Image:     u.Image,
		Address:   u.Address,
		CreatedAt: u.CreatedAt.AsTime(),
		UpdatedAt: u.UpdatedAt.AsTime(),
	}
}

// RegisterReq is the registration request struct
type RegisterReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ToPBRegisterReq converts a RegisterReq to a pb.UserReq
func ToPBRegisterReq(r RegisterReq) *pb.UserReq {
	return &pb.UserReq{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}

func ToUserResFromUserReq(u *pb.UserReq) types.UserResModel {
	return types.UserResModel{
		ID:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		Phone:     u.Phone,
		Image:     u.Image,
		Address:   u.Address,
		CreatedAt: u.CreatedAt.AsTime(),
		UpdatedAt: u.UpdatedAt.AsTime(),
	}
}