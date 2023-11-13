package user

import "context"

type User struct {
	ID         string `json:"id,omitempty"`
	CustomerID string `json:"customer_id"`
	Status     string `json:"status"`
	CreatedOn  int64  `json:"created_on,omitempty"`
	Address    string `json:"address"`
}

// Repository describes the persistence on order model
type Repository interface {
	CreateUser(ctx context.Context, user User) error
	//GetUserByID(ctx context.Context, id string) (User, error)
	//	ChangeUserStatus(ctx context.Context, id string, status string) error
}
