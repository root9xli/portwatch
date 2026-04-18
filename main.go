package main

import (
	"flag"
	"log"
	"time"

	"github.com/portwatch/monitor"
)

func main() {
	interval := flag.Duration("interval", 10*time.Second, "polling interval")
	configFile := flag.String("config", "portwatch.yaml", "path to config file")
	flag.Parse()

	cfg, err := monitor.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	m := monitor.New(cfg)
	log.Printf("portwatch started (interval=%s)", *interval)

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := m.Run(); err != nil {
			log.Printf("monitor error: %v", err)
		}
	}
}
