package config

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST" default:"0.0.0.0"`
	Port     string `envconfig:"REDIS_PORT" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
}

type RMQ struct {
	UseTLS   bool   `envconfig:"AMQP_USE_TLS" default:"false"`
	Host     string `envconfig:"AMQP_HOST" default:"localhost"`
	Port     int    `envconfig:"AMQP_PORT" default:"5672"`
	User     string `envconfig:"AMQP_USER"`
	Password string `envconfig:"AMQP_PASSWORD"`
}

type SentryConfig struct {
	DSN           string  `envconfig:"SENTRY_DSN"`
	EnableTracing bool    `envconfig:"SENTRY_ENABLE_TRACING" default:"true"`
	SampleRate    float64 `envconfig:"SENTRY_SAMPLE_RATE" default:"0.2"`
}
