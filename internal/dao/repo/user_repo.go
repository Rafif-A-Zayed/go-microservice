package repo

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"user-management/internal"
)

type UserSQLRepository struct {
	db     *gorm.DB
	logger log.Logger
}

// NewUserQLRepository creates new instance of UserQLRepository.
/*func NewUserSQLRepository(db *gorm.DB, logger log.Logger) (*user.Repository, error) {
	return &UserSQLRepository{db: db, logger: log.With(logger, "rep", "userDB")}, nil
}*/

// New returns a concrete repository backed by CockroachDB
func NewUserSQLRepository(db *gorm.DB, logger log.Logger) (*UserSQLRepository, error) {
	// return  repository
	return &UserSQLRepository{
		db:     db,
		logger: log.With(logger, "rep", "userDB"),
	}, nil
}

// CreateUser inserts a new order and its order items into db
func (repo *UserSQLRepository) CreateUser(ctx context.Context, user user.User) error {

	res := repo.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(UserModel{
			ID:         user.ID,
			CustomerID: user.CustomerID,
			Status:     user.Status,
			Address:    user.Address,
		})

	if res.Error != nil {
		return fmt.Errorf("failed to create user in DB: %w", res.Error)
	}

	return nil

}
