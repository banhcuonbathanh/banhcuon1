package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	middleware "english-ai-full/ecomm-api"
	mapping_user "english-ai-full/ecomm-api/mapping"
	"english-ai-full/ecomm-api/types"
	pb "english-ai-full/ecomm-grpc/proto"
	logg "english-ai-full/logger"
	"english-ai-full/token"

	// "english-ai-full/util"

	// "english-ai-full/util"

	"github.com/go-chi/chi"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handlercontroller struct {
    ctx        context.Context
    client     pb.EcommUserClient 
    TokenMaker *token.JWTMaker
	logger   *logg.Logger
 
}

func NewHandler(client pb.EcommUserClient, secretKey string) *Handlercontroller {
    return &Handlercontroller{
        ctx:        context.Background(),
        client:     client,
        TokenMaker: token.NewJWTMaker(secretKey),
		logger:   logg.NewLogger(),
     
    }
}

func (h *Handlercontroller) CreateUser(w http.ResponseWriter, r *http.Request) {
    h.logger.Info("create users handler started")
    var u types.UserReqModel
    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        h.logger.Error(fmt.Sprintf("error decoding request body: %v", err))
        http.Error(w, "error decoding request body", http.StatusBadRequest)
        return
    }

    h.logger.Info("calling CreateUser from gRPC client")
    user, err := h.client.CreateUser(h.ctx, mapping_user.ToPBUserReq(u))
    if err != nil {
        h.logger.Error(fmt.Sprintf("error creating user: %v", err))
        http.Error(w, "error creating user", http.StatusInternalServerError)
        return
    }

    res := mapping_user.ToUserResFromUserReq(user)
    h.logger.Info("user successfully created")

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(res); err != nil {
        h.logger.Error(fmt.Sprintf("error encoding response: %v", err))
        http.Error(w, "error encoding response", http.StatusInternalServerError)
    }
}


