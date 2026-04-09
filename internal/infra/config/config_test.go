package config

import (
	"os"
	"testing"
)

func TestLoadRequiresMongoAndKafka(t *testing.T) {
	t.Setenv("MONGO_URI", "")
	t.Setenv("KAFKA_BROKERS", "")

	_, err := Load()
	if err == nil {
		t.Fatalf("expected error when MONGO_URI and KAFKA_BROKERS are missing")
	}

	t.Setenv("MONGO_URI", "mongodb://x:27017")
	t.Setenv("KAFKA_BROKERS", "kafka:9092")
	_, err = Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestLoadSplitsKafkaBrokers(t *testing.T) {
	t.Setenv("MONGO_URI", "mongodb://x:27017")
	t.Setenv("KAFKA_BROKERS", "a:1, b:2,,  c:3 ")
	t.Setenv("MONGO_DB", "db")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.KafkaBrokers) != 3 {
		t.Fatalf("expected 3 brokers, got %d", len(cfg.KafkaBrokers))
	}
}

func TestLoadDefaultHTTPAddr(t *testing.T) {
	_ = os.Unsetenv("HTTP_ADDR")
	t.Setenv("MONGO_URI", "mongodb://x:27017")
	t.Setenv("KAFKA_BROKERS", "kafka:9092")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.HTTPAddr != "0.0.0.0:8080" {
		t.Fatalf("unexpected default HTTP addr: %s", cfg.HTTPAddr)
	}
}

