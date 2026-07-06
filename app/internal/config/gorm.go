package config

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(v *viper.Viper, log *zerolog.Logger) *gorm.DB {
	dsn := v.GetString("DATABASE_URL")

	idleConnection := v.GetInt("DATABASE_POOL_IDLE")
	maxConnection := v.GetInt("DATABASE_POOL_MAX")
	maxLifeTimeConnection := v.GetInt("DATABASE_POOL_LIFETIME")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(&zerologWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second * 2,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		}),
	})

	if err != nil {
		log.Fatal().Msgf("failed to connect database: %v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatal().Msgf("failed to connect database: %v", err)
	}

	connection.SetMaxIdleConns(idleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	log.Info().Msg("✅ Connection to MySQL established via GORM")

	return db
}

type zerologWriter struct {
	Logger *zerolog.Logger
}

func (l *zerologWriter) Printf(message string, args ...interface{}) {
	l.Logger.Info().Msgf(message, args...)
}
