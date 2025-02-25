// Copyright 2024 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	v1alpha1 "github.com/kubeflow/trainer/pkg/apis/trainer/v1alpha1"
	internal "github.com/kubeflow/trainer/pkg/client/applyconfiguration/internal"
	trainerv1alpha1 "github.com/kubeflow/trainer/pkg/client/applyconfiguration/trainer/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=trainer.kubeflow.org, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("ClusterTrainingRuntime"):
		return &trainerv1alpha1.ClusterTrainingRuntimeApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ContainerOverride"):
		return &trainerv1alpha1.ContainerOverrideApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("CoschedulingPodGroupPolicySource"):
		return &trainerv1alpha1.CoschedulingPodGroupPolicySourceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("DatasetConfig"):
		return &trainerv1alpha1.DatasetConfigApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("InputModel"):
		return &trainerv1alpha1.InputModelApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("JobSetTemplateSpec"):
		return &trainerv1alpha1.JobSetTemplateSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("JobStatus"):
		return &trainerv1alpha1.JobStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("MLPolicy"):
		return &trainerv1alpha1.MLPolicyApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("MLPolicySource"):
		return &trainerv1alpha1.MLPolicySourceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ModelConfig"):
		return &trainerv1alpha1.ModelConfigApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("MPIMLPolicySource"):
		return &trainerv1alpha1.MPIMLPolicySourceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("OutputModel"):
		return &trainerv1alpha1.OutputModelApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PodGroupPolicy"):
		return &trainerv1alpha1.PodGroupPolicyApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PodGroupPolicySource"):
		return &trainerv1alpha1.PodGroupPolicySourceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PodSpecOverride"):
		return &trainerv1alpha1.PodSpecOverrideApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("PodSpecOverrideTargetJob"):
		return &trainerv1alpha1.PodSpecOverrideTargetJobApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("RuntimeRef"):
		return &trainerv1alpha1.RuntimeRefApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("TorchElasticPolicy"):
		return &trainerv1alpha1.TorchElasticPolicyApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("TorchMLPolicySource"):
		return &trainerv1alpha1.TorchMLPolicySourceApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Trainer"):
		return &trainerv1alpha1.TrainerApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("TrainingRuntime"):
		return &trainerv1alpha1.TrainingRuntimeApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("TrainingRuntimeSpec"):
		return &trainerv1alpha1.TrainingRuntimeSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("TrainJob"):
		return &trainerv1alpha1.TrainJobApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("TrainJobSpec"):
		return &trainerv1alpha1.TrainJobSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("TrainJobStatus"):
		return &trainerv1alpha1.TrainJobStatusApplyConfiguration{}

	}
	return nil
}

func NewTypeConverter(scheme *runtime.Scheme) *testing.TypeConverter {
	return &testing.TypeConverter{Scheme: scheme, TypeResolver: internal.Parser()}
}
