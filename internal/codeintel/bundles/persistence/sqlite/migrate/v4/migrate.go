package v4

import (
	"context"

	"github.com/keegancsmith/sqlf"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/store"
)

// Migrate v4: Rename meta.numResultChunks to meta.num_result_chunks and drop the version columns.
func Migrate(ctx context.Context, s *store.Store, serializer serialization.Serializer) error {
	return s.ExecAll(
		ctx,
		sqlf.Sprintf(`CREATE TABLE t_meta (num_result_chunks int NOT NULL)`),
		sqlf.Sprintf(`INSERT INTO t_meta (num_result_chunks) SELECT numResultChunks FROM meta`),
		sqlf.Sprintf(`DROP TABLE meta`),
		sqlf.Sprintf(`ALTER TABLE t_meta RENAME TO meta`),
	)
}
