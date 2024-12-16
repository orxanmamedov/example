package mysql

import (
	"database/sql"
	"fmt"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io/fs"
	"time"
)

type NodeConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	MaxOpen  uint
	TimeOut  time.Duration
}

type Config struct {
	DBconfig    NodeConfig
	Metrics     bool
	Migrations  bool
	MigrationFS fs.FS
}

var ErrNoRows = sql.ErrNoRows

type MySQL interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*Stmt, error)
	Master() *Storage
	Slave() *Storage
}

type Storage struct {
	db *sqlx.DB

	master *sqlx.DB
	slave  *sqlx.DB
}

func New(cfg Config) (*Storage, error) {
	masterDB, err := sqlx.Connect("mysql", dsn(cfg.DBconfig))
	if err != nil {
		return nil, fmt.Errorf("master DB: %v", err)
	}
	masterDB.DB.SetMaxOpenConns(int(cfg.DBconfig.MaxOpen))

	s := &Storage{
		db:     masterDB,
		master: masterDB,
	}

	return s, err
}

func (s *Storage) Master() *Storage {
	return s
}

func (s *Storage) Slave() *Storage {
	return &Storage{
		db: s.slave,
	}
}

func (s *Storage) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	//started := time.Now()
	err := s.db.GetContext(ctx, dest, query, args...)
	//s.m.witre(ctx, started, query, err)

	return err
}

func (s *Storage) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	//started := time.Now()
	err := s.db.SelectContext(ctx, dest, query, args...)
	//s.m.write(ctx, started, query, err)

	return err
}

func (s *Storage) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	//started := time.Now()
	res, err := s.db.ExecContext(ctx, query, args...)
	//s.m.write(ctx, started, query, err)

	return res, err
}

func (s *Storage) Transaction(ctx context.Context, t func(tx *sqlx.Tx) error) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	err = t(tx)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return errors.Wrapf(err, "rollback error: %v", txErr)
		}
		return err
	}
	return tx.Commit()
}

func (s *Storage) TransactionLevel(ctx context.Context, level sql.IsolationLevel, t func(tx *sqlx.Tx) error) error {
	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{Isolation: level})
	if err != nil {
		return err
	}
	err = t(tx)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return errors.Wrapf(err, "rollback error: %v", txErr)
		}
		return err
	}
	return tx.Commit()
}

func (s *Storage) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	stmt, err := s.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{stmt: stmt}, nil
}

type Stmt struct {
	query string
	stmt  *sqlx.Stmt
}

func (s *Stmt) Close() error {
	return s.stmt.Close()
}

func (s *Stmt) GetContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	//started := time.Now()
	err := s.stmt.GetContext(ctx, dest, args...)
	//s.m.write(ctx, started, s.query, err)

	return err
}

func (s *Stmt) SelectContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	//started := time.Now()
	err := s.stmt.SelectContext(ctx, dest, args...)
	//s.m.write(ctx, started, s.query, err)

	return err
}

func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	//started := time.Now()
	res, err := s.stmt.ExecContext(ctx, args...)
	//s.m.write(ctx, started, s.query, err)

	return res, err
}

func (s *Storage) Close() error {
	if err := s.master.Close(); err != nil {
		return errors.Wrap(err, "db master")
	}

	if err := s.slave.Close(); err != nil {
		return errors.Wrap(err, "db slave")
	}
	return nil
}

const HealthCheckTimeout = 5 * time.Second

func (s *Storage) HealthCheck(timeout time.Duration) func() error {
	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if s == nil || s.master == nil || s.slave == nil {
			return fmt.Errorf("database is nil")
		}
		if err := s.master.PingContext(ctx); err != nil {
			return errors.Wrap(err, "db master")
		}
		if err := s.slave.PingContext(ctx); err != nil {
			return errors.Wrap(err, "db slave")
		}
		return nil
	}
}

func dsn(cfg NodeConfig) string {
	c := mysqlDriver.NewConfig()
	c.User = cfg.User
	c.Passwd = cfg.Password
	c.Net = "tcp"
	c.Addr = fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	c.DBName = cfg.Database
	c.ParseTime = true
	c.Timeout = cfg.TimeOut

	return c.FormatDSN()
}
