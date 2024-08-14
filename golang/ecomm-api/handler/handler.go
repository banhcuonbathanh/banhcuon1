package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	pb "english-ai-full/ecomm-grpc/proto"
	"english-ai-full/ecomm-grpc/data/user"
	"english-ai-full/token"

	"github.com/go-chi/chi"
)

type handlercontroller struct {
	ctx        context.Context
	client     pb.EcommUserClient 
	TokenMaker *token.JWTMaker
}

func NewHandler(client pb.EcommUserClient, secretKey string) *handlercontroller {
	return &handlercontroller{
		ctx:        context.Background(),
		client:     client,
		TokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (h *handlercontroller) createUser(w http.ResponseWriter, r *http.Request) {
	var u user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	user, err := h.client.CreateUser(h.ctx, toPBUserReq(u))
	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}
	res := toUserRes(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *handlercontroller) getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing ID", http.StatusBadRequest)
		return
	}
	user, err := h.client.GetUser(h.ctx, &pb.UserReq{Id: i})
	if err != nil {
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
	res := toUserRes(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handlercontroller) listUsers(w http.ResponseWriter, r *http.Request) {
	lur, err := h.client.ListUsers(h.ctx, &pb.UserReq{})
	if err != nil {
		http.Error(w, "error listing users", http.StatusInternalServerError)
		return
	}
	var res []UserRes
	for _, u := range lur.GetUsers() {
		res = append(res, toUserRes(u))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handlercontroller) updateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing ID", http.StatusBadRequest)
		return
	}
	var u UserReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	u.ID = i
	updated, err := h.client.UpdateUser(h.ctx, toPBUserReq(u))
	if err != nil {
		http.Error(w, "error updating user", http.StatusInternalServerError)
		return
	}
	res := toUserRes(updated)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handlercontroller) deleteUser(w http.ResponseWriter, r *http.Request) {
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
func (h *handlercontroller) login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginReq
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	user, err := h.client.Login(h.ctx, &pb.LoginReq{Email: loginReq.Email, Password: loginReq.Password})
	if err != nil {
		http.Error(w, "error logging in", http.StatusUnauthorized)
		return
	}
	res := toUserRes(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *handlercontroller) register(w http.ResponseWriter, r *http.Request) {
	var registerReq RegisterReq
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	success, err := h.client.Register(h.ctx, toPBRegisterReq(registerReq))
	if err != nil {
		http.Error(w, "error registering user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]bool{"success": success})
}

// You'll need to implement these helper functions:
// toPBUserReq, toUserRes, toPBRegisterReq

// Also, define these structs:
// UserReq, UserRes, LoginReq, RegisterReq