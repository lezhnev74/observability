package observability

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"log"
	"time"
)

// Clickhouse ingests data from various sources.
// It uses Metric data model for the data.
// It uses async insert, so the whole app works as a simple proxy for data.

const createTable = `
	CREATE TABLE observability
	(
    	timestamp DateTime('UTC') CODEC(Delta, ZSTD),
		metric LowCardinality(String),
		value Int64,
		tag1 String,
		tag2 String,
		tag3 String
	)
	ENGINE = MergeTree
	ORDER BY (metric, timestamp)
	TTL timestamp + INTERVAL 15 DAY;
`

type Metric struct {
	Timestamp time.Time
	Metric    string
	Value     int64
	Tag1      string
	Tag2      string
	Tag3      string
}

func chConnect(host, port, database, username, password string) (driver.Conn, error) {
	return clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", host, port)},
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
	})
}

// ingest starts a new connection to CH and returns the ingestion channel.
func ingest(ctx context.Context, ch driver.Conn) chan<- Metric {
	ingestCh := make(chan Metric, 1000)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ch.Close()
				return
			case m := <-ingestCh:
				insertQuery := fmt.Sprintf(
					`INSERT INTO observability VALUES(%d,'%s',%d,'%s','%s','%s')`,
					m.Timestamp.Unix(),
					m.Metric,
					m.Value,
					m.Tag1,
					m.Tag2,
					m.Tag3,
				)
				err := ch.AsyncInsert(ctx, insertQuery, false)
				if err != nil {
					log.Printf("unable to insert to clickhouse: %s", err)
				}
			}
		}
	}()
	return ingestCh
}
