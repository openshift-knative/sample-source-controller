/*
Copyright 2019 The Knative Authors.

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

package v1alpha1


import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/eventing/pkg/apis/duck"
	"knative.dev/pkg/apis"
)

const (
	// SampleConditionReady has status True when the SampleSource is ready to send events.
	SampleConditionReady = apis.ConditionReady

	// SampleConditionSinkProvided has status True when the SampleSource has been configured with a sink target.
	SampleConditionSinkProvided apis.ConditionType = "SinkProvided"

	// SampleConditionDeployed has status True when the SampleSource has had it's deployment created.
	SampleConditionDeployed apis.ConditionType = "Deployed"

	// SampleConditionEventTypeProvided has status True when the SampleSource has been configured with its event types.
	SampleConditionEventTypeProvided apis.ConditionType = "EventTypesProvided"
)

var SampleCondSet = apis.NewLivingConditionSet(
	SampleConditionSinkProvided,
	SampleConditionDeployed,
	SampleConditionEventTypeProvided,
)

// GetCondition returns the condition currently associated with the given type, or nil.
func (s *SampleSourceStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return SampleCondSet.Manage(s).GetCondition(t)
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (s *SampleSourceStatus) InitializeConditions() {
	SampleCondSet.Manage(s).InitializeConditions()
}

// MarkSinkWarnDeprecated sets the condition that the source has a sink configured and warns ref is deprecated.
func (s *SampleSourceStatus) MarkSinkWarnRefDeprecated(uri string) {
	s.SinkURI = uri
	if len(uri) > 0 {
		c := apis.Condition{
			Type:     SampleConditionSinkProvided,
			Status:   corev1.ConditionTrue,
			Severity: apis.ConditionSeverityError,
			Message:  "Using deprecated object ref fields when specifying spec.sink. These will be removed in a future release. Update to spec.sink.ref.",
		}
		SampleCondSet.Manage(s).SetCondition(c)
	} else {
		SampleCondSet.Manage(s).MarkUnknown(SampleConditionSinkProvided, "SinkEmpty", "Sink has resolved to empty.%s", "")
	}
}

// MarkSink sets the condition that the source has a sink configured.
func (s *SampleSourceStatus) MarkSink(uri string) {
	s.SinkURI = uri
	if len(uri) > 0 {
		SampleCondSet.Manage(s).MarkTrue(SampleConditionSinkProvided)
	} else {
		SampleCondSet.Manage(s).MarkUnknown(SampleConditionSinkProvided, "SinkEmpty", "Sink has resolved to empty.%s", "")
	}
}

// MarkNoSink sets the condition that the source does not have a sink configured.
func (s *SampleSourceStatus) MarkNoSink(reason, messageFormat string, messageA ...interface{}) {
	SampleCondSet.Manage(s).MarkFalse(SampleConditionSinkProvided, reason, messageFormat, messageA...)
}

// PropagateDeploymentAvailability uses the availability of the provided Deployment to determine if
// SampleConditionDeployed should be marked as true or false.
func (s *SampleSourceStatus) PropagateDeploymentAvailability(d *appsv1.Deployment) {
	if duck.DeploymentIsAvailable(&d.Status, false) {
		SampleCondSet.Manage(s).MarkTrue(SampleConditionDeployed)
	} else {
		// I don't know how to propagate the status well, so just give the name of the Deployment
		// for now.
		SampleCondSet.Manage(s).MarkFalse(SampleConditionDeployed, "DeploymentUnavailable", "The Deployment '%s' is unavailable.", d.Name)
	}
}

// MarkEventTypes sets the condition that the source has set its event type.
func (s *SampleSourceStatus) MarkEventTypes() {
	SampleCondSet.Manage(s).MarkTrue(SampleConditionEventTypeProvided)
}

// MarkNoEventTypes sets the condition that the source does not its event type configured.
func (s *SampleSourceStatus) MarkNoEventTypes(reason, messageFormat string, messageA ...interface{}) {
	SampleCondSet.Manage(s).MarkFalse(SampleConditionEventTypeProvided, reason, messageFormat, messageA...)
}

// IsReady returns true if the resource is ready overall.
func (s *SampleSourceStatus) IsReady() bool {
	return SampleCondSet.Manage(s).IsHappy()
}
