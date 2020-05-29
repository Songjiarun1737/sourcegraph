package v2

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/types"
	"github.com/sourcegraph/sourcegraph/internal/sqliteutil"
)

func Migrate(ctx context.Context, db *sqlx.DB, serializer serialization.Serializer) error {
	if err := migrateDefinitionReferences(ctx, db, serializer, "definitions"); err != nil {
		return err
	}
	if err := migrateDefinitionReferences(ctx, db, serializer, "references"); err != nil {
		return err
	}

	return nil
}

type DefinitionReferenceRow struct {
	Scheme         string
	Identifier     string
	URI            string
	StartLine      int
	StartCharacter int
	EndLine        int
	EndCharacter   int
}

func migrateDefinitionReferences(ctx context.Context, db *sqlx.DB, serializer serialization.Serializer, tableName string) error {
	rows, err := scanDefinitionReferenceRows(db.QueryContext(ctx, `SELECT * FROM "`+tableName+`"`))
	if err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `CREATE TABLE "t_`+tableName+`" ("scheme" text NOT NULL, "identifier" text NOT NULL, "data" blob NOT NULL)`); err != nil {
		return err
	}

	inserter := sqliteutil.NewBatchInserter(db, "definitions", "scheme", "identifier", "data")

	for _, row := range groupDefinitionReferenceRows(rows) {
		data, err := serializer.MarshalLocations(row.Locations)
		if err != nil {
			return err
		}

		if err := inserter.Insert(ctx, row.Scheme, row.Identifier, data); err != nil {
			return err
		}
	}

	if err := inserter.Flush(ctx); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `DELETE TABLE "`+tableName+`"`); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, `RENAME TABLE "t_`+tableName+`" TO "`+tableName+`"`); err != nil {
		return err
	}

	return nil
}

// scanDefinitionReferenceRow populates a DefinitionReferenceRow value from the given scanner.
func scanDefinitionReferenceRow(rows *sql.Rows) (row DefinitionReferenceRow, err error) {
	err = rows.Scan(
		&row.Scheme,
		&row.Identifier,
		&row.URI,
		&row.StartLine,
		&row.StartCharacter,
		&row.EndLine,
		&row.EndCharacter,
	)
	return row, err
}

// scanDefinitionReferenceRows reads the given set of definition/reference rows and returns
// a slice of resulting values. This method should be called directly with the return value
// of `*db.query`.
func scanDefinitionReferenceRows(rows *sql.Rows, err error) ([]DefinitionReferenceRow, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var definitionReferenceRows []DefinitionReferenceRow
	for rows.Next() {
		row, err := scanDefinitionReferenceRow(rows)
		if err != nil {
			return nil, err
		}

		definitionReferenceRows = append(definitionReferenceRows, row)
	}

	return definitionReferenceRows, nil // TODO(efritz) - need to update rows.Err() everywhere
}

func groupDefinitionReferenceRows(rows []DefinitionReferenceRow) []types.MonikerLocations {
	uniques := map[string]types.MonikerLocations{}
	for _, row := range rows {
		key := makeKey(row.Scheme, row.Identifier)
		uniques[key] = types.MonikerLocations{
			Scheme:     row.Scheme,
			Identifier: row.Identifier,
			Locations: append(uniques[key].Locations, types.Location{
				URI:            row.URI,
				StartLine:      row.StartLine,
				StartCharacter: row.StartCharacter,
				EndLine:        row.EndLine,
				EndCharacter:   row.EndCharacter,
			}),
		}
	}

	monikerLocations := make([]types.MonikerLocations, 0, len(uniques))
	for _, v := range uniques {
		if len(v.Locations) > 0 {
			monikerLocations = append(monikerLocations, v)
		}
	}

	return monikerLocations
}

func makeKey(parts ...string) string {
	return strings.Join(parts, ":")
}
