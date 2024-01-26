package config

import (
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Configuration viper.Viper

func LoadConfiguration() {
	configuration := *viper.New()
	configuration.SetConfigName("application")
	configuration.AddConfigPath(os.ExpandEnv("infrastructure/config"))
	configuration.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	configuration.SetConfigType("yaml")
	if err := configuration.ReadInConfig(); err != nil {
		panic(err)
	}
	configuration.WatchConfig()
	configuration.OnConfigChange(func(e fsnotify.Event) {
		//	glog.Info("App Config file changed %s:", e.Name)
	})
	Configuration = configuration
}
