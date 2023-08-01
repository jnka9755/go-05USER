package user

import (
	"context"

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

	UpdateReq struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
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
		// Get:    makeGetEndpoint(b),
		// GetAll: makeGetAllEndpoint(b, config),
		// Update: makeUpdateEndpoint(b),
		// Delete: makeDeleteEndpoint(b),
	}
}

func makeCreateEndpoint(b Business) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(CreateReq)

		if req.FirstName == "" {
			return nil, response.BadRequest("first name is required")
		}

		if req.LastName == "" {
			return nil, response.BadRequest("last name is required")
		}

		responseUser, err := b.Create(ctx, &req)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("Success", responseUser, nil), nil
	}
}

// func makeGetEndpoint(b Business) Controller {
// 	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
// 		path := mux.Vars(r)
// 		id := path["id"]

// 		user, err := b.Get(id)

// 		if err != nil {

// 			w.WriteHeader(400)
// 			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
// 			return
// 		}

// 		json.NewEncoder(w).Encode(&Response{Status: 200, Data: user})
// 	}
// }

// func makeGetAllEndpoint(b Business, config Config) Controller {
// 	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

// 		value := r.URL.Query()

// 		filters := Filters{
// 			FirstName: value.Get("first_name"),
// 			LastName:  value.Get("last_name"),
// 		}

// 		limit, _ := strconv.Atoi(value.Get("limit"))
// 		page, _ := strconv.Atoi(value.Get("page"))

// 		count, err := b.Count(filters)

// 		if err != nil {
// 			w.WriteHeader(500)
// 			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
// 			return
// 		}

// 		meta, err := meta.New(page, limit, count, config.LimPageDef)

// 		if err != nil {
// 			w.WriteHeader(500)
// 			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
// 			return
// 		}

// 		users, err := b.GetAll(filters, meta.Offset(), meta.Limit())

// 		if err != nil {
// 			w.WriteHeader(400)
// 			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
// 			return
// 		}

// 		json.NewEncoder(w).Encode(&Response{Status: 200, Data: users, Meta: meta})
// 	}
// }

// func makeUpdateEndpoint(b Business) Controller {
// 	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

// 		var request UpdateReq

// 		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 			w.WriteHeader(400)
// 			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "Invalid request format"})
// 			return
// 		}

// 		if request.FirstName != nil && *request.FirstName == "" {
// 			w.WriteHeader(400)
// 			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "first_name is required"})
// 			return
// 		}

// 		if request.LastName != nil && *request.LastName == "" {
// 			w.WriteHeader(400)
// 			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "last_name is required"})
// 			return
// 		}

// 		path := mux.Vars(r)
// 		id := path["id"]

// 		if err := b.Update(id, &request); err != nil {
// 			w.WriteHeader(404)
// 			json.NewEncoder(w).Encode(&Response{Status: 404, Err: "User doesn't exist"})
// 			return
// 		}

// 		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "Successful update"})
// 	}
// }

// func makeDeleteEndpoint(b Business) Controller {
// 	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
// 		path := mux.Vars(r)
// 		id := path["id"]

// 		if err := b.Delete(id); err != nil {
// 			w.WriteHeader(404)
// 			json.NewEncoder(w).Encode(&Response{Status: 404, Err: "User doesn't exist"})
// 		}

// 		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "Successful delete"})
// 	}
// }
