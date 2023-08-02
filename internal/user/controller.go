package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jnka9755/go-05META/meta"
	"github.com/jnka9755/go-05RESPONSE/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)
	Endpoints  struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	GetReq struct {
		ID string
	}

	GetAllReq struct {
		FirstName string
		LastName  string
		Limit     int
		Page      int
	}

	UpdateReq struct {
		ID        string
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	DeleteReq struct {
		ID string
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(b Business, config Config) Endpoints {

	return Endpoints{
		Create: makeCreateEndpoint(b),
		Get:    makeGetEndpoint(b),
		GetAll: makeGetAllEndpoint(b, config),
		Update: makeUpdateEndpoint(b),
		Delete: makeDeleteEndpoint(b),
	}
}

func makeCreateEndpoint(b Business) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(CreateReq)

		if req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}

		if req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		responseUser, err := b.Create(ctx, &req)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("Success create user", responseUser, nil), nil
	}
}

func makeGetEndpoint(b Business) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetReq)

		user, err := b.Get(ctx, req.ID)

		if err != nil {
			return nil, response.NotFound(err.Error())
		}

		return response.OK("Success get user", user, nil), nil
	}
}

func makeGetAllEndpoint(b Business, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetAllReq)

		filters := Filters{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		count, err := b.Count(ctx, filters)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		users, err := b.GetAll(ctx, filters, meta.Offset(), meta.Limit())

		if err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("Success get users", users, meta), nil
	}
}

func makeUpdateEndpoint(b Business) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(UpdateReq)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}

		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		err := b.Update(ctx, &req)

		if err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK(fmt.Sprintf("Success update user with ID -> '%s'", req.ID), nil, nil), nil
	}
}

func makeDeleteEndpoint(b Business) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(DeleteReq)

		err := b.Delete(ctx, req.ID)

		if err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK(fmt.Sprintf("Success delete user with ID -> '%s'", req.ID), nil, nil), nil
	}
}
