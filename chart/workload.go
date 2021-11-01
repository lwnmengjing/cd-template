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
	//port
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
	//env
	env := make([]*k8s.EnvVar, len(config.Cfg.ImportEnvNames))
	for i := range config.Cfg.ImportEnvNames {
		v := os.Getenv(config.Cfg.ImportEnvNames[i])
		env[i] = &k8s.EnvVar{
			Name:  &config.Cfg.ImportEnvNames[i],
			Value: &v,
		}
	}
	//config
	volumeMounts := make([]*k8s.VolumeMount, 0)
	volumes := make([]*k8s.Volume, len(config.Cfg.ConfigData))
	if len(config.Cfg.ConfigData) > 0 {
		readOnly := true
		for i := range config.Cfg.ConfigData {
			volumes[i] = &k8s.Volume{
				Name: &config.Cfg.ConfigData[i].Name,
				ConfigMap: &k8s.ConfigMapVolumeSource{
					Name: &config.Cfg.ConfigData[i].Name,
				},
			}

			for path := range config.Cfg.ConfigData[i].Data {
				volumeMounts = append(volumeMounts, &k8s.VolumeMount{
					MountPath: &path,
					Name:      &config.Cfg.ConfigData[i].Name,
					ReadOnly:  &readOnly,
				})
			}
		}
	}
	serviceAccountName := ""
	if config.Cfg.ServiceAccount {
		serviceAccountName = config.Cfg.App + "-" + config.Cfg.Service
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
						ServiceAccountName: &serviceAccountName,
						Containers: &[]*k8s.Container{{
							Name:         jsii.String(config.Cfg.Service),
							Image:        jsii.String(config.Cfg.Image.String()),
							Ports:        &ports,
							Env:          &env,
							VolumeMounts: &volumeMounts,
						}},
						Volumes: &volumes,
					},
				},
			},
		})
	}
	return chart
}
