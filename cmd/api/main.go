package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
	"url-shortener/config"
	_ "url-shortener/docs"
	"url-shortener/internal/server"
	"url-shortener/pkg/db/postgres"
	"url-shortener/pkg/db/redis"
	"url-shortener/pkg/http"
	_ "url-shortener/pkg/http"
	"url-shortener/pkg/logger"
	redisStorage "url-shortener/pkg/storage/redis"
	"url-shortener/pkg/utils"
)

// @title URL Shortener API
// @version 1.0

// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
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
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	app := fiber.New(fiber.Config{
		StructValidator: utils.NewStructValidator(validator.New()),
		ReadTimeout:     time.Second * cfg.Server.ReadTimeout,
		WriteTimeout:    time.Second * cfg.Server.WriteTimeout,
		ServerHeader:    cfg.Server.ServerHeader,
		ErrorHandler:    http.ErrorHandler,
		ProxyHeader:     cfg.Server.ProxyHeader,
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

	rs := redisStorage.New(redisStorage.Config{Client: rdb})
	defer rs.Close()

	srv := server.NewServer(app, cfg, psqlDB, rs, appLogger)
	if err = srv.Run(); err != nil {
		appLogger.Fatalf("srv.Run: %v", err)
	}
}
