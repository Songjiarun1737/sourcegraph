package v4

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization"
)

func Migrate(ctx context.Context, db *sqlx.DB, serializer serialization.Serializer) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE t_meta (num_result_chunks int NOT NULL)`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO t_meta (num_result_chunks) SELECT numResultChunks FROM meta`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `DROP TABLE meta`); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, `ALTER TABLE t_meta RENAME TO meta`); err != nil {
		return err
	}

	return nil
}
