package main

import (
	"fmt"

	"github.com/20326/flexbox/config"
	"github.com/20326/flexbox/logger"
)

var (
	Version = "v1.0.0"
)

type MyConfig struct {
	Logger loggerSection `yaml:"logger"`
}

type loggerSection struct {
	Level      string `yaml:"level"`
	FileName   string `yaml:"fileName"`
	MaxSize    int64  `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	Cron       string `yaml:"cron"`
}

func main() {

	var cfg MyConfig

	if errs := config.New("").LoadFile("logger.yaml", &cfg).End(); len(errs) > 0 {
		panic("not found config file")
	}

	fmt.Printf("cfg: %v\n", cfg)

	log := logger.NewMultiLogger(
		logger.WithLevel(cfg.Logger.Level),
		logger.WithFileName(cfg.Logger.FileName),
		logger.WithCronRunner(cfg.Logger.Cron),
		logger.WithMaxSize(cfg.Logger.MaxSize),
		logger.WithMaxBackups(cfg.Logger.MaxBackups),
		logger.WithCronRunner(cfg.Logger.Cron),
	)

	log.Info().Str("version", Version).Msg("start")
	log.Info().Msg("init logger ok")

	fmt.Println("finished")
}
