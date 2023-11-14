package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	user "user-management/internal"
	usertransport "user-management/internal/transport"
	logger "user-management/internal/util"
)

var (
	ErrBadRouting = errors.New("bad routing")
)

// NewService wires Go kit endpoints to the HTTP transport.
func NewHandler(
	svcEndpoints usertransport.Endpoints, logger log.Logger,
) http.Handler {
	// set-up router and initialize http endpoints

	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeErrorResponse),
		kithttp.ServerFinalizer(newServerFinalizer(logger)),
	}

	// HTTP Post - /orders
	r.Methods("POST").Path("/users").Handler(kithttp.NewServer(
		svcEndpoints.Create,
		decodeCreateRequest,
		encodeResponse,
		options...,
	))

	// HTTP Post - /orders/{id}
	r.Methods("GET").Path("/users/{id}").Handler(kithttp.NewServer(
		svcEndpoints.GetByID,
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	))

	// HTTP Post - /orders/status
	r.Methods("POST").Path("/users/status").Handler(kithttp.NewServer(
		svcEndpoints.ChangeStatus,
		decodeChangeStausRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req usertransport.CreateRequest
	if e := json.NewDecoder(r.Body).Decode(&req.User); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return usertransport.GetByIDRequest{ID: id}, nil
}

func decodeChangeStausRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req usertransport.ChangeStatusRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeErrorResponse(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case user.ErrUserNotFound:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func newServerFinalizer(lgr log.Logger) kithttp.ServerFinalizerFunc {
	return func(ctx context.Context, code int, r *http.Request) {
		logger.Info(lgr, "status", code, "path", r.RequestURI, "method", r.Method)

	}
}
