package serviceielts

import (
	// "context"
	// "fmt"

	pb "english-ai-full/ecomm-grpc/proto/claude"
	repository "english-ai-full/ecomm-grpc/repository/ielts_repository"
)

type IELTSService struct {
	claudeRepo *repository.ClaudeRepository

	pb.UnimplementedIELTSServiceServer
}

func NewIELTSService(claudeRepo *repository.ClaudeRepository) *IELTSService {
	return &IELTSService{
		claudeRepo: claudeRepo,
	
	}
}

// func (s *IELTSService) EvaluateIELTS(ctx context.Context, req *pb.EvaluationResponseFromToDataBase) (*pb.EvaluationResponseFromToDataBase, error) {
// 	// evaluation, err := s.completeIELTSEvaluation(ctx, req)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// req.EvaluationFromClaude = evaluation
// 	// return s.claudeRepo.CreateClaudeEvaluation(ctx, req)
// }

// func (s *IELTSService) completeIELTSEvaluation(ctx context.Context, req *pb.EvaluationResponseFromToDataBase) (string, error) {
// 	// prompt := fmt.Sprintf(
// 	// 	"Please evaluate the following IELTS response:\n\nQuestion: %s\n\nAnswer: %s\n\nProvide a detailed evaluation including strengths, weaknesses, and an overall band score.",
// 	// 	req.Question,
// 	// 	req.Answer,
// 	// )

// 	// evaluation, err := s.claudeService.CallAnthropic(ctx, prompt, 1000) // Adjust maxTokens as needed
// 	if err != nil {
// 		return "", fmt.Errorf("error calling Anthropic API: %w", err)
// 	}

// 	return evaluation, nil
// }