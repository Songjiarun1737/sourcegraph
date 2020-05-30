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

var migrations = []struct {
	Version       string
	MigrationFunc func(ctx context.Context, s *store.Store, serializer serialization.Serializer) error
}{
	{"v00000", v0.Migrate},
	{"v00001", v1.Migrate},
	{"v00002", v2.Migrate},
	{"v00003", v3.Migrate},
	{"v00004", v4.Migrate},
}

var UnknownSchemaVersion = migrations[0].Version
var CurrentSchemaVersion = migrations[len(migrations)-1].Version

func Migrate(ctx context.Context, s *store.Store, serializer serialization.Serializer) error {
	version, err := getVersion(ctx, s)
	if err != nil {
		fmt.Printf("WAAS BAD?\n")
		return err
	}

	//
	// TODO - should copy file, replace, etc
	//

	found := false
	for _, migration := range migrations {
		if migration.Version == version {
			found = true
			continue
		}
		if !found {
			continue
		}

		if err := migration.MigrationFunc(ctx, s, serializer); err != nil {
			fmt.Printf("FAILY WHALEY %v\n", migration.Version)
			return err
		}
	}

	if !found {
		return fmt.Errorf("unrecognized schema version %s", version)
	}

	return s.ExecAll(ctx, sqlf.Sprintf("UPDATE schema_version SET version = %s", CurrentSchemaVersion))
}

func getVersion(ctx context.Context, s *store.Store) (string, error) {
	version, exists, err := store.ScanFirstString(s.Query(ctx, sqlf.Sprintf("SELECT version FROM schema_version LIMIT 1")))
	fmt.Printf("ERR: %v\n", err)
	if err != nil {
		// TODO - better matching
		if strings.Contains(err.Error(), "no such table: schema_version") {
			fmt.Printf("UHHHOK\n")
			return UnknownSchemaVersion, nil
		}

		return "", err
	}
	if !exists {
		return "", fmt.Errorf("No version in table")
	}

	return version, nil
}
