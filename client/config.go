package client

import (
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.SetConfigName("config")        // name of config file (without extension)
	viper.AddConfigPath("/etc/golook/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.golook") // call multiple times to add many search paths
	viper.AddConfigPath(".")             // call multiple times to add many search paths
	err := viper.ReadInConfig()          // Find and read the config file
	if err != nil {                      // Handle errors reading the config file
		log.Fatalf("Fatal error config file: %s \n", err)
	}
}
