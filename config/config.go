package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Logger   Logger
	Server   Server
	Postgres Postgres
	Redis    Redis
}

type Logger struct {
	Development bool
	Encoding    string
	Level       string
}

type Server struct {
	AppVersion   string
	ServerHeader string
	Port         string
	MetricsPort  string
	ProxyHeader  string
	Mode         string
	JwtSecretKey string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Debug        bool
}

type Postgres struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlDbname   string
	PostgresqlPassword string
	PgDriver           string
}

type Redis struct {
	RedisAddr     string
	RedisPassword string
	RedisDb       int
	MinIdleConns  int
	PoolSize      int
	PoolTimeout   time.Duration
}

func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct: %v", err)
		return nil, err
	}

	return &c, nil
}
