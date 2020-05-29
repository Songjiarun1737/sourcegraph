package v3

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization"
)

func Migrate(ctx context.Context, db *sqlx.DB, serializer serialization.Serializer) error {
	_, err := db.ExecContext(ctx, `ALTER TABLE resultChunks RENAME TO result_chunks`)
	return err
}
