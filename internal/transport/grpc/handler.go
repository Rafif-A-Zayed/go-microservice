package grpc

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	user "user-management/internal"
	usertransport "user-management/internal/transport"
)

type GrpcServer struct {
	createUser kitgrpc.Handler
	getUser    kitgrpc.Handler
	logger     log.Logger
}

func (s GrpcServer) mustEmbedUnimplementedUserServer() {
	//TODO implement me
	panic("implement me")
}

// NewService wires Go kit endpoints to the GRPC transport.
func NewHandler(
	svcEndpoints usertransport.Endpoints, logger log.Logger,
) UserServer {

	options := []kitgrpc.ServerOption{
		kitgrpc.ServerFinalizer(newServerFinalizer(logger)),
	}

	return &GrpcServer{
		createUser: kitgrpc.NewServer(
			svcEndpoints.Create, decodeCreateRequest, encodeCreateResponse, options...,
		),
		getUser: kitgrpc.NewServer(
			svcEndpoints.GetByID, decodeGetRequest, encodeGetResponse, options...,
		),
		logger: logger,
	}

}

func (s GrpcServer) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	_, rep, err := s.createUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*CreateUserResponse), nil
}

func decodeCreateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*CreateUserRequest)
	return usertransport.CreateRequest{
		User: user.User{
			CustomerID: req.CustomerId,
			Address:    req.Address,
			Status:     req.Status,
		},
	}, nil
}

func encodeCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(usertransport.CreateResponse)
	err := getError(res.Err)
	if err == nil {
		return &CreateUserResponse{}, nil
	}
	return nil, err
}

func (s GrpcServer) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	_, rep, err := s.getUser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*GetUserResponse), nil
}

func decodeGetRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*GetUserRequest)
	return usertransport.GetByIDRequest{
		ID: req.UserId,
	}, nil
}
func encodeGetResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(usertransport.GetByIDResponse)
	err := getError(res.Err)
	if err == nil {
		return &GetUserResponse{}, nil
	}

	return &GetUserResponse{
		UserId:     res.User.ID,
		Address:    res.User.Address,
		Status:     res.User.Status,
		CustomerId: res.User.CustomerID,
	}, err
}

func getError(err error) error {
	switch err {
	case nil:
		return nil
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

func newServerFinalizer(logger log.Logger) kitgrpc.ServerFinalizerFunc {
	return func(ctx context.Context, err error) {
		level.Info(logger).Log("status", err)
	}
}
