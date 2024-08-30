package mapping_user

import (
	"english-ai-full/ecomm-api/types"
	pb "english-ai-full/ecomm-grpc/proto"

	// "google.golang.org/protobuf/types/known/timestamppb"
)

// UserReqModel represents the user request model

// ToPBUserReq converts a UserReqModel to a pb.UserReq
func ToPBUserReq(u types.UserReqModel) *pb.UserReq {
	var phone int64
	if u.Phone != nil {
		phone = *u.Phone
	}
	return &pb.UserReq{
		Id:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
		Phone:    phone,
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

// func sessionToSessionRes(s *types.Session) *pb.SessionRes {
//     return &pb.SessionRes{
//         Id:           s.ID,
//         UserEmail:    s.UserEmail,
//         RefreshToken: s.RefreshToken,
//         IsRevoked:    s.IsRevoked,
//         ExpiresAt:    timestamppb.New(s.ExpiresAt),
//     }
// }

// Helper function to convert SessionReq to Session
// func sessionReqToSession(req *pb.SessionReq) *types.Session {
//     return &types.Session{
//         ID:           req.Id,
//         UserEmail:    req.UserEmail,
//         RefreshToken: req.RefreshToken,
//         IsRevoked:    req.IsRevoked,
//         ExpiresAt:    req.ExpiresAt.AsTime(),
//     }
// }

// func ToUserRes(u *pb.UserRes) types.UserResModel {
// 	return types.UserResModel{
// 		ID:        u.Id,
// 		Name:      u.Name,
// 		Email:     u.Email,
// 		IsAdmin:   u.IsAdmin,
// 		Phone:     u.Phone,
// 		Image:     u.Image,
// 		Address:   u.Address,
// 		CreatedAt: u.CreatedAt.AsTime(),
// 		UpdatedAt: u.UpdatedAt.AsTime(),
// 	}
// }