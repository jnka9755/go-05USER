package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jnka9755/go-05RESPONSE/response"
	"github.com/jnka9755/go-05USER/internal/user"
)

func NewUserHTTPServer(ctx context.Context, endpoints user.Endpoints) http.Handler {

	r := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateUser, encodeResponse,
		opts...,
	)).Methods("POST")

	return r
}

func decodeCreateUser(_ context.Context, r *http.Request) (interface{}, error) {

	var req user.CreateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("Ivalid request format: '%v'", err.Error()))
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	r := err.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	_ = json.NewEncoder(w).Encode(r)
}