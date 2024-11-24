package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	R2_Endpoint        string        `yaml:"R2_ENDPOINT"`
	R2_Bucket          string        `yaml:"R2_BUCKET"`
	R2_Region          string        `yaml:"R2_REGION"`
	R2_Access_Key      string        `yaml:"R2_ACCESS_KEY"`
	R2_Secret_Key      string        `yaml:"R2_SECRET_KEY"`
	R2_Upload_Expiry   time.Duration `yaml:"R2_UPLOAD_EXPIRY_MINUTES"`
	R2_Download_Expiry time.Duration `yaml:"R2_DOWNLOAD_EXPIRY_MINUTES"`
	Cache_Dir          string        `yaml:"CACHE_DIR"`
}

var Configs Config

func Init() {
	result, err := readEnviromentVars()
	if err != nil {
		log.Printf("Error reading enviroment vars: %v", err)
		result, err = readConfigFile()
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}
	}
	Configs = result
}

func readEnviromentVars() (Config, error) {
	Configs = Config{}
	var err error
	Configs.R2_Endpoint, err = readEnvThrowEmpty("R2_ENDPOINT")
	if err != nil {
		return Configs, err
	}
	Configs.R2_Bucket, err = readEnvThrowEmpty("R2_BUCKET")
	if err != nil {
		return Configs, err
	}
	Configs.R2_Region, err = readEnvThrowEmpty("R2_REGION")
	if err != nil {
		return Configs, err
	}
	Configs.R2_Access_Key, err = readEnvThrowEmpty("R2_ACCESS_KEY")
	if err != nil {
		return Configs, err
	}
	Configs.R2_Secret_Key, err = readEnvThrowEmpty("R2_SECRET_KEY")
	if err != nil {
		return Configs, err
	}
	Configs.R2_Upload_Expiry, err = envToDuration("R2_UPLOAD_EXPIRY_MINUTES")
	if err != nil {
		return Configs, err
	}
	Configs.R2_Download_Expiry, err = envToDuration("R2_DOWNLOAD_EXPIRY_MINUTES")
	if err != nil {
		return Configs, err
	}
	Configs.Cache_Dir, err = readEnvThrowEmpty("CACHE_DIR")
	if err != nil {
		return Configs, err
	}

	return Configs, nil
}

func envToDuration(key string) (time.Duration, error) {
	str, err := readEnvThrowEmpty(key)
	if err != nil {
		return time.Minute, err
	}
	var mins int
	mins, err = strconv.Atoi(str)
	if err != nil {
		return time.Minute, err
	}
	return time.Duration(mins) * time.Minute, nil
}

func readEnvThrowEmpty(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("enviroment variable %v not set", key)
	}
	return val, nil
}

func readConfigFile() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("error reading config file: %v", err)
		return Config{}, err
	}

	uplExpiry := time.Duration(viper.GetInt("R2_UPLOAD_EXPIRY_MINUTES")) * time.Minute
	downExpiry := time.Duration(viper.GetInt("R2_DOWNLOAD_EXPIRY_MINUTES")) * time.Minute

	return Config{
		R2_Endpoint:        viper.GetString("R2_ENDPOINT"),
		R2_Bucket:          viper.GetString("R2_BUCKET"),
		R2_Region:          viper.GetString("R2_REGION"),
		R2_Access_Key:      viper.GetString("R2_ACCESS_KEY"),
		R2_Secret_Key:      viper.GetString("R2_SECRET_KEY"),
		R2_Upload_Expiry:   uplExpiry,
		R2_Download_Expiry: downExpiry,
		Cache_Dir:          viper.GetString("CACHE_DIR"),
	}, nil
}
