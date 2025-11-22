package config

import "os"

type Config struct {
	HTTPAddr    string
	JWTSecret   string
	Environment string

	MemberServiceAddr       string
	MatchingServiceAddr     string
	MessagingServiceAddr    string
	NotificationServiceAddr string
	MediaServiceAddr        string
	LocationServiceAddr     string
}

func Load() *Config {
	return &Config{
		HTTPAddr:    getEnv("HTTP_ADDR", ":8000"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		Environment: getEnv("ENVIRONMENT", "development"),

		MemberServiceAddr:       getEnv("MEMBER_SERVICE_ADDR", "localhost:9090"),
		MatchingServiceAddr:     getEnv("MATCHING_SERVICE_ADDR", "localhost:9091"),
		MessagingServiceAddr:    getEnv("MESSAGING_SERVICE_ADDR", "localhost:9092"),
		NotificationServiceAddr: getEnv("NOTIFICATION_SERVICE_ADDR", "localhost:9093"),
		MediaServiceAddr:        getEnv("MEDIA_SERVICE_ADDR", "localhost:9094"),
		LocationServiceAddr:     getEnv("LOCATION_SERVICE_ADDR", "localhost:9095"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
