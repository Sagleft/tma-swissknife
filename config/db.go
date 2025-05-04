package config

import "fmt"

type DB struct {
	Host          string `envconfig:"DB_HOST" default:"localhost"`
	Port          int    `envconfig:"DB_PORT" default:"3306"`
	Name          string `envconfig:"DB_NAME"`
	User          string `envconfig:"DB_USER"`
	Password      string `envconfig:"DB_PASSWORD"`
	ConnTimeoutMS int    `envconfig:"DB_CONN_TIMEOUT" default:"5000"`

	MaxOpenConns        int `envconfig:"DB_MAX_OPEN_CONNS" default:"10"`
	MaxIdleConns        int `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetimeMins int `envconfig:"DB_CONN_MAX_LIFETIME_MINS" default:"5"`

	GormDebugMode bool `envconfig:"DB_GORM_DEBUG_MODE" default:"false"`
}

func (db DB) GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?timeout=%dms&parseTime=true",
		db.User, db.Password, db.Host, db.Port, db.Name, db.ConnTimeoutMS,
	)
}
