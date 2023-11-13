package implementation

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
	"time"

	usersvc "user-management/internal"
)

// service implements the Order Service
type Service struct {
	repository usersvc.Repository
	logger     log.Logger
}

// NewService creates and returns a new Order service instance
func NewService(rep usersvc.Repository, logger log.Logger) (*Service, error) {
	// return  repository
	return &Service{
		repository: rep,
		logger:     logger,
	}, nil
}

// Create makes an order
func (s *Service) Create(ctx context.Context, user usersvc.User) (string, error) {
	logger := log.With(s.logger, "method", "Create")
	uuid1, _ := uuid.NewV4()
	id := uuid1.String()
	user.ID = id
	user.Status = "Pending"
	user.CustomerID = "U12336286"
	user.Address = "user address"
	user.CreatedOn = time.Now().Unix()

	if err := s.repository.CreateUser(ctx, user); err != nil {
		level.Error(logger).Log("err", err)
		return "", usersvc.ErrCmdRepository
	}
	return id, nil
}

// GetByID returns an order given by id
func (s *Service) GetByID(ctx context.Context, id string) (usersvc.User, error) {
	logger := log.With(s.logger, "method", "GetByID")
	level.Info(logger).Log(id)
	user := usersvc.User{
		ID:         "",
		CustomerID: "U12336286",
		Address:    "user address",
		Status:     "Pending",
	}
	/*order, err := s.repository.GetOrderByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		if err == sql.ErrNoRows {
			return order, ordersvc.ErrOrderNotFound
		}
		return order, ordersvc.ErrQueryRepository
	}*/
	return user, nil
}

// ChangeStatus changes the status of an order
func (s *Service) ChangeStatus(ctx context.Context, id string, status string) error {
	logger := log.With(s.logger, "method", "ChangeStatus")
	level.Info(logger).Log(id, status)
	/*	if err := s.repository.ChangeOrderStatus(ctx, id, status); err != nil {
		level.Error(logger).Log("err", err)
		return ordersvc.ErrCmdRepository
	}*/
	return nil
}
