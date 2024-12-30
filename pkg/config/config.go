package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Abhishekjha321/community_service/log"
	"github.com/Abhishekjha321/community_service/storage/db/postgres"
	"github.com/spf13/viper"
)

const (
	FilePath = "CONFIG_PATH"
	FileName = "CONFIG_FILE"
)

type RedisConfig struct {
	HostNames string
	Password  string
}

type Logger struct {
	Filename string
}
type Server struct {
	Port int
}

var Config = &struct {
	Name                     string
	AppEnv                   string
	AppVersion               string
	OtelExporterOtlpEndPoint string
	Server                   Server
	PostgresMaster           *postgres.PGMaster
	PostgresSlave            *postgres.PGSlave
	Logger                   Logger
	RedisConfig              RedisConfig
	UserInfoDelay time.Duration
}{}

func Initialize() error {
	configPath, ok := os.LookupEnv(FilePath)
	if !ok {
		return fmt.Errorf("error: config path env variable not found, env var: %s", FilePath)
	}

	configFileName, ok := os.LookupEnv(FileName)
	if !ok {
		return fmt.Errorf("error: config file env variable not found, env var: %s", FileName)
	}
	logger.Initialize(configFileName, Config.Name)

	// read config file
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configFileName)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return fmt.Errorf("error while reading config file: %w", err)
	}

	// bind config file to Config struct
	configBindErr := viper.Unmarshal(&Config)
	if configBindErr != nil {
		return fmt.Errorf("error when binding config: %w", configBindErr)
	}

	logger.GetLogger().WithContext(context.Background()).Info("Logger setup completed")

	otelExporterOtlpEndPoint := Config.OtelExporterOtlpEndPoint
	logger.GetLogInstance(context.Background(), "Initialize:").Infof(" otelExporterOtlpEndPoint: %v", otelExporterOtlpEndPoint)
	// telemetry.Initialize(Config.Name, otelExporterOtlpEndPoint, Config.AppEnv, Config.AppVersion)

	logger.GetLogger().WithContext(context.Background()).Info("Postgres initialized successfully")
	return nil
}
