package Python_Api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	pb_python_ielts "english-ai-full/ecomm-grpc/proto/python_proto/claude"
)

type PythonIeltsHandlerController struct {
	ctx    context.Context
	client pb_python_ielts.IELTSServiceClient
}

func NewPythonIeltsHandler(client pb_python_ielts.IELTSServiceClient) *PythonIeltsHandlerController {
	return &PythonIeltsHandlerController{ctx: context.Background(), client: client}
}
func (h *PythonIeltsHandlerController) TestPythonGRPC(w http.ResponseWriter, r *http.Request) {
	// Log the start of the function
	log.Println("TestPythonGRPC called")

	// Create a new EvaluationResponseFromToDataBase message
	request := &pb_python_ielts.EvaluationResponseFromToDataBase{
		StudentResponse:    r.URL.Query().Get("student_response"),
		Passage:            r.URL.Query().Get("passage"),
		Question:           r.URL.Query().Get("question"),
		ComplexSentences:   r.URL.Query().Get("complex_sentences"),
		AdvancedVocabulary: r.URL.Query().Get("advanced_vocabulary"),
		CohesiveDevices:    r.URL.Query().Get("cohesive_devices"),
	}

	// Log the request details
	log.Printf("Request: %+v\n", request)

	resp, err := h.client.EvaluateIELTS(h.ctx, request)
	if err != nil {
		// Log the error
		log.Printf("Error calling Python gRPC service: %v\n", err)
		http.Error(w, "Error calling Python gRPC service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Log the response details
	log.Printf("Response: %+v\n", resp)

	// Create a response map with all fields from EvaluationResponseFromToDataBase
	response := map[string]string{
		"student_response":      resp.StudentResponse,
		"passage":               resp.Passage,
		"question":              resp.Question,
		"complex_sentences":     resp.ComplexSentences,
		"advanced_vocabulary":   resp.AdvancedVocabulary,
		"cohesive_devices":      resp.CohesiveDevices,
		"evaluation_from_claude": resp.EvaluationFromClaude,
	}

	// Log the final response map
	log.Printf("Final response map: %+v\n", response)

	json.NewEncoder(w).Encode(response)
}
