package util

import "github.com/spf13/viper"

type Config struct {
	DBDriver      string `mapstructure:"dbdriver"`
	DBSource      string `mapstructure:"dbsource"`
	ServerAddress string `mapstructure:"serveraddress"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigFile("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return

}
