package Python_Api

import (
	"context"
	"encoding/json"

	"net/http"

	pb_python "english-ai-full/ecomm-grpc/proto/python_proto"
)

type PythonHandler struct {
	ctx    context.Context
	client pb_python.GreeterClient
}

func NewPythonHandler(client pb_python.GreeterClient) *PythonHandler {
	return &PythonHandler{ctx: context.Background(), client: client}
}

func (h *PythonHandler) TestPythonGRPC(w http.ResponseWriter, r *http.Request) {


	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	resp, err := h.client.SayHello(h.ctx, &pb_python.HelloRequest{Name: name})
	if err != nil {
		http.Error(w, "Error calling Python gRPC service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": resp.Message})
}