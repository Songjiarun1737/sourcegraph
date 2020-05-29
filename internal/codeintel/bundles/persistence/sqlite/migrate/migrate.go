package migrate

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization"
	v0 "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/migrate/v0"
	v1 "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/migrate/v1"
	v2 "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/migrate/v2"
	v3 "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/migrate/v3"
	v4 "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/migrate/v4"
)

var migrations = []struct {
	Version       string
	MigrationFunc func(ctx context.Context, db *sqlx.DB, serializer serialization.Serializer) error
}{
	{"v00000", v0.Migrate},
	{"v00001", v1.Migrate},
	{"v00002", v2.Migrate},
	{"v00003", v3.Migrate},
	{"v00004", v4.Migrate},
}

var UnknownSchemaVersion = migrations[0].Version
var CurrentSchemaVersion = migrations[len(migrations)-1].Version

func Migrate(ctx context.Context, db *sqlx.DB, serializer serialization.Serializer) error {
	version, err := getVersion(ctx, db)
	if err != nil {
		return err
	}

	// TODO - should copy file, replace, etc

	found := false
	for _, migration := range migrations {
		if migration.Version == version {
			found = true
			continue
		}
		if !found {
			continue
		}

		if err := migration.MigrationFunc(ctx, db, serializer); err != nil {
			return err
		}
	}

	if !found {
		return fmt.Errorf("unrecognized schema version %s", version)
	}

	if _, err := db.ExecContext(ctx, "UPDATE schema_version SET version = ?", CurrentSchemaVersion); err != nil {
		return err
	}

	return nil
}

func getVersion(ctx context.Context, db *sqlx.DB) (version string, _ error) {
	if err := db.QueryRowContext(ctx, "SELECT version FROM schema_version LIMIT 1").Scan(&version); err != nil {
		// TODO - better matching
		if strings.Contains(err.Error(), "no such table: schema_version") {
			return UnknownSchemaVersion, nil
		}

		return "", err
	}

	return version, nil
}
