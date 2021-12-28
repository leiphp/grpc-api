package config

import "github.com/spf13/viper"

func setDefault(v *viper.Viper, defaultData map[string]string) {
	for key, value := range defaultData {
		v.SetDefault(key, value)
	}
}
