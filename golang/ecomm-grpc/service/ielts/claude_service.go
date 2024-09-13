package ielts

import (
	"context"

	repository "english-ai-full/ecomm-grpc/repository/ielts_repository"

	pb "english-ai-full/ecomm-grpc/proto/claude"
)

type IELTSService struct {
	claudeRepo *repository.ClaudeRepository
	pb.UnimplementedIELTSServiceServer
}

func NewIELTSService(claudeRepo *repository.ClaudeRepository) *IELTSService {
	return &IELTSService{claudeRepo: claudeRepo}
}

func (s *IELTSService) EvaluateIELTS(ctx context.Context, req *pb.EvaluationResponseFromToDataBase) (*pb.EvaluationResponseFromToDataBase, error) {
	evaluation, err := s.completeIELTSEvaluation(req)
	if err != nil {
		return nil, err
	}

	req.EvaluationFromClaude = evaluation
	return s.claudeRepo.CreateClaudeEvaluation(ctx, req)
}

// func (s *IELTSService) completeIELTSEvaluation(req *pb.EvaluationResponseFromToDataBase) (string, error) {
// 	grammarEval, err := s.evaluateGrammar(req.StudentResponse, req.ComplexSentences)
// 	if err != nil {
// 		return "", err
// 	}

// 	vocabEval, err := s.evaluateVocabulary(req.StudentResponse, req.AdvancedVocabulary)
// 	if err != nil {
// 		return "", err
// 	}

// 	coherenceEval, err := s.evaluateCoherence(req.StudentResponse, req.CohesiveDevices)
// 	if err != nil {
// 		return "", err
// 	}

// 	contentEval, err := s.evaluateContent(req.StudentResponse, req.Passage, req.Question)
// 	if err != nil {
// 		return "", err
// 	}

// 	return s.synthesizeEvaluation(grammarEval, vocabEval, coherenceEval, contentEval)
// }

// func (s *IELTSService) evaluateGrammar(response, complexSentences string) (string, error) {
// 	prompt := fmt.Sprintf(`
// 	As an IELTS examiner, evaluate the grammar in this response. Focus on:
// 	1. Sentence structure variety
// 	2. Use of complex grammatical constructions
// 	3. Accuracy of grammar usage

// 	Student's response:
// 	%s

// 	Complex structures to look for:
// 	%s

// 	Provide:
// 	1. A brief assessment of grammar usage
// 	2. Examples of well-used complex structures
// 	3. Suggestions for improvement
// 	`, response, complexSentences)

// 	return anthropic.CallAnthropic(prompt, 1000)
// }