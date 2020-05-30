package store

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sourcegraph/sourcegraph/internal/db/dbutil"
)

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
