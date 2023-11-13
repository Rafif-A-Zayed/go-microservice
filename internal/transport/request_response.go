package transport

import "user-management/internal"

// CreateRequest holds the request parameters for the Create method.
type CreateRequest struct {
	User user.User
}

// CreateResponse holds the response values for the Create method.
type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

// GetByIDRequest holds the request parameters for the GetByID method.
type GetByIDRequest struct {
	ID string
}

// GetByIDResponse holds the response values for the GetByID method.
type GetByIDResponse struct {
	User user.User `json:"user"`
	Err  error     `json:"error,omitempty"`
}

// ChangeStatusRequest holds the request parameters for the ChangeStatus method.
type ChangeStatusRequest struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// ChangeStatusResponse holds the response values for the ChangeStatus method.
type ChangeStatusResponse struct {
	Err error `json:"error,omitempty"`
}
