package util

import (
	"log"

	"github.com/spf13/viper"
)

// Creating based on .env
type Config struct {
	DBDriver      string `mapstructure:"DB_Driver"`
	DBSource      string `mapstructure:"DB_Source"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	TokenAPI      string `mapstructure:"TOKEN_SYMETRIC_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("credential")
	viper.SetConfigType("env") // it can json,xml,ini, and etc

	// overrride values that read from config file
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// Log the error but continue, expecting environment variables to be used
		log.Printf("Warning: Could not read config file: %v\n", err)
	}

	err = viper.Unmarshal(&config)
	return
}
