/*
 * @Author: lwnmengjing
 * @Date: 2021/10/29 10:30 下午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/10/29 10:30 下午
 */

package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
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
	WorkloadType   string          `json:"workloadType"`
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

// NewConfig set config
func NewConfig(path string) {
	rb, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(rb, &Cfg)
	if err != nil {
		log.Fatalln(err)
	}
}
