package postgres

import (
	"database/sql"
	"embed"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"sync"
)

//go:embed sql
var migrations embed.FS

type Migrator struct {
	adminSecret      string
	createDefaultApp bool
	mx               sync.Mutex
}

func NewMigrator() *Migrator {
	m := &Migrator{}
	return m
}

func (m *Migrator) Migrate(db *sql.DB) error {
	if mErr := m.runMigrations(db); mErr != nil {
		return mErr
	}
	return nil
}

func (m *Migrator) runMigrations(db *sql.DB) error {

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	goose.SetBaseFS(migrations)

	if err := goose.Up(db, "sql"); err != nil {
		return errors.Wrap(err, "database sql failed")
	}

	return nil
}
