package service

import (
	"context"
	"log"
	"time"

	proto "english-ai-full/ecomm-grpc/proto/comment"
repository "english-ai-full/ecomm-grpc/repository/comment_repository"
	"english-ai-full/ecomm-api/types"
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

func (cs *CommentServiceStruct) CreateComment(ctx context.Context, req *proto.GetCommentsRequest) (*proto.Comment, error) {
	log.Println("CreateComment")
	newComment := &types.CommentModel{
		Content:   req.Content,
		AuthorID:  req.AuthorId,
		ParentID:  req.ParentId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdComment, err := cs.commentRepo.CreateComment(ctx, newComment.Content, newComment.AuthorID, newComment.ParentID)
	if err != nil {
		log.Println("Error creating comment: proto service", err)
		return nil, err
	}

	log.Println("Comment created successfully. ID:", createdComment.Id)
	return convertPbCommentToCommentRes(createdComment), nil
}

func (cs *CommentServiceStruct) GetComments(ctx context.Context, req *proto.GetCommentsReq) (*proto.GetCommentsRes, error) {
	log.Println("GetComments")
	comments, err := cs.commentRepo.GetComments(ctx, req.ParentId)
	if err != nil {
		log.Println("Error getting comments: proto service", err)
		return nil, err
	}

	return &proto.GetCommentsRes{
		Comments: comments,
	}, nil
}

func (cs *CommentServiceStruct) UpdateComment(ctx context.Context, req *proto.UpdateCommentReq) (*proto.CommentRes, error) {
	log.Println("UpdateComment")
	updatedComment, err := cs.commentRepo.UpdateComment(ctx, req.Id, req.Content)
	if err != nil {
		log.Println("Error updating comment: proto service", err)
		return nil, err
	}

	return convertPbCommentToCommentRes(updatedComment), nil
}

func (cs *CommentServiceStruct) DeleteComment(ctx context.Context, req *proto.DeleteCommentReq) (*emptypb.Empty, error) {
	log.Println("DeleteComment")
	err := cs.commentRepo.DeleteComment(ctx, req.Id)
	if err != nil {
		log.Println("Error deleting comment: proto service", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (cs *CommentServiceStruct) GetCommentByID(ctx context.Context, req *proto.GetCommentByIDReq) (*proto.CommentRes, error) {
	log.Println("GetCommentByID")
	comment, err := cs.commentRepo.GetCommentByID(ctx, req.Id)
	if err != nil {
		log.Println("Error getting comment by ID: proto service", err)
		return nil, err
	}

	return convertPbCommentToCommentRes(comment), nil
}

func convertPbCommentToCommentRes(pbComment *proto.Comment) *proto.CommentRes {
	return &proto.CommentRes{
		Id:        pbComment.Id,
		Content:   pbComment.Content,
		AuthorId:  pbComment.AuthorId,
		ParentId:  pbComment.ParentId,
		CreatedAt: timestamppb.New(pbComment.CreatedAt.AsTime()),
		UpdatedAt: timestamppb.New(pbComment.UpdatedAt.AsTime()),
	}
}