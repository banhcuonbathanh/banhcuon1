package comment_api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"


	"english-ai-full/ecomm-api/types"
	pb "english-ai-full/ecomm-grpc/proto/comment"
	service "english-ai-full/ecomm-grpc/service/comment_service"
	"english-ai-full/token"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CommentHandlerController struct {
	ctx    context.Context
	client pb.CommentServiceClient
	TokenMaker *token.JWTMaker
}

func NewCommentHandler(client pb.CommentServiceClient, secretKey string) *CommentHandlerController {
	return &CommentHandlerController{
		ctx:    context.Background(),
		client: client,
		TokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (h *CommentHandlerController) CreateComment(w http.ResponseWriter, r *http.Request) {
	var req types.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	// Convert types.CreateCommentRequest to pb.CreateCommentRequest
	protoReq := &pb.CreateCommentRequest{
		Content:  req.Content,
		AuthorId: req.AuthorID,
		ParentId: req.ParentID,
	}

	createdComment, err := h.client.CreateComment(h.ctx, protoReq)
	if err != nil {
		http.Error(w, "error creating comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert proto.Comment to types.CommentModel
	responseComment := service.ConvertProtoToModelComment(createdComment)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseComment)
}

func (h *CommentHandlerController) GetComments(w http.ResponseWriter, r *http.Request) {
	parentID := chi.URLParam(r, "parentId")
	parentIDInt, err := strconv.ParseInt(parentID, 10, 64)
	if err != nil {
		http.Error(w, "invalid parent ID", http.StatusBadRequest)
		return
	}

	comments, err := h.client.GetComments(h.ctx, &pb.GetCommentsRequest{ParentId: parentIDInt})
	if err != nil {
		http.Error(w, "error getting comments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandlerController) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var req types.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	// Convert types.UpdateCommentRequest to pb.CommentRes
	protoReq := &pb.UpdateCommentRequest{
		Id:      req.ID,
		Content: req.Content,
	}

	updatedComment, err := h.client.UpdateComment(h.ctx, protoReq)
	if err != nil {
		http.Error(w, "error updating comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert pb.CommentRes to types.CommentModel
	responseComment := service.ConvertProtoToModelComment(updatedComment)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseComment)
}

func (h *CommentHandlerController) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid comment ID", http.StatusBadRequest)
		return
	}

	_, err = h.client.DeleteComment(h.ctx, &pb.DeleteCommentRequest{Id: idInt})
	if err != nil {
		http.Error(w, "error deleting comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CommentHandlerController) GetCommentByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid comment ID", http.StatusBadRequest)
		return
	}

	comment, err := h.client.GetCommentByID(h.ctx, &pb.GetCommentByIDRequest{Id: idInt})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			http.Error(w, "comment not found", http.StatusNotFound)
		} else {
			http.Error(w, "error getting comment: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}