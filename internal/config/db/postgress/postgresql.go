package postgress

import (
	"database/sql"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"user-management/internal/dao/repo"
)

// InitDB initialises PG database connection.
func InitDB(cfg Config, lgr log.Logger) (*gorm.DB, error) {
	var (
		master *sql.DB
		err    error
	)

	err = level.Info(lgr).Log("Connecting to DB... User: %s, host: %s:%d, database:%s", map[string]interface{}{
		"user":     cfg.User,
		"host":     cfg.Master,
		"database": cfg.Name,
		"port":     cfg.Port,
	})
	if err != nil {
		return nil, err
	}

	dns := cfg.MasterDSN()
	master, err = openDB(dns, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not open master db: %w", err)
	}

	gormDB, err := openGORM(master, cfg)

	if err != nil {
		err1 := master.Close()
		if err1 != nil {
			return nil, fmt.Errorf("could not close DB connection string: %w", err1)
		}

		return nil, fmt.Errorf("could open grom connection string: %w", err)
	}

	err = gormDB.AutoMigrate(&repo.UserModel{})
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}

func openDB(dsn string, cfg Config) (*sql.DB, error) {
	pgxcfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("could not parse DB connection string: %w", err)
	}

	db := stdlib.OpenDB(*pgxcfg)
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.MaxConnectionLifetime)

	/*defer func(db *sql.DB) error {
		err := db.Close()
		if err != nil {
			return fmt.Errorf("could not parse DB connection string: %w", err)
		}
		return nil
	}(db)*/
	return db, nil
}

func openGORM(master *sql.DB, cfg Config) (*gorm.DB, error) {
	dialector := postgres.New(postgres.Config{Conn: master})
	logLevel := gormLogger.LogLevel(cfg.LogLevel)
	config := &gorm.Config{
		Logger:                 gormLogger.Default.LogMode(logLevel),
		SkipDefaultTransaction: true,
	}

	gormDB, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, fmt.Errorf("could not open gorm DB: %w", err)
	}
	return gormDB, nil
}
