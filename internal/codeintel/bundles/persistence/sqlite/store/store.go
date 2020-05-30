package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/keegancsmith/sqlf"
	"github.com/pkg/errors"
	"github.com/sourcegraph/sourcegraph/internal/db/dbutil"
	"github.com/sourcegraph/sourcegraph/internal/sqliteutil"
)

type Store struct {
	db dbutil.DB
}

var _ sqliteutil.Execable = &Store{}

func New(db dbutil.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Query(ctx context.Context, query *sqlf.Query) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query.Query(sqlf.SimpleBindVar), query.Args()...)
}

//
//

func (s *Store) ExecAll(ctx context.Context, queries ...*sqlf.Query) error {
	for _, query := range queries {
		// NOTE: NEED TO HAVE ACTUAL EXEC, CANNOT CREATE TABLE WITH QUERY
		fmt.Printf("DO QUERY %s\n", query)
		// TODO - need to dispose of rows?
		if _, err := s.Query(ctx, query); err != nil {
			fmt.Printf("NOP#: %v\n", err)
			return err
		}
	}

	return nil
}

// TODO - rework this interface
func (s *Store) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	// TODO - need to dispose of rows?
	_, err := s.Query(ctx, sqlf.Sprintf(query, args...))
	return nil, err
}

//
//

var ErrNotTransactable = errors.New("db: not transactable")

func (s *Store) Transact(ctx context.Context) (*Store, error) {
	if _, ok := s.db.(dbutil.Tx); ok {
		// Already in a Tx
		return s, nil
	}

	tb, ok := s.db.(dbutil.TxBeginner)
	if !ok {
		// Not a Tx nor a TxBeginner
		return nil, ErrNotTransactable
	}

	tx, err := tb.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "db: BeginTx")
	}

	return &Store{db: tx}, nil
}

func (s *Store) Done(err error) error {
	if tx, ok := s.db.(dbutil.Tx); ok {
		if err != nil {
			if rollErr := tx.Rollback(); rollErr != nil {
				err = multierror.Append(err, rollErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = multierror.Append(err, commitErr)
			}
		}
	}

	return err
}
