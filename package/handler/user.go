package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

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

	r.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetUser,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllUser,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateUser,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	r.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteUser,
		encodeResponse,
		opts...,
	)).Methods("DELETE")

	return r
}

func decodeCreateUser(_ context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	var req user.CreateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("Ivalid request format: '%v'", err.Error()))
	}

	return req, nil
}

func decodeGetUser(_ context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	p := mux.Vars(r)
	req := user.GetReq{
		ID: p["id"],
	}

	return req, nil
}

func decodeGetAllUser(_ context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := user.GetAllReq{
		FirstName: v.Get("first_name"),
		LastName:  v.Get("last_name"),
		Limit:     limit,
		Page:      page,
	}

	return req, nil
}

func decodeUpdateUser(_ context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	var req user.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("Invalid request format: '%v'", err.Error()))
	}

	p := mux.Vars(r)
	req.ID = p["id"]

	return req, nil
}

func decodeDeleteUser(_ context.Context, r *http.Request) (interface{}, error) {

	if err := authorization(r.Header.Get("Authorization")); err != nil {
		return nil, response.Forbidden(err.Error())
	}

	p := mux.Vars(r)
	req := user.DeleteReq{
		ID: p["id"],
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

func authorization(token string) error {

	if token != os.Getenv("TOKEN") {
		return errors.New("Invalid token")
	}

	return nil
}