func (h *Handlercontroller) FindByEmail(w http.ResponseWriter, r *http.Request) {
	log.Println("User FindByEmail handlercontroller")
	email := chi.URLParam(r, "email")
	user, err := h.client.FindByEmail(h.ctx, &pb.UserReq{Email: email})
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}


   res := mapping_user.ToUserResFromPbUserRes(user)
	log.Println("User FindByEmail handlercontroller res", res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handlercontroller) ListUsers(w http.ResponseWriter, r *http.Request) {
	lur, err := h.client.FindAllUsers(h.ctx, &emptypb.Empty{})
	if err != nil {
		http.Error(w, "failed to fetch users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var res []types.UserResModel
	for _, u := range lur.GetUsers() {
		userRes := mapping_user.ToUserRes(u)
		res = append(res, userRes)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}



func (h *Handlercontroller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing ID", http.StatusBadRequest)
		return
	}
	_, err = h.client.DeleteUser(h.ctx, &pb.UserReq{Id: i})
	if err != nil {
		http.Error(w, "error deleting user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Add these functions to handle login and register
func (h *Handlercontroller) Login(w http.ResponseWriter, r *http.Request) {
    h.logger.Info("user login handler started")
    var u types.LoginUserReq

    if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
        h.logger.Error(fmt.Sprintf("error decoding request body: %v", err))
        http.Error(w, "error decoding request body", http.StatusBadRequest)
        return
    }

    ur, err := h.client.FindByEmail(h.ctx, &pb.UserReq{Email: u.Email})
    if err != nil {
        h.logger.Error(fmt.Sprintf("error getting user by email: %v", err))
        http.Error(w, "error getting user", http.StatusInternalServerError)
        return
    }

    h.logger.Info("user found, proceeding to create tokens")

	// err = util.CheckPassword(u.Password, ur.Password)
	// if err != nil {
	// 	http.Error(w, "wrong password", http.StatusUnauthorized)
	// 	return
	// }

	// create a json web token (JWT) and return it as response
	accessToken, accessClaims, err := h.TokenMaker.CreateToken(ur.GetId(), ur.GetEmail(), ur.GetRole(), 15*time.Minute)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}
	log.Println("User login handlercontroller accessToken", accessToken)
	refreshToken, refreshClaims, err := h.TokenMaker.CreateToken(ur.GetId(), ur.GetEmail(), ur.GetRole(), 24*time.Hour)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}
	log.Println("User login handlercontroller refreshToken", refreshToken)
	log.Println("User login handlercontroller refreshClaims.RegisteredClaims.ID", refreshClaims.RegisteredClaims.ID)
	log.Println("User login handlercontroller ur.Email", ur.Email)
	log.Println("User login handlercontroller timestamppb.New(refreshClaims.RegisteredClaims.ExpiresAt.Time", timestamppb.New(refreshClaims.RegisteredClaims.ExpiresAt.Time))


	if len(refreshClaims.RegisteredClaims.ID) > 255 {
        log.Printf("Warning: ID exceeds 255 characters")
    }
    if len(refreshToken) > 255 {
        log.Printf("Warning: Refresh token exceeds 255 characters")
    }
	session, err := h.client.CreateSession(h.ctx, &pb.SessionReq{
		Id:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    ur.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    timestamppb.New(refreshClaims.RegisteredClaims.ExpiresAt.Time),
	})
	if err != nil {
		http.Error(w, "error creating session", http.StatusInternalServerError)
		return
	}
	log.Println("User login handlercontroller session", session)
	res := types.LoginUserRes{
		SessionID:             session.GetId(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		User:                  mapping_user.ToUserRes(ur),
	}
	log.Println("User login handlercontroller res", res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}


func (h *Handlercontroller) LogoutUser(w http.ResponseWriter, r *http.Request) {
    h.logger.Info("logout handler started")

    claims, ok := r.Context().Value(middleware.AuthKey{}).(*token.UserClaims)
    if !ok {
        h.logger.Error("failed to retrieve user claims from context")
        http.Error(w, "failed to retrieve user claims", http.StatusUnauthorized)
        return
    }

    // Attempt to delete the session
    h.logger.Info(fmt.Sprintf("attempting to delete session with ID: %s", claims.RegisteredClaims.ID))
    _, err := h.client.DeleteSession(h.ctx, &pb.SessionReq{
        Id: claims.RegisteredClaims.ID,
    })
    if err != nil {
        h.logger.Error(fmt.Sprintf("error deleting session: %v", err))
        http.Error(w, "error deleting session", http.StatusInternalServerError)
        return
    }

    // Successfully deleted session
    h.logger.Info(fmt.Sprintf("session with ID %s successfully deleted", claims.RegisteredClaims.ID))
    w.WriteHeader(http.StatusNoContent)
}


func (h *Handlercontroller) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	var req types.RenewAccessTokenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	refreshClaims, err := h.TokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "error verifying token", http.StatusUnauthorized)
		return
	}

	session, err := h.client.GetSession(h.ctx, &pb.SessionReq{
		Id: refreshClaims.RegisteredClaims.ID,
	})
	if err != nil {
		http.Error(w, "error getting session", http.StatusInternalServerError)
		return
	}

	if session.IsRevoked {
		http.Error(w, "session revoked", http.StatusUnauthorized)
		return
	}

	if session.GetUserEmail() != refreshClaims.Email {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	accessToken, accessClaims, err := h.TokenMaker.CreateToken(refreshClaims.ID, refreshClaims.Email, refreshClaims.Role, 15*time.Minute)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	res := types.RenewAccessTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *Handlercontroller) RevokeSession(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.AuthKey{}).(*token.UserClaims)
	
	_, err := h.client.RevokeSession(h.ctx, &pb.SessionReq{
		Id: claims.RegisteredClaims.ID,
	})
	if err != nil {
		http.Error(w, "error revoking session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}