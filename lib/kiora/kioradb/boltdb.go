package kioradb

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sinkingpoint/kiora/internal/stubs"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/sinkingpoint/msgpack/v5"
	"go.etcd.io/bbolt"
)

var _ = DB(&BoltDB{})

type BoltDB struct {
	db *bbolt.DB
}

func NewBoltDB(path string) (*BoltDB, error) {
	db, err := bbolt.Open(path, 0600, &bbolt.Options{
		Timeout:  1 * time.Second,
		OpenFile: stubs.OS.OpenFile,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to open bolt db")
	}

	return &BoltDB{db: db}, nil
}

func (b *BoltDB) StoreAlerts(ctx context.Context, alerts ...model.Alert) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
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
	})
}

func (b *BoltDB) QueryAlerts(ctx context.Context, query *query.AlertQuery) []model.Alert {
	alerts := make([]model.Alert, 0)
	b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("alerts"))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var alert model.Alert
			if err := msgpack.Unmarshal(v, &alert); err != nil {
				return errors.Wrap(err, "failed to unmarshal alert")
			}

			if query.Filter.MatchesAlert(ctx, &alert) {
				alerts = append(alerts, alert)
			}

			return nil
		})
	})

	return alerts
}

func (b *BoltDB) StoreSilences(ctx context.Context, silences ...model.Silence) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
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
	})
}

// QuerySilences queries the database for silences matching the given query.
func (b *BoltDB) QuerySilences(ctx context.Context, query query.SilenceFilter) []model.Silence {
	silences := make([]model.Silence, 0)
	b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("silences"))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			var silence model.Silence
			if err := msgpack.Unmarshal(v, &silence); err != nil {
				return errors.Wrap(err, "failed to unmarshal silence")
			}

			if query.MatchesSilence(ctx, &silence) {
				silences = append(silences, silence)
			}

			return nil
		})
	})

	return silences
}

func (b *BoltDB) Close() error {
	return b.db.Close()
}
