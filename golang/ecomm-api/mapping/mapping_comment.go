package mapping_user

import (
	"english-ai-full/ecomm-api/types"
	pb "english-ai-full/ecomm-grpc/proto/comment"

)
func ConvertProtoToModelComment(comment *pb.Comment) *types.CommentModel {
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
func convertProtoToModelComments(comments []*pb.Comment) []*types.CommentModel {
	modelComments := make([]*types.CommentModel, len(comments))
	for i, comment := range comments {
		modelComments[i] = ConvertProtoToModelComment(comment)
	}
	return modelComments
}