package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	R2_Endpoint        string        `mapstructure:"R2_ENDPOINT"`
	R2_Bucket          string        `mapstructure:"R2_BUCKET"`
	R2_Region          string        `mapstructure:"R2_REGION"`
	R2_Access_Key      string        `mapstructure:"R2_ACCESS_KEY"`
	R2_Secret_Key      string        `mapstructure:"R2_SECRET_KEY"`
	R2_Upload_Expiry   time.Duration `mapstructure:"R2_UPLOAD_EXPIRY_MINUTES"`
	R2_Download_Expiry time.Duration `mapstructure:"R2_DOWNLOAD_EXPIRY_MINUTES"`
	Cache_Dir          string        `mapstructure:"CACHE_DIR"`
}

var Configs Config

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	viper.BindEnv("R2_ENDPOINT")
	viper.BindEnv("R2_BUCKET")
	viper.BindEnv("R2_REGION")
	viper.BindEnv("R2_ACCESS_KEY")
	viper.BindEnv("R2_SECRET_KEY")
	viper.BindEnv("R2_UPLOAD_EXPIRY_MINUTES")
	viper.BindEnv("R2_DOWNLOAD_EXPIRY_MINUTES")
	viper.BindEnv("CACHE_DIR")

	viper.SetDefault("R2_UPLOAD_EXPIRY_MINUTES", "30m")
	viper.SetDefault("R2_DOWNLOAD_EXPIRY_MINUTES", "30m")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found or unreadable: %v", err)
	}

	if err := viper.Unmarshal(&Configs); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	log.Printf("Configuration loaded: %+v", Configs)
}
