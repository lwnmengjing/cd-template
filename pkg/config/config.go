/*
 * @Author: lwnmengjing
 * @Date: 2021/10/29 10:30 下午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/10/29 10:30 下午
 */

package config

import (
	"flag"
	"log"
	"strings"

	"github.com/lwnmengjing/core-go/config"
)

var Cfg Config

type Config struct {
	Namespace      string          `json:"namespace" yaml:"namespace"`
	App            string          `json:"app" yaml:"app"`
	Service        string          `json:"service" yaml:"service"`
	Version        string          `json:"version" yaml:"version"`
	Replicas       uint            `json:"replicas" yaml:"replicas"`
	ServiceAccount bool            `json:"serviceAccount" yaml:"serviceAccount"`
	Image          Image           `json:"image" yaml:"image"`
	Ports          []Port          `json:"ports" yaml:"ports"`
	Metrics        Metrics         `json:"metrics" yaml:"metrics"`
	ImportEnvNames []string        `json:"importEnvNames" yaml:"importEnvNames"`
	ConfigData     []ConfigmapData `json:"configData" yaml:"configData"`
	WorkloadType   string          `json:"workloadType" yaml:"workloadType"`
	Command        []*string       `json:"command" yaml:"command"`
	Args           []*string       `json:"args" yaml:"args"`
}

type Port struct {
	Port       uint   `json:"port" yaml:"port"`
	TargetPort uint   `json:"targetPort" yaml:"targetPort"`
	Name       string `json:"name" yaml:"name"`
}

type Image struct {
	Path   string `json:"path" yaml:"path"`
	Tag    string `json:"tag" yaml:"tag"`
	Secret string `json:"secret" yaml:"secret"`
}

func (e Image) String() string {
	if e.Tag == "" {
		return e.Path
	}
	return e.Path + ":" + e.Tag
}

type Metrics struct {
	Scrape bool   `json:"scrape" yaml:"scrape"`
	Path   string `json:"path" yaml:"path"`
	Port   uint   `json:"port" yaml:"port"`
}

type ConfigmapData struct {
	Name string            `json:"name" yaml:"name"`
	Path string            `json:"path" yaml:"path"`
	Data map[string]string `json:"data" yaml:"data"`
}

var (
	namespace      = flag.String("namespace", "default", "deploy namespace")
	app            = flag.String("app", "", "application")
	service        = flag.String("service", "", "service")
	version        = flag.String("version", "v1", "service version")
	port           = flag.Uint("port", 8000, "port")
	portName       = flag.String("portName", "http", "port name")
	image          = flag.String("image", "", "image:tag")
	importEnvNames = flag.String("importEnvNames", "", "import env names, split ','")
)

// NewConfig set config
func NewConfig(path *string) {
	var err error
	config.DefaultConfig, err = config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	if path != nil && *path != "" {
		err = config.LoadFile(*path)
		if err != nil {
			log.Fatalln(err)
		}
		err = config.Scan(&Cfg)
		if err != nil {
			log.Fatalln(err)
		}
	}

	Cfg.Namespace = config.Get("namespace").String(*namespace)
	Cfg.App = config.Get("app").String(*app)
	Cfg.Service = config.Get("service").String(*service)
	if len(Cfg.Ports) == 0 {
		Cfg.Ports = []Port{
			{
				Port:       *port,
				Name:       *portName,
				TargetPort: *port,
			},
		}
	}
	Cfg.Image.Path = config.Get("image", "path").String(*image)
	Cfg.Version = config.Get("version").String(*version)
	if len(Cfg.ImportEnvNames) == 0 {
		Cfg.ImportEnvNames = strings.Split(*importEnvNames, ",")
	}
}
