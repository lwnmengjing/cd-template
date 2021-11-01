/*
 * @Author: lwnmengjing
 * @Date: 2021/11/1 11:28 上午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/11/1 11:28 上午
 */

package chart

import (
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s"
	"github.com/lwnmengjing/cd-template-go/imports/k8s"
	"github.com/lwnmengjing/cd-template-go/pkg/config"
	"log"
	"os"
	"strings"
)

func NewConfigmapChart(scope constructs.Construct, id string, props *cdk8s.ChartProps) {
	if len(config.Cfg.ConfigData) == 0 {
		return
	}

	chart := cdk8s.NewChart(scope, jsii.String(id), props)

	list := make([]*k8s.KubeConfigMapProps, len(config.Cfg.ConfigData))

	for i := range config.Cfg.ConfigData {
		data := make(map[string]*string)
		for path := range config.Cfg.ConfigData[i].Data {
			if strings.Index(config.Cfg.ConfigData[i].Data[path], string(os.PathSeparator)) > -1 &&
				strings.Index(config.Cfg.ConfigData[i].Data[path], "\n") == -1 {
				//path
				rb, err := os.ReadFile(config.Cfg.ConfigData[i].Data[path])
				if err != nil {
					log.Fatalf("read %s error, %s", config.Cfg.ConfigData[i].Data[path], err.Error())
					return
				}
				data[path] = jsii.String(string(rb))
				continue
			}
			data[path] = jsii.String(config.Cfg.ConfigData[i].Data[path])
		}
		list[i] = &k8s.KubeConfigMapProps{
			Data: &data,
			Metadata: &k8s.ObjectMeta{
				Name:      &config.Cfg.ConfigData[i].Name,
				Namespace: props.Namespace,
			},
		}
	}
	k8s.NewKubeConfigMapList(chart, jsii.String("configmap"), &k8s.KubeConfigMapListProps{
		Items: &list,
	})
}
