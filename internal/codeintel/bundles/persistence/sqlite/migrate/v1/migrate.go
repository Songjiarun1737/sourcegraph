package v1

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization"
)

func Migrate(ctx context.Context, db *sqlx.DB, serializer serialization.Serializer) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE schema_version ("version" TEXT NOT NULL)`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES (?);`, "v00001"); err != nil {
		return err
	}
	return nil
}
