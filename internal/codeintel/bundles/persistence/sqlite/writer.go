package sqlite

import (
	"context"
	"database/sql"

	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	persistence "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization"
	jsonserializer "github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/serialization/json"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/persistence/sqlite/migrate"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/bundles/types"
	"github.com/sourcegraph/sourcegraph/internal/sqliteutil"
)

type sqliteWriter struct {
	db                    *sqlx.DB
	tx                    *sql.Tx
	serializer            serialization.Serializer
	scheamVersionInserter *sqliteutil.BatchInserter
	metaInserter          *sqliteutil.BatchInserter
	documentInserter      *sqliteutil.BatchInserter
	resultChunkInserter   *sqliteutil.BatchInserter
	definitionInserter    *sqliteutil.BatchInserter
	referenceInserter     *sqliteutil.BatchInserter
}

var _ persistence.Writer = &sqliteWriter{}

func NewWriter(filename string) (_ persistence.Writer, err error) {
	db, err := sqlx.Open("sqlite3_with_pcre", filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if closeErr := db.Close(); closeErr != nil {
				err = multierror.Append(err, closeErr)
			}
		}
	}()

	if _, err := db.Exec(`CREATE TABLE "schema_version" ("version" text NOT NULL)`); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE "meta" ("num_result_chunks" integer NOT NULL)`); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE "documents" ("path" text PRIMARY KEY NOT NULL, "data" blob NOT NULL)`); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE "result_chunks" ("id" integer PRIMARY KEY NOT NULL, "data" blob NOT NULL)`); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE "definitions" ("scheme" text NOT NULL, "identifier" text NOT NULL, "data" blob NOT NULL, PRIMARY KEY (scheme, identifier))`); err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE "references" ("scheme" text NOT NULL, "identifier" text NOT NULL, "data" blob NOT NULL, PRIMARY KEY (scheme, identifier))`); err != nil {
		return nil, err
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &sqliteWriter{
		db:                    db,
		tx:                    tx,
		serializer:            jsonserializer.New(),
		scheamVersionInserter: sqliteutil.NewBatchInserter(tx, "schema_version", "version"),
		metaInserter:          sqliteutil.NewBatchInserter(tx, "meta", "numResultChunks"),
		documentInserter:      sqliteutil.NewBatchInserter(tx, "documents", "path", "data"),
		resultChunkInserter:   sqliteutil.NewBatchInserter(tx, "result_chunks", "id", "data"),
		definitionInserter:    sqliteutil.NewBatchInserter(tx, "definitions", "scheme", "identifier", "data"),
		referenceInserter:     sqliteutil.NewBatchInserter(tx, `references`, "scheme", "identifier", "data"),
	}, nil
}

func (w *sqliteWriter) WriteMeta(ctx context.Context, numResultChunks int) error {
	if err := w.scheamVersionInserter.Insert(ctx, migrate.CurrentSchemaVersion); err != nil {
		return errors.Wrap(err, "scheamVersionInserter.Insert")
	}
	if err := w.metaInserter.Insert(ctx, numResultChunks); err != nil {
		return errors.Wrap(err, "metaInserter.Insert")
	}
	return nil
}

func (w *sqliteWriter) WriteDocuments(ctx context.Context, documents map[string]types.DocumentData) error {
	for path, document := range documents {
		data, err := w.serializer.MarshalDocumentData(document)
		if err != nil {
			return errors.Wrap(err, "serializer.MarshalDocumentData")
		}

		if err := w.documentInserter.Insert(ctx, path, data); err != nil {
			return errors.Wrap(err, "documentInserter.Insert")
		}
	}
	return nil
}

func (w *sqliteWriter) WriteResultChunks(ctx context.Context, resultChunks map[int]types.ResultChunkData) error {
	for id, resultChunk := range resultChunks {
		data, err := w.serializer.MarshalResultChunkData(resultChunk)
		if err != nil {
			return errors.Wrap(err, "serializer.MarshalResultChunkData")
		}

		if err := w.resultChunkInserter.Insert(ctx, id, data); err != nil {
			return errors.Wrap(err, "resultChunkInserter.Insert")
		}
	}
	return nil
}

func (w *sqliteWriter) WriteDefinitions(ctx context.Context, monikerLocations []types.MonikerLocations) error {
	for _, ml := range monikerLocations {
		data, err := w.serializer.MarshalLocations(ml.Locations)
		if err != nil {
			return errors.Wrap(err, "serializer.MarshalLocations")
		}

		if err := w.definitionInserter.Insert(ctx, ml.Scheme, ml.Identifier, data); err != nil {
			return errors.Wrap(err, "definitionInserter.Insert")
		}
	}
	return nil
}

func (w *sqliteWriter) WriteReferences(ctx context.Context, monikerLocations []types.MonikerLocations) error {
	for _, ml := range monikerLocations {
		data, err := w.serializer.MarshalLocations(ml.Locations)
		if err != nil {
			return errors.Wrap(err, "serializer.MarshalLocations")
		}

		if err := w.referenceInserter.Insert(ctx, ml.Scheme, ml.Identifier, data); err != nil {
			return errors.Wrap(err, "referenceInserter.Insert")
		}
	}
	return nil
}

func (w *sqliteWriter) Flush(ctx context.Context) error {
	if err := w.scheamVersionInserter.Flush(ctx); err != nil {
		return errors.Wrap(err, "scheamVersionInserter.Flush")
	}
	if err := w.metaInserter.Flush(ctx); err != nil {
		return errors.Wrap(err, "metaInserter.Flush")
	}
	if err := w.documentInserter.Flush(ctx); err != nil {
		return errors.Wrap(err, "documentInserter.Flush")
	}
	if err := w.resultChunkInserter.Flush(ctx); err != nil {
		return errors.Wrap(err, "resultChunkInserter.Flush")
	}
	if err := w.definitionInserter.Flush(ctx); err != nil {
		return errors.Wrap(err, "definitionInserter.Flush")
	}
	if err := w.referenceInserter.Flush(ctx); err != nil {
		return errors.Wrap(err, "referenceInserter.Flush")
	}
	if err := w.tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (w *sqliteWriter) Close() (err error) {
	return w.db.Close()
}
