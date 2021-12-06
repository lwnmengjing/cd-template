/*
 * @Author: lwnmengjing
 * @Date: 2021/11/1 11:28 上午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/11/1 11:28 上午
 */

package chart

import (
	"log"
	"os"
	"strings"

	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s"
	"github.com/lwnmengjing/cd-template-go/imports/k8s"
	"github.com/lwnmengjing/cd-template-go/pkg/config"
)

func NewConfigmapChart(scope constructs.Construct, id string, props *cdk8s.ChartProps) cdk8s.Chart {
	if len(config.Cfg.ConfigData.Data) == 0 {
		return nil
	}

	chart := cdk8s.NewChart(scope, jsii.String(id), props)

	data := make(map[string]*string)
	for path := range config.Cfg.ConfigData.Data {
		if strings.Index(config.Cfg.ConfigData.Data[path], string(os.PathSeparator)) > -1 &&
			strings.Index(config.Cfg.ConfigData.Data[path], "\n") == -1 {
			//path
			rb, err := os.ReadFile(config.Cfg.ConfigData.Data[path])
			if err != nil {
				log.Fatalf("read %s error, %s", config.Cfg.ConfigData.Data[path], err.Error())
				return nil
			}
			data[path] = jsii.String(string(rb))
			continue
		}
		data[path] = jsii.String(config.Cfg.ConfigData.Data[path])
	}
	cm := k8s.NewKubeConfigMap(chart, jsii.String("configmap"), &k8s.KubeConfigMapProps{
		Data: &data,
		Metadata: &k8s.ObjectMeta{
			//Name:      &config.Cfg.ConfigData[i].Name,
			Namespace: props.Namespace,
		},
	})
	config.Cfg.ConfigData.Name = *cm.Name()
	return chart
}
