package v4

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization"
)

func Migrate(ctx context.Context, db *sqlx.DB, serializer serialization.Serializer) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE meta DROP COLUMN lsifVersion`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE meta DROP COLUMN sourcegraphVersion`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE meta RENAME COLUMN numResultChunks TO num_result_chunks`); err != nil {
		return err
	}

	return nil
}
