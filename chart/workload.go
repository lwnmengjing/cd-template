/*
 * @Author: lwnmengjing
 * @Date: 2021/10/29 11:21 下午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/10/29 11:21 下午
 */

package chart

import (
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s"
	"github.com/lwnmengjing/cd-template-go/imports/k8s"
	"github.com/lwnmengjing/cd-template-go/pkg/config"
	"os"
)

func NewWorkloadChart(scope constructs.Construct, id string, props *cdk8s.ChartProps) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, jsii.String(id), props)
	workloadName := config.Cfg.Service + "-" + config.Cfg.Version
	ports := make([]*k8s.ContainerPort, 0)
	var samePort bool
	for i := range config.Cfg.Ports {
		samePort = !config.Cfg.Metrics.Scrape || config.Cfg.Ports[i].Port == config.Cfg.Metrics.Port
		ports = append(ports, &k8s.ContainerPort{
			ContainerPort: jsii.Number(float64(config.Cfg.Ports[i].Port)),
		})
	}
	if !samePort {
		ports = append(ports, &k8s.ContainerPort{
			ContainerPort: jsii.Number(float64(config.Cfg.Metrics.Port)),
		})
	}
	env := make([]*k8s.EnvVar, len(config.Cfg.ImportEnvNames))
	for i := range config.Cfg.ImportEnvNames {
		v := os.Getenv(config.Cfg.ImportEnvNames[i])
		env[i] = &k8s.EnvVar{
			Name:  &config.Cfg.ImportEnvNames[i],
			Value: &v,
		}
	}
	switch config.Cfg.WorkloadType {
	default:
		k8s.NewKubeDeployment(chart, jsii.String("deployment"), &k8s.KubeDeploymentProps{
			Metadata: &k8s.ObjectMeta{
				Name:      &workloadName,
				Namespace: &config.Cfg.Namespace,
				Labels:    props.Labels,
			},
			Spec: &k8s.DeploymentSpec{
				Replicas: jsii.Number(float64(config.Cfg.Replicas)),
				Selector: &k8s.LabelSelector{
					MatchLabels: props.Labels,
				},
				Template: &k8s.PodTemplateSpec{
					Metadata: &k8s.ObjectMeta{
						Labels: props.Labels,
					},
					Spec: &k8s.PodSpec{
						Containers: &[]*k8s.Container{{
							Name:  jsii.String(config.Cfg.Service),
							Image: jsii.String(config.Cfg.Image.String()),
							Ports: &ports,
							Env:   &env,
						}},
					},
				},
			},
		})
	}
	return chart
}
