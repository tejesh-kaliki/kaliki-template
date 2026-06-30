// Package config loads runtime configuration from a YAML file with env overrides.
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)
// MailConfig drives the transactional mail package. Provider selects the
// transport at runtime ("smtp" -> Mailpit/SMTP, "ses" -> AWS SES, anything else
// -> the log transport, which prints the message to the server log). Switching
// providers never changes whether credentials leak to clients — they don't.
type MailConfig struct {
	Provider      string `yaml:"provider"`
	SenderAddress string `yaml:"sender_address"`
	SenderName    string `yaml:"sender_name"`
	AppName       string `yaml:"app_name"`
	// BaseURL is the public app URL used to build links (token verification method).
	BaseURL string `yaml:"base_url"`
	SMTP SMTPConfig `yaml:"smtp"`
	SES  SESConfig  `yaml:"ses"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type SESConfig struct {
	Region          string `yaml:"region"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
}

type Config struct {
	Server        ServerConfig        `yaml:"server"`
	Database      DatabaseConfig      `yaml:"database"`
	Observability ObservabilityConfig `yaml:"observability"`
	Token         TokenConfig         `yaml:"token"`
	Mail          MailConfig          `yaml:"mail"`
	Redis         RedisConfig         `yaml:"redis"`
	Kafka         KafkaConfig         `yaml:"kafka"`
	Razorpay      RazorpayConfig      `yaml:"razorpay"`
	Firebase      FirebaseConfig      `yaml:"firebase"`
	Storage       StorageConfig       `yaml:"storage"`
}

type RedisConfig struct {
	URL string `yaml:"url"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}

type RazorpayConfig struct {
	KeyID     string `yaml:"key_id"`
	KeySecret string `yaml:"key_secret"`
}

type FirebaseConfig struct {
	CredentialsFile string `yaml:"credentials_file"`
}

type StorageConfig struct {
	Region   string `yaml:"region"`
	Bucket   string `yaml:"bucket"`
	Endpoint string `yaml:"endpoint"` // optional, for S3-compatible services
}

// TokenConfig holds JWT signing settings (used by the auth module).
type TokenConfig struct {
	Secret string `yaml:"secret"`
	// ExpiryHours is the access-token TTL. Keep it short; clients refresh.
	ExpiryHours int `yaml:"expiry_hours"`
	// RefreshExpiryHours is the refresh-token TTL (defaults to 720h / 30 days).
	RefreshExpiryHours int `yaml:"refresh_expiry_hours"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	URL string `yaml:"url"`
}

// ObservabilityConfig is wired in by default. Leaving Endpoint empty selects a
// no-op tracer, so the app runs with zero external dependencies until you point
// it at a real OTel collector.
type ObservabilityConfig struct {
	ServiceName string `yaml:"service_name"`
	Endpoint    string `yaml:"endpoint"`
}

// Load reads the YAML at path, then applies env overrides for anything that is
// commonly injected by the deployment environment.
func Load(path string) (*Config, error) {
	cfg := &Config{
		Server:        ServerConfig{Port: "8080"},
		Observability: ObservabilityConfig{ServiceName: "backend"},
		Token:         TokenConfig{ExpiryHours: 1, RefreshExpiryHours: 720},
	}

	if b, err := os.ReadFile(path); err == nil {
		if err := yaml.Unmarshal(b, cfg); err != nil {
			return nil, err
		}
	}

	if v := os.Getenv("SERVER_PORT"); v != "" {
		cfg.Server.Port = v
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.Database.URL = v
	}
	if v := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); v != "" {
		cfg.Observability.Endpoint = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.Token.Secret = v
	}
	if v := os.Getenv("MAIL_PROVIDER"); v != "" {
		cfg.Mail.Provider = v
	}
	if v := os.Getenv("MAIL_SENDER_ADDRESS"); v != "" {
		cfg.Mail.SenderAddress = v
	}
	if v := os.Getenv("MAIL_BASE_URL"); v != "" {
		cfg.Mail.BaseURL = v
	}
	if v := os.Getenv("SMTP_HOST"); v != "" {
		cfg.Mail.SMTP.Host = v
	}
	if v := os.Getenv("SMTP_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.Mail.SMTP.Port = p
		}
	}
	if v := os.Getenv("SMTP_USERNAME"); v != "" {
		cfg.Mail.SMTP.Username = v
	}
	if v := os.Getenv("SMTP_PASSWORD"); v != "" {
		cfg.Mail.SMTP.Password = v
	}
	if v := os.Getenv("SES_REGION"); v != "" {
		cfg.Mail.SES.Region = v
	}
	if v := os.Getenv("SES_ACCESS_KEY_ID"); v != "" {
		cfg.Mail.SES.AccessKeyID = v
	}
	if v := os.Getenv("SES_SECRET_ACCESS_KEY"); v != "" {
		cfg.Mail.SES.SecretAccessKey = v
	}
	if v := os.Getenv("REDIS_URL"); v != "" {
		cfg.Redis.URL = v
	}
	if v := os.Getenv("KAFKA_BROKERS"); v != "" {
		cfg.Kafka.Brokers = strings.Split(v, ",")
	}
	if v := os.Getenv("RAZORPAY_KEY_ID"); v != "" {
		cfg.Razorpay.KeyID = v
	}
	if v := os.Getenv("RAZORPAY_KEY_SECRET"); v != "" {
		cfg.Razorpay.KeySecret = v
	}
	if v := os.Getenv("FIREBASE_CREDENTIALS_FILE"); v != "" {
		cfg.Firebase.CredentialsFile = v
	}
	if v := os.Getenv("S3_BUCKET"); v != "" {
		cfg.Storage.Bucket = v
	}
	if v := os.Getenv("S3_REGION"); v != "" {
		cfg.Storage.Region = v
	}
	if v := os.Getenv("S3_ENDPOINT"); v != "" {
		cfg.Storage.Endpoint = v
	}

	// Fail fast if the insecure dev secret survived into a real deployment.
	if os.Getenv("APP_ENV") == "production" &&
		(cfg.Token.Secret == "" || cfg.Token.Secret == "dev-insecure-change-me") {
		return nil, errors.New("refusing to start: set a strong JWT_SECRET " +
			"(the dev default is not allowed when APP_ENV=production)")
	}

	return cfg, nil
}
