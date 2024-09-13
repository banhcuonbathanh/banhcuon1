package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"

	pb "english-ai-full/ecomm-grpc/proto/claude"
)

type ClaudeRepository struct {
	db *pgxpool.Pool
}

func NewClaudeRepository(db *pgxpool.Pool) *ClaudeRepository {
	return &ClaudeRepository{
		db: db,
	}
}

func (r *ClaudeRepository) CreateClaudeEvaluation(ctx context.Context, evaluationRequest *pb.EvaluationResponseFromToDataBase) (*pb.EvaluationResponseFromToDataBase, error) {
	log.Println("Creating new evaluation in database")

	query := `
		INSERT INTO claude_evaluations (
			student_response, passage, question, complex_sentences, 
			advanced_vocabulary, cohesive_devices, evaluation_from_claude
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRow(ctx, query,
		evaluationRequest.StudentResponse,
		evaluationRequest.Passage,
		evaluationRequest.Question,
		evaluationRequest.ComplexSentences,
		evaluationRequest.AdvancedVocabulary,
		evaluationRequest.CohesiveDevices,
		evaluationRequest.EvaluationFromClaude,
	).Scan(&id)

	if err != nil {
		log.Println("Error creating evaluation:", err)
		return nil, fmt.Errorf("error creating evaluation: %w", err)
	}

	return &pb.EvaluationResponseFromToDataBase{
		StudentResponse:      evaluationRequest.StudentResponse,
		Passage:              evaluationRequest.Passage,
		Question:             evaluationRequest.Question,
		ComplexSentences:     evaluationRequest.ComplexSentences,
		AdvancedVocabulary:   evaluationRequest.AdvancedVocabulary,
		CohesiveDevices:      evaluationRequest.CohesiveDevices,
		EvaluationFromClaude: 	evaluationRequest.EvaluationFromClaude,
	}, nil
}