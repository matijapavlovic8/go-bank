package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	. "go-bank-v2/internal/types"
)

type Store interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccountBalance(*Account, float64) error
	GetAccounts(int) ([]*Account, error)
	GetAccountByNumber(int) (*Account, error)
	GetAllAccounts() ([]*Account, error)
	CreateUser(*User) error
	DeleteUser(int) error
	GetUserByID(int) (*User, error)
	Migrate() error
}

type PostgresqlStore struct {
	db *sql.DB
}

func NewPostgresStore(config DbConfig) (*PostgresqlStore, error) {

	db, err := sql.Open("postgres", connectionString(config))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresqlStore{
		db: db,
	}, nil
}

func connectionString(config DbConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
}

func (s *PostgresqlStore) Migrate() error {
	m := NewMigrator()

	err := m.Migrate(s.db)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresqlStore) close() {
	err := s.db.Close()
	if err != nil {
		return
	}
}
