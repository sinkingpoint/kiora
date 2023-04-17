package serf

import (
	"context"

	"github.com/hashicorp/serf/serf"
	"github.com/rs/zerolog"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"
	"github.com/sinkingpoint/msgpack/v5"
)

var _ = serf.UserDelegate(&DBDelegate{})

type DBDump struct {
	Alerts   []model.Alert
	Silences []model.Silence
}

type DBDelegate struct {
	db     kioradb.DB
	logger zerolog.Logger
}

func NewDBDelegate(db kioradb.DB, logger zerolog.Logger) *DBDelegate {
	return &DBDelegate{
		db:     db,
		logger: logger.With().Str("component", "db-delegate").Logger(),
	}
}

func (d *DBDelegate) LocalState(join bool) []byte {
	dump := DBDump{}
	dump.Alerts = d.db.QueryAlerts(context.Background(), query.NewAlertQuery(query.MatchAll()))
	dump.Silences = d.db.QuerySilences(context.Background(), query.MatchAll())

	bytes, _ := msgpack.Marshal(dump)
	return bytes
}

func (d *DBDelegate) MergeRemoteState(buf []byte, join bool) {
	dump := DBDump{}
	if err := msgpack.Unmarshal(buf, &dump); err != nil {
		d.logger.Err(err).Msg("failed to unmarshal DB dump")
		return
	}

	if err := d.db.StoreSilences(context.Background(), dump.Silences...); err != nil {
		d.logger.Err(err).Msg("failed to store silences")

		// We return here because if we failed to store silences, then submitting alerts may cause false positives.
		return
	}

	if err := d.db.StoreAlerts(context.Background(), dump.Alerts...); err != nil {
		d.logger.Err(err).Msg("failed to store alerts")
	}
}
