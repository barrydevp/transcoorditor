package common

import (
	// "fmt"
	"github.com/spf13/viper"
)

func loadDefaults() {
	// for API Server
	viper.SetDefault("PORT", "8000")
	viper.SetDefault("SERVER_READ_TIMEOUT", "300")

	// mongodb
	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")
	viper.SetDefault("MONGODB_DB", "transcoorditor")

}

func InitEnv(envFile string) error {
	loadDefaults()

	viper.SetConfigFile(envFile)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// fmt.Println("Error while reading env from", envFile, err)

		return err
	}

	// Override config parameters from environment variables if specified
	// for _, key := range viper.AllKeys() {
	// 	viper.BindEnv(key)
	// }

	return nil
}
