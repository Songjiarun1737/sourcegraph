package v2

import (
	"context"

	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/store"
)

// Migrate v2: Modify the storage of definition and references. Prior to this version, both
// tables stored scheme, identifier, and location fields as normalized rows. This version
// modifies the tables to store arrays of encoded locations keyed by (scheme, identifier)
// pairs. This makes the storage more uniform with documents and result chunks, and tends
// to save a good amount of space on disk due to the reduce number of tuples.
func Migrate(ctx context.Context, s *store.Store, serializer serialization.Serializer) error {
	for _, tableName := range []string{"definitions", "references"} {
		if err := migrateDefinitionReferences(ctx, s, serializer, tableName); err != nil {
			return err
		}
	}

	return nil
}
