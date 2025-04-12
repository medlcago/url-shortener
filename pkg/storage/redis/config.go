package redis

import "github.com/redis/go-redis/v9"

type Config struct {
	Client *redis.Client
}

func defaultConfig(cfg ...Config) Config {
	if len(cfg) > 0 {
		return cfg[0]
	}
	return defaultCfg
}

var defaultCfg = Config{
	Client: redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	}),
}
