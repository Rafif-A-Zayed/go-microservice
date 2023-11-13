package middleware

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
	"user-management/internal"
)

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next user.Service) user.Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   user.Service
	logger log.Logger
}

func (mw loggingMiddleware) Create(ctx context.Context, user user.User) (id string, err error) {
	defer func(begin time.Time) {
		err := mw.logger.Log("method", "Create", "CustomerID", user.CustomerID, "took", time.Since(begin), "err", err)
		if err != nil {
			return
		}
	}(time.Now())
	return mw.next.Create(ctx, user)
}

func (mw loggingMiddleware) GetByID(ctx context.Context, id string) (user user.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetByID", "OrderID", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetByID(ctx, id)
}

func (mw loggingMiddleware) ChangeStatus(ctx context.Context, id string, status string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "ChangeStatus", "OrderID", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.ChangeStatus(ctx, id, status)
}
