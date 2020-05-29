package v0

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization"
)

func Migrate(ctx context.Context, db *sqlx.DB, serializer serialization.Serializer) error {
	return nil
}
