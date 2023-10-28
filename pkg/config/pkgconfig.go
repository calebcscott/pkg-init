package config

import (
	"github.com/spf13/viper"
)


type PkgConfig struct {
    CacheDir string
    TemplateMap map[string]string
}


func NewConfig() PkgConfig {
    cache := viper.GetString("data-dir")
    templateMap := viper.GetStringMapString("templates")

    return PkgConfig{
        CacheDir: cache,
        TemplateMap: templateMap,
    }
}


