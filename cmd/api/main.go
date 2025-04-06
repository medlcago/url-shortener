package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
	"url-shortener/config"
	"url-shortener/internal/server"
	"url-shortener/pkg/db/postgres"
	"url-shortener/pkg/db/redis"
	"url-shortener/pkg/http"
	"url-shortener/pkg/logger"
	"url-shortener/pkg/storage"
	"url-shortener/pkg/utils"
)

func main() {
	configPath := utils.GetConfigPath(os.Getenv("config_env"))
	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}
	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewLogrusLogger(cfg)
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	app := fiber.New(fiber.Config{
		StructValidator: utils.NewStructValidator(validator.New()),
		ReadTimeout:     time.Second * cfg.Server.ReadTimeout,
		WriteTimeout:    time.Second * cfg.Server.WriteTimeout,
		ServerHeader:    cfg.Server.ServerHeader,
		ErrorHandler:    http.ErrorHandler,
	})

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("NewPsqlDB: %v", err)
	}
	defer psqlDB.Close()

	rdb, err := redis.NewRedisClient(cfg)
	if err != nil {
		appLogger.Fatalf("NewRedisClient: %v", err)
	}
	defer rdb.Close()

	redisStorage := storage.WithRedisClient(rdb)

	srv := server.NewServer(app, cfg, psqlDB, redisStorage, appLogger)
	if err = srv.Run(); err != nil {
		appLogger.Fatalf("srv.Run: %v", err)
	}
}
