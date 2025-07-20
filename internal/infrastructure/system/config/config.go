package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr      string
	DSN           string
	JWTSecret     string
	JWTTTL        time.Duration
	ArgonTime     uint32
	ArgonMemoryMB uint32
	ArgonThreads  uint8
	ArgonKeyLen   uint32
}

func Load() (Config, error) {
	var cfg Config
	var err error

	cfg.HTTPAddr = getEnv("APP_HTTP_ADDR", ":8080")
	cfg.DSN = os.Getenv("PG_URL")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		return cfg, errors.New("JWT_SECRET required")
	}

	ttlStr := getEnv("JWT_TTL", "1h")
	cfg.JWTTTL, err = time.ParseDuration(ttlStr)
	if err != nil {
		return cfg, errors.New("invalid JWT_TTL")
	}

	cfg.ArgonTime = uint32(mustUint("ARGON_TIME", 1))
	cfg.ArgonMemoryMB = uint32(mustUint("ARGON_MEMORY_MB", 64))
	cfg.ArgonThreads = uint8(mustUint("ARGON_THREADS", 2))
	cfg.ArgonKeyLen = uint32(mustUint("ARGON_KEY_LEN", 32))

	if cfg.DSN == "" {
		return cfg, errors.New("PG_URL required")
	}
	return cfg, nil
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func mustUint(k string, def uint64) uint64 {
	v := getEnv(k, "")
	if v == "" {
		return def
	}
	n, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return def
	}
	return n
}
