package transport

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"user-management/internal"
)

// Endpoints holds all Go kit endpoints for the Order service.
type Endpoints struct {
	Create       endpoint.Endpoint
	GetByID      endpoint.Endpoint
	ChangeStatus endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the user service.
func MakeEndpoints(s user.Service) Endpoints {
	return Endpoints{
		Create:       makeCreateEndpoint(s),
		GetByID:      makeGetByIDEndpoint(s),
		ChangeStatus: makeChangeStatusEndpoint(s),
	}
}

func makeCreateEndpoint(s user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest) // type assertion
		id, err := s.Create(ctx, req.User)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetByIDEndpoint(s user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		userRes, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{User: userRes, Err: err}, nil
	}
}

func makeChangeStatusEndpoint(s user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeStatusRequest)
		err := s.ChangeStatus(ctx, req.ID, req.Status)
		return ChangeStatusResponse{Err: err}, nil
	}
}
