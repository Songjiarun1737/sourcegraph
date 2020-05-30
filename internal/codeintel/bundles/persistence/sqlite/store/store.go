package store

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/keegancsmith/sqlf"
	"github.com/sourcegraph/sourcegraph/internal/db/dbutil"
	"github.com/sourcegraph/sourcegraph/internal/sqliteutil"
)

type ExecableDB interface {
	dbutil.DB
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type Store struct {
	db ExecableDB
}

var _ sqliteutil.Execable = &Store{}

func Open(filename string) (*Store, func() error, error) {
	db, err := sqlx.Open("sqlite3_with_pcre", filename)
	if err != nil {
		return nil, nil, err
	}

	return &Store{db: db}, db.Close, nil
}

func (s *Store) Query(ctx context.Context, query *sqlf.Query) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query.Query(sqlf.SimpleBindVar), query.Args()...)
}

func (s *Store) Exec(ctx context.Context, query *sqlf.Query) (sql.Result, error) {
	return s.db.ExecContext(ctx, query.Query(sqlf.SimpleBindVar), query.Args()...)
}

func (s *Store) ExecAll(ctx context.Context, queries ...*sqlf.Query) error {
	for _, query := range queries {
		if _, err := s.Exec(ctx, query); err != nil {
			return err
		}
	}

	return nil
}

// TODO - rework this interface
func (s *Store) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}
