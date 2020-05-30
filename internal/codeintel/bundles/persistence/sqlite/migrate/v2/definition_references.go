package v2

import (
	"context"
	"fmt"

	"github.com/keegancsmith/sqlf"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/store"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/types"
	"github.com/sourcegraph/sourcegraph/internal/sqliteutil"
)

func migrateDefinitionReferences(ctx context.Context, s *store.Store, serializer serialization.Serializer, tableName string) error {
	tempTableName := fmt.Sprintf("t_%s", tableName)

	monikerLocations, err := groupOldData(ctx, s, tableName)
	if err != nil {
		return err
	}
	if err := createTempTable(ctx, s, tempTableName); err != nil {
		return err
	}
	if err := populateTable(ctx, s, serializer, tempTableName, monikerLocations); err != nil {
		return err
	}
	if err := replaceTable(ctx, s, tableName, tempTableName); err != nil {
		return err
	}

	return nil
}

func groupOldData(ctx context.Context, s *store.Store, tableName string) ([]types.MonikerLocations, error) {
	rows, err := scanDefinitionReferenceRows(s.Query(ctx, sqlf.Sprintf(`
		SELECT
			scheme,
			identifier,
			documentPath,
			startLine,
			startCharacter,
			endLine,
			endCharacter
		FROM "`+tableName+`"
	`)))
	if err != nil {
		return nil, err
	}

	return groupDefinitionReferenceRows(rows), nil
}

func createTempTable(ctx context.Context, s *store.Store, tempTableName string) error {
	return s.ExecAll(ctx, sqlf.Sprintf(`
		CREATE TABLE "`+tempTableName+`" (
			"scheme" text NOT NULL,
			"identifier" text NOT NULL,
			"data" blob NOT NULL
		)
	`))
}

func populateTable(ctx context.Context, s sqliteutil.Execable, serializer serialization.Serializer, tempTableName string, monikerLocations []types.MonikerLocations) error {
	inserter := sqliteutil.NewBatchInserter(s, tempTableName, "scheme", "identifier", "data")

	for _, ml := range monikerLocations {
		data, err := serializer.MarshalLocations(ml.Locations)
		if err != nil {
			return err
		}

		if err := inserter.Insert(ctx, ml.Scheme, ml.Identifier, data); err != nil {
			return err
		}
	}

	return inserter.Flush(ctx)
}

func replaceTable(ctx context.Context, s *store.Store, targetTable, tempTable string) error {
	return s.ExecAll(
		ctx,
		sqlf.Sprintf(`DROP TABLE "`+targetTable+`"`),
		sqlf.Sprintf(`ALTER TABLE "`+tempTable+`" RENAME TO "`+targetTable+`"`),
	)
}
