package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)


type PkgConfig struct {
    CacheDir string
    TemplateMap map[string]string
    LanguageSupport map[string]string
}


func NewConfig() PkgConfig {
    var err error


    cache := viper.GetString("data-dir")
    if cache != "" {
        home := os.Getenv("HOME")
        replacer := strings.NewReplacer("~", home, "${home}", home)
        cache = replacer.Replace(cache)

    } else {
        cache, err = os.UserCacheDir()
        if err != nil {
            fmt.Println("Could not find suitable cache, defaulting to ~/.cache")
            cache = "~/.cache"
        }

        cache = filepath.Join(cache, "pkg-init")
    }
    if _, err := os.Stat(cache); err != nil {
        fmt.Println("Attempting to ensure cache dir exists: ", cache)
        os.MkdirAll(cache, 0755)
    }
    templateMap := viper.GetStringMapString("templates")
    langs := viper.GetStringMapString("languages")

    return PkgConfig{
        CacheDir: cache,
        TemplateMap: templateMap,
        LanguageSupport: langs,
    }
}


