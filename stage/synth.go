package stage

import (
	"path/filepath"

	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s"
	"github.com/lwnmengjing/cd-template/chart"
	"github.com/lwnmengjing/cd-template/pkg/config"
)

// Synth chart generate
func Synth(stage string) {
	app := cdk8s.NewApp(&cdk8s.AppProps{Outdir: jsii.String(filepath.Join("dist", stage))})
	chart.NewServiceChart(app, config.Cfg.App+"-"+config.Cfg.Service+"-service", &cdk8s.ChartProps{
		Labels: &map[string]*string{
			"app":     &config.Cfg.Service,
			"version": &config.Cfg.Version,
		},
	})
	needConfigmap := false
	if len(config.Cfg.Config) > 0 {
		for i := range config.Cfg.Config {
			if len(config.Cfg.Config[i].Data) > 0 {
				needConfigmap = true
			}
		}
	}
	if needConfigmap {
		chart.NewConfigmapChart(app, config.Cfg.App+"-"+config.Cfg.Service+"-configmap", &cdk8s.ChartProps{
			Labels: &map[string]*string{
				"app":     &config.Cfg.Service,
				"version": &config.Cfg.Version,
			},
		})
	}
	chart.NewWorkloadChart(app, config.Cfg.App+"-"+config.Cfg.Service+"-workload", &cdk8s.ChartProps{
		Labels: &map[string]*string{
			"app":     &config.Cfg.Service,
			"version": &config.Cfg.Version,
		},
	})
	if config.Cfg.Hpa {
		chart.NewHpaChart(app, config.Cfg.App+"-"+config.Cfg.Service+"-hpa", &cdk8s.ChartProps{
			Labels: &map[string]*string{
				"app":     &config.Cfg.Service,
				"version": &config.Cfg.Version,
			},
		})
	}
	app.Synth()
}
