package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"english-ai-full/ecomm-grpc/models"

	pb "english-ai-full/ecomm-grpc/proto"
	"english-ai-full/token"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
    UserClient pb.UserServiceClient
    TokenMaker *token.JWTMaker
}

func NewHandler(client pb.UserServiceClient, secretKey string) *Handler {
    return &Handler{
        UserClient: client,
        TokenMaker: token.NewJWTMaker(secretKey),
    }
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    req := &pb.CreateUserRequest{
        Username: user.Username,
        Email:    user.Email,
    }

    resp, err := h.UserClient.CreateUser(r.Context(), req)
    if err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(resp.User)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    resp, err := h.UserClient.GetUser(r.Context(), &pb.GetUserRequest{Id: id})
    if err != nil {
        http.Error(w, "Failed to get user", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp.User)
}

// Implement other handler methods (UpdateUser, DeleteUser, etc.) as needed

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Implement login logic here
	// This should validate credentials and create a JWT token
}

func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Implement logout logic here
	// This should invalidate the JWT token
}

// Add other handler methods as needed