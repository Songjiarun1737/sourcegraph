package migrate

import (
	"context"
	"fmt"
	"strings"

	"github.com/keegancsmith/sqlf"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization"
	v0 "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/migrate/v0"
	v1 "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/migrate/v1"
	v2 "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/migrate/v2"
	v3 "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/migrate/v3"
	v4 "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/migrate/v4"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/store"
)

type MigrationFunc func(ctx context.Context, s *store.Store, serializer serialization.Serializer) error

type MigrationSpec struct {
	Version       string
	MigrationFunc MigrationFunc
}

var migrations = []MigrationSpec{
	{"v00000", v0.Migrate},
	{"v00001", v1.Migrate},
	{"v00002", v2.Migrate},
	{"v00003", v3.Migrate},
	{"v00004", v4.Migrate},
}

var UnknownSchemaVersion = migrations[0].Version
var CurrentSchemaVersion = migrations[len(migrations)-1].Version

func Migrate(ctx context.Context, s *store.Store, serializer serialization.Serializer) (err error) {
	version, err := getVersion(ctx, s)
	if err != nil {
		return err
	}

	found := false
	for _, migration := range migrations {
		if migration.Version == version {
			found = true
			continue
		}
		if !found {
			continue
		}

		if err := runMigration(ctx, s, serializer, migration.Version, migration.MigrationFunc); err != nil {
			return err
		}
	}

	if !found {
		return fmt.Errorf("unrecognized schema version %s", version)
	}

	return nil
}

func runMigration(ctx context.Context, store *store.Store, serializer serialization.Serializer, version string, migrationFunc MigrationFunc) (err error) {
	tx, err := store.Transact(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Done(err)
	}()

	if err := migrationFunc(ctx, tx, serializer); err != nil {
		return err
	}

	if err := tx.ExecAll(ctx, sqlf.Sprintf("UPDATE schema_version SET version = %s", version)); err != nil {
		return err
	}

	return nil
}

func getVersion(ctx context.Context, s *store.Store) (string, error) {
	version, exists, err := store.ScanFirstString(s.Query(ctx, sqlf.Sprintf("SELECT version FROM schema_version LIMIT 1")))
	if err != nil {
		// TODO - better matching
		if strings.Contains(err.Error(), "no such table: schema_version") {
			return UnknownSchemaVersion, nil
		}

		return "", err
	}
	if !exists {
		return "", fmt.Errorf("No version in table")
	}

	return version, nil
}
