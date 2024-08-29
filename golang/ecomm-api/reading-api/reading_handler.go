package reading_api

import (
	"context"
	"encoding/json"
	"log"

	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"google.golang.org/protobuf/types/known/emptypb"

	mapping_user "english-ai-full/ecomm-api/mapping"
	"english-ai-full/ecomm-api/types"
	pb "english-ai-full/ecomm-grpc/proto/reading"
	"english-ai-full/token"
)

type ReadingHandlerController struct {
	ctx    context.Context
	client pb.EcommReadingClient
	TokenMaker *token.JWTMaker
}

func NewReadingHandler(client pb.EcommReadingClient,secretKey string) *ReadingHandlerController {
	return &ReadingHandlerController{
		ctx:    context.Background(),
		client: client,
		TokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (h *ReadingHandlerController) CreateReading(w http.ResponseWriter, r *http.Request) {
	var reading types.ReadingReqModel
	if err := json.NewDecoder(r.Body).Decode(&reading); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}
	log.Println("handler CreateReading before")
	createdReading, err := h.client.CreateReading(h.ctx, mapping_user.ToPBReadingReq(reading))
	if err != nil {
		log.Println("handler CreateReading err ", err)
		http.Error(w, "error creating reading in handler", http.StatusInternalServerError)
		return
	}
	log.Println("handler CreateReading after")
	res := mapping_user.ToReadingResFromPbReadingRes(createdReading)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *ReadingHandlerController) FindByID(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    i, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        http.Error(w, "error parsing ID", http.StatusBadRequest)
        return
    }

    // Convert i to a string
    idString := strconv.FormatInt(i, 10)
    reading, err := h.client.FindByID(h.ctx, &pb.ReadingRes{Id: idString})
    if err != nil {
        http.Error(w, "error getting reading", http.StatusInternalServerError)
        return
    }

    res := mapping_user.ToReadingResFromPbReadingRes(reading)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(res)
}


func (h *ReadingHandlerController) ListReadings(w http.ResponseWriter, r *http.Request) {
	lrr, err := h.client.FindAllReading(h.ctx, &emptypb.Empty{})
	if err != nil {
		http.Error(w, "failed to fetch readings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var res []types.ReadingResModel
	for _, r := range lrr.GetReadings() {
		readingRes := mapping_user.ToReadingResFromPbReadingRes(r)
		res = append(res, readingRes)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ReadingHandlerController) UpdateReading(w http.ResponseWriter, r *http.Request) {
	var reading types.ReadingReqModel
	if err := json.NewDecoder(r.Body).Decode(&reading); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	updatedReading, err := h.client.UpdateReading(h.ctx, mapping_user.ToPBReadingReq(reading))
	if err != nil {
		http.Error(w, "error updating reading", http.StatusInternalServerError)
		return
	}

	res := mapping_user.ToReadingResFromPbReadingRes(updatedReading)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *ReadingHandlerController) DeleteReading(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    i, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        http.Error(w, "error parsing ID", http.StatusBadRequest)
        return
    }

    // Convert i to a string
    idString := strconv.FormatInt(i, 10)

    _, err = h.client.DeleteReading(h.ctx, &pb.ReadingRes{Id: idString})
    if err != nil {
        http.Error(w, "error deleting reading", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}


func (h *ReadingHandlerController) FindReadingByPage(w http.ResponseWriter, r *http.Request) {
	pageNumber, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("size"))

	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	result, err := h.client.FindReadingByPage(h.ctx, &pb.PageRequestReading{
		PageNumber: int32(pageNumber),
		PageSize:   int32(pageSize),
	})
	if err != nil {
		http.Error(w, "failed to fetch readings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := struct {
		Readings   []types.ReadingResModel `json:"readings"`
		TotalCount int32                   `json:"totalCount"`
	}{
		Readings:   mapping_user.ToReadingResListFromPbReadingResList(result.Readings),
		TotalCount: result.TotalCount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}