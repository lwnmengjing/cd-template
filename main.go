package main

import (
	"flag"
	"github.com/lwnmengjing/cd-template/pkg/config"
	"github.com/lwnmengjing/cd-template/stage"
)

var configPath = flag.String("config", "", "config path")

func main() {
	flag.Parse()
	config.NewConfig(configPath)
	stage.Synth("prod")
	config.Cfg.Hpa = false
	config.Cfg.Resources = nil
	stage.Synth("test")
}
