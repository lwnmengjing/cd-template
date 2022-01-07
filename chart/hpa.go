package chart

import (
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s"
	"github.com/lwnmengjing/cd-template-go/imports/k8s"
	"github.com/lwnmengjing/cd-template-go/pkg/config"
)

func NewHpaChart(scope constructs.Construct, id string, props *cdk8s.ChartProps) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, jsii.String(id), props)

	k8s.NewKubeHorizontalPodAutoscalerV2Beta1(chart, jsii.String("hpa"), &k8s.KubeHorizontalPodAutoscalerV2Beta1Props{
		Metadata: &k8s.ObjectMeta{
			Labels: props.Labels,
			Name:   &config.Cfg.Service,
		},
		Spec: &k8s.HorizontalPodAutoscalerSpecV2Beta1{
			MinReplicas: jsii.Number(float64(config.Cfg.Replicas)),
			MaxReplicas: jsii.Number(float64(config.Cfg.MaxReplicas)),
			ScaleTargetRef: &k8s.CrossVersionObjectReferenceV2Beta1{
				Kind:       jsii.String(config.Cfg.WorkloadType),
				Name:       jsii.String(config.Cfg.Service),
				ApiVersion: jsii.String("apps/v1"),
			},
			Metrics: &[]*k8s.MetricSpecV2Beta1{
				{
					Type: jsii.String("Resource"),
					Resource: &k8s.ResourceMetricSourceV2Beta1{
						Name:                     jsii.String("cpu"),
						TargetAverageUtilization: jsii.Number(80),
					},
				},
				{
					Type: jsii.String("Resource"),
					Resource: &k8s.ResourceMetricSourceV2Beta1{
						Name:                     jsii.String("memory"),
						TargetAverageUtilization: jsii.Number(80),
					},
				},
			},
		},
	})

	return chart
}
