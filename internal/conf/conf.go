package conf

import (
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util"
	"github.com/spf13/viper"
)

var conf *config

func Init() {
	listDir := []string{".", "../", "../../", "../../../", "../../../../"}

	for _, dir := range listDir {
		viper.SetConfigName("env")
		viper.SetConfigType("json")
		viper.AddConfigPath(dir)
		err := viper.ReadInConfig()
		if err == nil {
			viper.SetConfigName("env.override")
			err = viper.MergeInConfig()
			util.Panic(err)

			if err = viper.Unmarshal(&conf); err != nil {
				panic(err)
			}

			newLogger()
			return
		}
	}

	panic("cannot load env")
}

func GetConfig() *config {
	return conf
}
