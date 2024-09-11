package service

import (
	"context"
	"log"

	"english-ai-full/ecomm-api/types"
	proto "english-ai-full/ecomm-grpc/proto/comment"
	repository "english-ai-full/ecomm-grpc/repository/comment_repository"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CommentServiceStruct struct {
	commentRepo *repository.CommentRepository
	proto.UnimplementedCommentServiceServer
}

func NewCommentServer(commentRepo *repository.CommentRepository) *CommentServiceStruct {
	return &CommentServiceStruct{
		commentRepo: commentRepo,
	}
}


func (cs *CommentServiceStruct) GetComments(ctx context.Context, req *proto.GetCommentsRequest) (*proto.CommentRes, error) {
	log.Println("GetComments")
	comments, err := cs.commentRepo.GetCommentByID(ctx, req.ParentId)
	if err != nil {
		log.Println("Error getting comments: proto service", err)
		return nil, err
	}

	return &proto.CommentRes{
		Id:        comments.Id,
		Content:   comments.Content,
		AuthorId:  comments.AuthorId,
		ParentId:  comments.ParentId,
		CreatedAt: timestamppb.New(comments.CreatedAt.AsTime()),
		UpdatedAt: timestamppb.New(comments.UpdatedAt.AsTime()),
	}, nil
}

func (cs *CommentServiceStruct) UpdateComment(ctx context.Context, req *proto.UpdateCommentRequest) (*proto.CommentRes, error) {
	log.Println("UpdateComment")
	updatedComment, err := cs.commentRepo.UpdateComment(ctx, req.Id, req.Content)
	if err != nil {
		log.Println("Error updating comment: proto service", err)
		return nil, err
	}

	return updatedComment, nil
}

func (cs *CommentServiceStruct) DeleteComment(ctx context.Context, req *proto.DeleteCommentRequest) (*emptypb.Empty, error) {
	log.Println("DeleteComment")
	err := cs.commentRepo.DeleteComment(ctx, req.Id)
	if err != nil {
		log.Println("Error deleting comment: proto service", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (cs *CommentServiceStruct) GetCommentByID(ctx context.Context, req *proto.GetCommentByIDRequest) (*proto.CommentRes, error) {
	log.Println("GetCommentByID")
	comment, err := cs.commentRepo.GetCommentByID(ctx, req.Id)
	if err != nil {
		log.Println("Error getting comment by ID: proto service", err)
		return nil, err
	}

	return comment, nil
}



func (cs *CommentServiceStruct) CreateComment(ctx context.Context, req *proto.CreateCommentRequest) (*proto.CommentRes, error) {
	log.Println("CreateComment")

	createdComment, err := cs.commentRepo.CreateComment(ctx, req.Content, req.AuthorId, req.ParentId)
	if err != nil {
		log.Println("Error creating comment: proto service", err)
		return nil, err
	}

	log.Println("Comment created successfully. ID:", createdComment.Id)
	return createdComment, nil
}

// This function may still be useful for other parts of your code
func ConvertProtoToModelComment(comment *proto.CommentRes) *types.CommentModel {
	return &types.CommentModel{
		ID:        comment.Id,
		Content:   comment.Content,
		AuthorID:  comment.AuthorId,
		ParentID:  comment.ParentId,
		Replies:   convertProtoToModelComments(comment.Replies),
		CreatedAt: comment.CreatedAt.AsTime(),
		UpdatedAt: comment.UpdatedAt.AsTime(),
	}
}

// Helper function to convert a slice of proto Comment to a slice of CommentModel
func convertProtoToModelComments(comments []*proto.CommentRes) []*types.CommentModel {
	modelComments := make([]*types.CommentModel, len(comments))
	for i, comment := range comments {
		modelComments[i] = ConvertProtoToModelComment(comment)
	}
	return modelComments
}

