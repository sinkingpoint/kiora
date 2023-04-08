package kioradb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/internal/stubs"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/sinkingpoint/msgpack/v5"
	"go.etcd.io/bbolt"
	"go.opentelemetry.io/otel"
)

var _ = DB(&BoltDB{})

// BoltDB is a DB implementation that stores data in a BoltDB database.
type BoltDB struct {
	db *bbolt.DB

	// cache is an in-memory cache of the database. This is used to speed up
	// queries to avoid having to serialize/deserialize data from the database for every request.
	cache *inMemoryDB
}

// NewBoltDB creates a new BoltDB database at the given path.
func NewBoltDB(path string) (*BoltDB, error) {
	backingDB, err := bbolt.Open(path, 0600, &bbolt.Options{
		Timeout:  1 * time.Second,
		OpenFile: stubs.OS.OpenFile,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to open bolt db")
	}

	db := &BoltDB{
		db:    backingDB,
		cache: NewInMemoryDB(),
	}

	db.refreshCache()

	return db, nil
}

// refreshCache clears out the in-memory cache and reloads it from the database.
func (b *BoltDB) refreshCache() error {
	log.Debug().Msg("loading boltdb into cache")
	b.cache.Clear()

	alerts := []model.Alert{}
	silences := []model.Silence{}
	if err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("alerts"))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var alert model.Alert
			if err := msgpack.Unmarshal(v, &alert); err != nil {
				return errors.Wrap(err, "failed to unmarshal alert")
			}

			alerts = append(alerts, alert)
			return nil
		})
	}); err != nil {
		return err
	}

	if err := b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("silences"))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var silence model.Silence
			if err := msgpack.Unmarshal(v, &silence); err != nil {
				return errors.Wrap(err, "failed to unmarshal silence")
			}

			silences = append(silences, silence)
			return nil
		})
	}); err != nil {
		return err
	}

	if err := b.cache.StoreAlerts(context.Background(), alerts...); err != nil {
		return err
	}

	if err := b.cache.StoreSilences(context.Background(), silences...); err != nil {
		return err
	}

	log.Debug().Msg("loaded boltdb into cache")

	return nil
}

func (b *BoltDB) StoreAlerts(ctx context.Context, alerts ...model.Alert) error {
	ctx, span := otel.Tracer("").Start(ctx, "BoltDB.StoreAlerts")
	defer span.End()

	if err := b.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("alerts"))
		if err != nil {
			return errors.Wrap(err, "failed to create alerts bucket")
		}

		for _, alert := range alerts {
			bytes, err := msgpack.Marshal(alert)
			if err != nil {
				return errors.Wrap(err, "failed to marshal alert")
			}

			if err := bucket.Put(alert.Labels.Bytes(), bytes); err != nil {
				return errors.Wrap(err, "failed to store alert")
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return b.cache.StoreAlerts(ctx, alerts...)
}

func (b *BoltDB) QueryAlerts(ctx context.Context, query *query.AlertQuery) []model.Alert {
	return b.cache.QueryAlerts(ctx, query)
}

func (b *BoltDB) StoreSilences(ctx context.Context, silences ...model.Silence) error {
	if err := b.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("silences"))
		if err != nil {
			return errors.Wrap(err, "failed to create silences bucket")
		}

		for _, silence := range silences {
			bytes, err := msgpack.Marshal(silence)
			if err != nil {
				return errors.Wrap(err, "failed to marshal silence")
			}

			if err := bucket.Put([]byte(silence.ID), bytes); err != nil {
				return errors.Wrap(err, "failed to store silence")
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return b.cache.StoreSilences(ctx, silences...)
}

// QuerySilences queries the database for silences matching the given query.
func (b *BoltDB) QuerySilences(ctx context.Context, query query.SilenceFilter) []model.Silence {
	return b.cache.QuerySilences(ctx, query)
}

func (b *BoltDB) Close() error {
	return b.db.Close()
}
