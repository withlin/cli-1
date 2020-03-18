/*
Copyright 2020 The Tekton Authors

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

// nolint: golint
package v1alpha1

import (
	"context"
	"fmt"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"knative.dev/pkg/apis"
)

var _ apis.Convertible = (*Task)(nil)

// ConvertUp implements api.Convertible
func (source *Task) ConvertUp(ctx context.Context, obj apis.Convertible) error {
	switch sink := obj.(type) {
	case *v1beta1.Task:
		sink.ObjectMeta = source.ObjectMeta
		return source.Spec.ConvertUp(ctx, &sink.Spec)
	default:
		return fmt.Errorf("unknown version, got: %T", sink)
	}
}

func (source *TaskSpec) ConvertUp(ctx context.Context, sink *v1beta1.TaskSpec) error {
	sink.Steps = source.Steps
	sink.Volumes = source.Volumes
	sink.StepTemplate = source.StepTemplate
	sink.Sidecars = source.Sidecars
	sink.Workspaces = source.Workspaces
	sink.Results = source.Results
	sink.Resources = source.Resources
	sink.Params = source.Params
	if source.Inputs != nil {
		if len(source.Inputs.Params) > 0 && len(source.Params) > 0 {
			// This shouldn't happen as it shouldn't pass validation
			return apis.ErrMultipleOneOf("inputs.params", "params")
		}
		if len(source.Inputs.Params) > 0 {
			sink.Params = make([]v1beta1.ParamSpec, len(source.Inputs.Params))
			for i, param := range source.Inputs.Params {
				sink.Params[i] = *param.DeepCopy()
			}
		}
		if len(source.Inputs.Resources) > 0 {
			if sink.Resources == nil {
				sink.Resources = &v1beta1.TaskResources{}
			}
			if len(source.Inputs.Resources) > 0 && source.Resources != nil && len(source.Resources.Inputs) > 0 {
				// This shouldn't happen as it shouldn't pass validation but just in case
				return apis.ErrMultipleOneOf("inputs.resources", "resources.inputs")
			}
			sink.Resources.Inputs = make([]v1beta1.TaskResource, len(source.Inputs.Resources))
			for i, resource := range source.Inputs.Resources {
				sink.Resources.Inputs[i] = v1beta1.TaskResource{ResourceDeclaration: v1beta1.ResourceDeclaration{
					Name:        resource.Name,
					Type:        resource.Type,
					Description: resource.Description,
					TargetPath:  resource.TargetPath,
					Optional:    resource.Optional,
				}}
			}
		}
	}
	if source.Outputs != nil && len(source.Outputs.Resources) > 0 {
		if sink.Resources == nil {
			sink.Resources = &v1beta1.TaskResources{}
		}
		if len(source.Outputs.Resources) > 0 && source.Resources != nil && len(source.Resources.Outputs) > 0 {
			// This shouldn't happen as it shouldn't pass validation but just in case
			return apis.ErrMultipleOneOf("outputs.resources", "resources.outputs")
		}
		sink.Resources.Outputs = make([]v1beta1.TaskResource, len(source.Outputs.Resources))
		for i, resource := range source.Outputs.Resources {
			sink.Resources.Outputs[i] = v1beta1.TaskResource{ResourceDeclaration: v1beta1.ResourceDeclaration{
				Name:        resource.Name,
				Type:        resource.Type,
				Description: resource.Description,
				TargetPath:  resource.TargetPath,
				Optional:    resource.Optional,
			}}
		}
	}
	return nil
}

// ConvertDown implements api.Convertible
func (sink *Task) ConvertDown(ctx context.Context, obj apis.Convertible) error {
	switch source := obj.(type) {
	case *v1beta1.Task:
		sink.ObjectMeta = source.ObjectMeta
		return sink.Spec.ConvertDown(ctx, &source.Spec)
	default:
		return fmt.Errorf("unknown version, got: %T", sink)
	}
}

func (sink *TaskSpec) ConvertDown(ctx context.Context, source *v1beta1.TaskSpec) error {
	sink.Steps = source.Steps
	sink.Volumes = source.Volumes
	sink.StepTemplate = source.StepTemplate
	sink.Sidecars = source.Sidecars
	sink.Workspaces = source.Workspaces
	sink.Results = source.Results
	sink.Params = source.Params
	sink.Resources = source.Resources
	return nil
}