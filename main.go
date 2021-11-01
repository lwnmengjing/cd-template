package main

import (
	"flag"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s"

	"github.com/lwnmengjing/cd-template-go/chart"
	"github.com/lwnmengjing/cd-template-go/pkg/config"
)

var configPath = flag.String("config", "config.local.yaml", "config path")

func main() {
	config.NewConfig(*configPath)
	app := cdk8s.NewApp(nil)
	chart.NewServiceChart(app, config.Cfg.App+"-"+config.Cfg.Service+"-service", &cdk8s.ChartProps{
		Labels: &map[string]*string{
			"app":     &config.Cfg.Service,
			"version": &config.Cfg.Version,
		},
		Namespace: &config.Cfg.Namespace,
	})
	chart.NewWorkloadChart(app, config.Cfg.App+"-"+config.Cfg.Service+"-workload", &cdk8s.ChartProps{
		Labels: &map[string]*string{
			"app":     &config.Cfg.Service,
			"version": &config.Cfg.Version,
		},
		Namespace: &config.Cfg.Namespace,
	})
	chart.NewConfigmapChart(app, config.Cfg.App+"-"+config.Cfg.Service+"-configmap", &cdk8s.ChartProps{
		Labels: &map[string]*string{
			"app":     &config.Cfg.Service,
			"version": &config.Cfg.Version,
		},
		Namespace: &config.Cfg.Namespace,
	})
	app.Synth()
}
