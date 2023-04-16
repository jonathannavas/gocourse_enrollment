package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/jonathannavas/go_lib_response/response"
	"github.com/jonathannavas/gocourse_enrollment/internal/enrollment"
)

func NewUserHTTPServer(ctx context.Context, endpoints enrollment.Endpoints) http.Handler {
	r := mux.NewRouter()
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}
	r.Handle("/enrollments", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateEnrollment, encodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/enrollments", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllEnrollment, encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/enrollments/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateEnrollment, encodeResponse,
		opts...,
	)).Methods("PATCH")

	return r
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func decodeCreateEnrollment(_ context.Context, r *http.Request) (interface{}, error) {
	var req enrollment.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("Invalid request format: '%v'", err.Error()))
	}
	return req, nil
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}

func decodeGetAllEnrollment(_ context.Context, r *http.Request) (interface{}, error) {
	v := r.URL.Query()
	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("limit"))

	req := enrollment.GetAllRequest{
		CourseID: v.Get("course_id"),
		UserID:   v.Get("user_id"),
		Limit:    limit,
		Page:     page,
	}

	return req, nil
}

func decodeUpdateEnrollment(_ context.Context, r *http.Request) (interface{}, error) {
	var req enrollment.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("Invalid request format: '%v'", err.Error()))
	}
	path := mux.Vars(r)
	req.ID = path["id"]
	return req, nil
}
