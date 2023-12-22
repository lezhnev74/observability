package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lezhnev74/observability"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	"time"
)

func main() {

	ctx := context.Background()
	cfg := loadCfg()

	ch, err := observability.ClickhouseConnect(
		cfg.Clickhouse.Host,
		cfg.Clickhouse.Port,
		cfg.Clickhouse.Database,
		cfg.Clickhouse.Username,
		cfg.Clickhouse.Password,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = observability.ClickhouseMigrate(ch, cfg.Clickhouse.Table)
	if err != nil {
		log.Fatal(err)
	}

	metricIngest := observability.ClickhouseInsert(ctx, ch, cfg.Clickhouse.Table)

	// Accept:
	go func() {
		path := fmt.Sprintf("%s:%s", cfg.Listen.Host, cfg.Listen.Port)
		pc, err := net.ListenPacket("udp", path)
		if err != nil {
			log.Fatalf("server failed: %s", err)
		}
		defer pc.Close()
		log.Printf("listen started")

		buf := make([]byte, 1024*1024)
		for {
			n, _, err := pc.ReadFrom(buf)
			if err != nil {
				log.Printf("unable to read incoming data: %s", err)
				break
			}
			if n == len(buf) {
				log.Printf("incoming data is too big")
			}

			m, err := observability.ParseMetric(buf[:n])
			if err != nil {
				log.Printf("unable to parse incoming data: %s", err)
				continue
			}
			metricIngest <- m
		}

		log.Printf("stop listening")
	}()

	for _, observe := range cfg.Observe {

		if observe.PollIntervalSec < 1 {
			log.Fatalf("poll interval is below 1s")
		}

		interval := time.Second * time.Duration(observe.PollIntervalSec)

		// FPM:
		if observe.FpmPoolStatsAddr != "" {
			f := func() {
				s, err := observability.PollFpmPoolStatus(observe.FpmPoolStatsAddr, observe.FpmPoolStatsPath)
				if err != nil {
					log.Print(err)
					return
				}
				for _, m := range observability.MapFpmToMetrics(*s) {
					metricIngest <- m
				}
			}
			observability.RepeatEvery(ctx, f, interval)
			continue
		}
	}

	<-ctx.Done()
}

func loadCfg() observability.Config {
	var cfgFile string
	flag.StringVar(&cfgFile, "c", "config.yml", "Specify the config file path")
	flag.Parse()

	cfgBytes, err := os.ReadFile(cfgFile)
	if err != nil {
		log.Fatalf("unable to read the config file as %s: %s", cfgFile, err)
	}

	var cfg observability.Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		log.Fatalf("unable to parse config: %s", err)
	}

	return cfg
}
