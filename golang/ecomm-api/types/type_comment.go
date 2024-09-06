package types


import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// CommentModel represents the structure for a comment
type CommentModel struct {
	ID        int64          `json:"id"`
	Content   string         `json:"content"`
	AuthorID  int64          `json:"author_id"`
	ParentID  int64          `json:"parent_id"`
	Replies   []*CommentModel `json:"replies,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// CreateCommentRequest represents the structure for creating a new comment
type CreateCommentRequest struct {
	Content  string `json:"content"`
	AuthorID int64  `json:"author_id"`
	ParentID int64  `json:"parent_id"`
}

// GetCommentsRequest represents the structure for requesting comments
type GetCommentsRequest struct {
	ParentID int64 `json:"parent_id"`
}

// GetCommentsResponse represents the structure for the response to a get comments request
type GetCommentsResponse struct {
	Comments []*CommentModel `json:"comments"`
}

// CreateReplyRequest represents the structure for creating a new reply
type CreateReplyRequest struct {
	Content  string `json:"content"`
	AuthorID int64  `json:"author_id"`
	ParentID int64  `json:"parent_id"`
}

// UpdateCommentRequest represents the structure for updating a comment
type UpdateCommentRequest struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

// GetCommentByIDRequest represents the structure for requesting a comment by ID
type GetCommentByIDRequest struct {
	ID int64 `json:"id"`
}

// DeleteCommentRequest represents the structure for deleting a comment
type DeleteCommentRequest struct {
	ID int64 `json:"id"`
}

// DeleteCommentResponse represents the structure for the response to a delete comment request
type DeleteCommentResponse struct {
	Success bool `json:"success"`
}

// ConvertToProtoComment converts a CommentModel to a proto Comment
func (c *CommentModel) ConvertToProtoComment() *Comment {
	return &Comment{
		Id:        c.ID,
		Content:   c.Content,
		AuthorId:  c.AuthorID,
		ParentId:  c.ParentID,
		Replies:   convertToProtoComments(c.Replies),
		CreatedAt: timestamppb.New(c.CreatedAt),
		UpdatedAt: timestamppb.New(c.UpdatedAt),
	}
}

// ConvertFromProtoComment converts a proto Comment to a CommentModel
func ConvertFromProtoComment(pc *Comment) *CommentModel {
	return &CommentModel{
		ID:        pc.Id,
		Content:   pc.Content,
		AuthorID:  pc.AuthorId,
		ParentID:  pc.ParentId,
		Replies:   convertFromProtoComments(pc.Replies),
		CreatedAt: pc.CreatedAt.AsTime(),
		UpdatedAt: pc.UpdatedAt.AsTime(),
	}
}

// Helper functions for converting slices of comments
func convertToProtoComments(comments []*CommentModel) []*Comment {
	protoComments := make([]*Comment, len(comments))
	for i, comment := range comments {
		protoComments[i] = comment.ConvertToProtoComment()
	}
	return protoComments
}

func convertFromProtoComments(protoComments []*Comment) []*CommentModel {
	comments := make([]*CommentModel, len(protoComments))
	for i, pc := range protoComments {
		comments[i] = ConvertFromProtoComment(pc)
	}
	return comments
}