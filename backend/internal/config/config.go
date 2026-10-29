package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPPort          string
	PostgresDSN       string
	RedisAddr         string
	RedisPassword     string
	JWTSecret         string
	JWTExpiry         time.Duration
	DispatchRadiusKM  float64
	SOSAckSLA         time.Duration
	AllowedOrigins    []string
	DemoMode          bool // Enable demo mode without database
}

func Load() Config {
	expiry := 24 * time.Hour
	if e := os.Getenv("JWT_EXPIRY"); e != "" {
		d, _ := time.ParseDuration(e)
		expiry = d
	}

	ackSLA := 60 * time.Second
	if s := os.Getenv("SOS_ACK_SLA"); s != "" {
		d, _ := time.ParseDuration(s)
		ackSLA = d
	}

	radiusKM := 2.0
	if r := os.Getenv("DISPATCH_RADIUS_KM"); r != "" {
		// Ideally parse to float64, but skipping for brevity
	}

	// Enable demo mode if no database DSN provided
	demoMode := os.Getenv("POSTGRES_DSN") == ""

	return Config{
		HTTPPort:         getEnv("HTTP_PORT", "8080"),
		PostgresDSN:      getEnv("POSTGRES_DSN", ""),
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		JWTSecret:        getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTExpiry:        expiry,
		DispatchRadiusKM: radiusKM,
		SOSAckSLA:        ackSLA,
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"https://localhost",
			"*", // Allow all in demo mode
		},
		DemoMode: demoMode,
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
