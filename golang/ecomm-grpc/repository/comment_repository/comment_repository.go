// repository/comment_repository.go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "english-ai-full/ecomm-grpc/proto/comment"
)

type CommentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (r *CommentRepository) CreateComment(ctx context.Context, content string, authorID int64, parentID int64) (*pb.Comment, error) {
	log.Println("Creating new comment in database")

	query := `
		INSERT INTO comments (content, author_id, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, content, author_id, parent_id, created_at, updated_at
	`

	now := time.Now()

	var id int64
	var returnedContent string
	var returnedAuthorID int64
	var returnedParentID int64
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, query, content, authorID, parentID, now, now).
		Scan(&id, &returnedContent, &returnedAuthorID, &returnedParentID, &createdAt, &updatedAt)
	if err != nil {
		log.Println("Error creating comment:", err)
		return nil, fmt.Errorf("error creating comment: %w", err)
	}

	return &pb.Comment{
		Id:        id,
		Content:   returnedContent,
		AuthorId:  returnedAuthorID,
		ParentId:  returnedParentID,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}, nil
}

func (r *CommentRepository) GetComments(ctx context.Context, parentID int64) ([]*pb.Comment, error) {
	log.Println("Fetching comments from database")

	query := `
		SELECT id, content, author_id, parent_id, created_at, updated_at
		FROM comments
		WHERE parent_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, parentID)
	if err != nil {
		log.Println("Error fetching comments:", err)
		return nil, fmt.Errorf("error fetching comments: %w", err)
	}
	defer rows.Close()

	var comments []*pb.Comment
	for rows.Next() {
		var id int64
		var authorID int64
		var content string
		var parentID sql.NullInt64
		var createdAt, updatedAt time.Time

		err := rows.Scan(&id, &content, &authorID, &parentID, &createdAt, &updatedAt)
		if err != nil {
			log.Println("Error scanning comment row:", err)
			return nil, fmt.Errorf("error scanning comment row: %w", err)
		}

		comments = append(comments, &pb.Comment{
			Id:        id,
			Content:   content,
			AuthorId:  authorID,
			ParentId:  parentID.Int64,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		})
	}

	return comments, nil
}

func (r *CommentRepository) UpdateComment(ctx context.Context, id int64, content string) (*pb.Comment, error) {
	log.Println("Updating comment in database")

	query := `
		UPDATE comments
		SET content = $2, updated_at = $3
		WHERE id = $1
		RETURNING id, content, author_id, parent_id, created_at, updated_at
	`

	now := time.Now()

	var returnedID int64
	var returnedAuthorID int64
	var returnedContent string
	var returnedParentID sql.NullInt64
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, query, id, content, now).
		Scan(&returnedID, &returnedContent, &returnedAuthorID, &returnedParentID, &createdAt, &updatedAt)
	if err != nil {
		log.Println("Error updating comment:", err)
		return nil, fmt.Errorf("error updating comment: %w", err)
	}

	return &pb.Comment{
		Id:        returnedID,
		Content:   returnedContent,
		AuthorId:  returnedAuthorID,
		ParentId:  returnedParentID.Int64,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}, nil
}

func (r *CommentRepository) DeleteComment(ctx context.Context, id int64) error {
	log.Println("Deleting comment from database")

	query := `DELETE FROM comments WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		log.Println("Error deleting comment:", err)
		return fmt.Errorf("error deleting comment: %w", err)
	}

	return nil
}

func (r *CommentRepository) GetCommentByID(ctx context.Context, id int64) (*pb.Comment, error) {
	log.Println("Fetching comment by ID from database")

	query := `
		SELECT id, content, author_id, parent_id, created_at, updated_at
		FROM comments
		WHERE id = $1
	`

	var returnedID int64
	var authorID int64
	var content string
	var parentID sql.NullInt64
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, query, id).
		Scan(&returnedID, &content, &authorID, &parentID, &createdAt, &updatedAt)
	if err != nil {
		log.Println("Error fetching comment:", err)
		return nil, fmt.Errorf("error fetching comment: %w", err)
	}

	return &pb.Comment{
		Id:        returnedID,
		Content:   content,
		AuthorId:  authorID,
		ParentId:  parentID.Int64,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}, nil
}