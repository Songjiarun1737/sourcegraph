package sqlite

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/keegancsmith/sqlf"
	pkgerrors "github.com/pkg/errors"
	persistence "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization"
	jsonserializer "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization/json"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/migrate"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/store"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/types"
)

var ErrNoMetadata = errors.New("no rows in meta table")

type sqliteReader struct {
	s          *store.Store
	serializer serialization.Serializer
	close      func() error
}

var _ persistence.Reader = &sqliteReader{}

func NewReader(filename string) (persistence.Reader, error) {
	db, err := sqlx.Open("sqlite3_with_pcre", filename)
	if err != nil {
		return nil, err
	}
	// TODO - close on error

	s := store.New(db)
	serializer := jsonserializer.New()

	if err := migrate.Migrate(context.Background(), s, serializer); err != nil {
		fmt.Printf("Cannot do it mannnnnnn\n")
		return nil, err
	}

	return &sqliteReader{
		s:          s,
		serializer: serializer,
		close:      db.Close,
	}, nil
}

func (r *sqliteReader) ReadMeta(ctx context.Context) (types.MetaData, error) {
	query := `SELECT num_result_chunks FROM meta LIMIT 1`

	numResultChunks, exists, err := store.ScanFirstInt(r.s.Query(ctx, sqlf.Sprintf(query)))
	if err != nil {
		return types.MetaData{}, err
	}
	if !exists {
		return types.MetaData{}, ErrNoMetadata
	}

	return types.MetaData{
		NumResultChunks: numResultChunks,
	}, nil
}

func (r *sqliteReader) ReadDocument(ctx context.Context, path string) (types.DocumentData, bool, error) {
	query := `SELECT data FROM documents WHERE path = %s LIMIT 1`

	data, exists, err := store.ScanFirstBytes(r.s.Query(ctx, sqlf.Sprintf(query, path)))
	if err != nil || !exists {
		return types.DocumentData{}, false, err
	}

	documentData, err := r.serializer.UnmarshalDocumentData(data)
	if err != nil {
		return types.DocumentData{}, false, pkgerrors.Wrap(err, "serializer.UnmarshalDocumentData")
	}
	return documentData, true, nil
}

func (r *sqliteReader) ReadResultChunk(ctx context.Context, id int) (types.ResultChunkData, bool, error) {
	query := `SELECT data FROM result_chunks WHERE id = %s LIMIT 1`

	data, exists, err := store.ScanFirstBytes(r.s.Query(ctx, sqlf.Sprintf(query, id)))
	if err != nil || !exists {
		return types.ResultChunkData{}, false, err
	}

	resultChunkData, err := r.serializer.UnmarshalResultChunkData(data)
	if err != nil {
		return types.ResultChunkData{}, false, pkgerrors.Wrap(err, "serializer.UnmarshalResultChunkData")
	}
	return resultChunkData, true, nil
}

func (r *sqliteReader) ReadDefinitions(ctx context.Context, scheme, identifier string, skip, take int) ([]types.Location, int, error) {
	return r.readDefinitionReferences(ctx, "definitions", scheme, identifier, skip, take)
}

func (r *sqliteReader) ReadReferences(ctx context.Context, scheme, identifier string, skip, take int) ([]types.Location, int, error) {
	return r.readDefinitionReferences(ctx, "references", scheme, identifier, skip, take)
}

func (r *sqliteReader) readDefinitionReferences(ctx context.Context, tableName, scheme, identifier string, skip, take int) ([]types.Location, int, error) {
	query := `SELECT data FROM "` + tableName + `" WHERE scheme = %s AND identifier = %s LIMIT 1`

	data, exists, err := store.ScanFirstBytes(r.s.Query(ctx, sqlf.Sprintf(query, scheme, identifier)))
	if err != nil || !exists {
		return nil, 0, err
	}

	locations, err := r.serializer.UnmarshalLocations(data)
	if err != nil {
		return nil, 0, pkgerrors.Wrap(err, "serializer.UnmarshalLocations")
	}

	//
	// TODO - refactor this all nice
	//

	slicedLocations := locations
	if skip != 0 && take != 0 {
		if skip >= len(locations) {
			skip = len(locations)
		}
		max := skip + take
		if max > len(locations) {
			max = len(locations)
		}
		slicedLocations = slicedLocations[skip:max]
	}

	return slicedLocations, len(locations), err
}

func (r *sqliteReader) Close() error {
	return r.close()
}

// // query performs QueryContext on the underlying connection.
// func (r *sqliteReader) query(ctx context.Context, query *sqlf.Query) (*sql.Rows, error) {
// 	return r.s.QueryContext(ctx, query.Query(sqlf.SimpleBindVar), query.Args()...)
// }

// // queryRow performs QueryRowContext on the underlying connection.
// func (r *sqliteReader) queryRow(ctx context.Context, query *sqlf.Query) *sql.Row {
// 	return r.db.QueryRowContext(ctx, query.Query(sqlf.SimpleBindVar), query.Args()...)
// }

// // scanBytes populates a byte slice value from the given scanner.
// func scanBytes(scanner *sql.Row) (value []byte, err error) {
// 	err = scanner.Scan(&value)
// 	return value, err
// }

// // scanInt populates an int value from the given scanner.
// func scanInt(scanner *sql.Row) (value int, err error) {
// 	err = scanner.Scan(&value)
// 	return value, err
// }
