/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resources

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/eventing/pkg/utils"
	"knative.dev/pkg/kmeta"

	"knative.dev/sample-source/pkg/apis/samples/v1alpha1"
)

// ReceiveAdapterArgs are the arguments needed to create a Sample Source Receive Adapter.
// Every field is required.
type ReceiveAdapterArgs struct {
	EventSource string
	Image       string
	Source      *v1alpha1.SampleSource
	Labels      map[string]string
	SinkURI     string
}

// MakeReceiveAdapter generates (but does not insert into K8s) the Receive Adapter Deployment for
// Sample sources.
func MakeReceiveAdapter(args *ReceiveAdapterArgs) *v1.Deployment {
	replicas := int32(1)
	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: args.Source.Namespace,
			Name:      utils.GenerateFixedName(args.Source, fmt.Sprintf("samplesource-%s", args.Source.Name)),
			Labels:    args.Labels,
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(args.Source),
			},
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: args.Labels,
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: args.Labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: args.Source.Spec.ServiceAccountName,
					Containers: []corev1.Container{
						{
							Name:  "receive-adapter",
							Image: args.Image,
							Env:   makeEnv(args.EventSource, args.SinkURI, &args.Source.Spec),
						},
					},
				},
			},
		},
	}
}

func makeEnv(eventSource, sinkURI string, spec *v1alpha1.SampleSourceSpec) []corev1.EnvVar {
	return []corev1.EnvVar{{
		Name:  "SINK_URI",
		Value: sinkURI,
	}, {
		Name:  "EVENT_SOURCE",
		Value: eventSource,
	}, {
		Name:  "INTERVAL",
		Value: spec.Interval,
	}, {
		Name: "NAMESPACE",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	}, {
		Name:  "METRICS_DOMAIN",
		Value: "knative.dev/eventing",
	}, {
		Name:  "K_METRICS_CONFIG",
		Value: "",
	}, {
		Name:  "K_LOGGING_CONFIG",
		Value: "",
	}}
}
