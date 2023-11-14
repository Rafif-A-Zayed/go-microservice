package middleware

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
	"user-management/internal"
	logger "user-management/internal/util"
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
		logger.Info(mw.logger, "method", "Create", "UserID", user.CustomerID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Create(ctx, user)
}

func (mw loggingMiddleware) GetByID(ctx context.Context, id string) (user user.User, err error) {
	defer func(begin time.Time) {
		logger.Info(mw.logger, "method", "GetByID", "UserID", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetByID(ctx, id)
}

func (mw loggingMiddleware) ChangeStatus(ctx context.Context, id string, status string) (err error) {
	defer func(begin time.Time) {
		logger.Info(mw.logger, "method", "ChangeStatus", "UserID", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.ChangeStatus(ctx, id, status)
}
