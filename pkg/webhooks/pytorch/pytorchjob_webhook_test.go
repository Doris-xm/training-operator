/*
Copyright 2024 The Kubeflow Authors.

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

package pytorch

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	corev1 "k8s.io/api/core/v1"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	trainingoperator "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
)

func TestValidateV1PyTorchJob(t *testing.T) {
	validPyTorchReplicaSpecs := map[trainingoperator.ReplicaType]*trainingoperator.ReplicaSpec{
		trainingoperator.PyTorchJobReplicaTypeMaster: {
			Replicas:      ptr.To[int32](1),
			RestartPolicy: trainingoperator.RestartPolicyOnFailure,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            "pytorch",
						Image:           "docker.io/kubeflowkatib/pytorch-mnist:v1beta1-45c5727",
						ImagePullPolicy: corev1.PullAlways,
						Command: []string{
							"python3",
							"/opt/pytorch-mnist/mnist.py",
							"--epochs=1",
						},
					}},
				},
			},
		},
		trainingoperator.PyTorchJobReplicaTypeWorker: {
			Replicas:      ptr.To[int32](1),
			RestartPolicy: trainingoperator.RestartPolicyOnFailure,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            "pytorch",
						Image:           "docker.io/kubeflowkatib/pytorch-mnist:v1beta1-45c5727",
						ImagePullPolicy: corev1.PullAlways,
						Command: []string{
							"python3",
							"/opt/pytorch-mnist/mnist.py",
							"--epochs=1",
						},
					}},
				},
			},
		},
	}

	testCases := map[string]struct {
		pytorchJob    *trainingoperator.PyTorchJob
		oldPytorchJob *trainingoperator.PyTorchJob
		wantErr       field.ErrorList
		wantWarnings  admission.Warnings
	}{
		"valid PyTorchJob": {
			pytorchJob: &trainingoperator.PyTorchJob{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: trainingoperator.PyTorchJobSpec{
					RunPolicy: trainingoperator.RunPolicy{
						ManagedBy: ptr.To(trainingoperator.KubeflowJobsController),
					},
					ElasticPolicy: &trainingoperator.ElasticPolicy{
						RDZVBackend: ptr.To(trainingoperator.BackendC10D),
					},
					PyTorchReplicaSpecs: validPyTorchReplicaSpecs,
				},
			},
		},
		"pytorchJob name does not meet DNS1035": {
			pytorchJob: &trainingoperator.PyTorchJob{
				ObjectMeta: metav1.ObjectMeta{
					Name: "0-test",
				},
				Spec: trainingoperator.PyTorchJobSpec{
					PyTorchReplicaSpecs: validPyTorchReplicaSpecs,
				},
			},
			wantErr: field.ErrorList{
				field.Invalid(field.NewPath("metadata").Child("name"), "", ""),
			},
		},
		"no containers": {
			pytorchJob: &trainingoperator.PyTorchJob{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: trainingoperator.PyTorchJobSpec{
					PyTorchReplicaSpecs: map[trainingoperator.ReplicaType]*trainingoperator.ReplicaSpec{
						trainingoperator.PyTorchJobReplicaTypeWorker: {
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{},
								},
							},
						},
					},
				},
			},
			wantErr: field.ErrorList{
				field.Required(pytorchReplicaSpecPath.
					Key(string(trainingoperator.PyTorchJobReplicaTypeWorker)).
					Child("template").
					Child("spec").
					Child("containers"), ""),
				field.Required(pytorchReplicaSpecPath.
					Key(string(trainingoperator.PyTorchJobReplicaTypeWorker)).
					Child("template").
					Child("spec").
					Child("containers"), ""),
			},
		},
		"image is empty": {
			pytorchJob: &trainingoperator.PyTorchJob{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: trainingoperator.PyTorchJobSpec{
					PyTorchReplicaSpecs: map[trainingoperator.ReplicaType]*trainingoperator.ReplicaSpec{
						trainingoperator.PyTorchJobReplicaTypeWorker: {
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  "pytorch",
											Image: "",
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: field.ErrorList{
				field.Required(pytorchReplicaSpecPath.
					Key(string(trainingoperator.PyTorchJobReplicaTypeWorker)).
					Child("template").
					Child("spec").
					Child("containers").
					Index(0).
					Child("image"), ""),
			},
		},
		"pytorchJob default container name doesn't present": {
			pytorchJob: &trainingoperator.PyTorchJob{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: trainingoperator.PyTorchJobSpec{
					PyTorchReplicaSpecs: map[trainingoperator.ReplicaType]*trainingoperator.ReplicaSpec{
						trainingoperator.PyTorchJobReplicaTypeWorker: {
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  "",
											Image: "gcr.io/kubeflow-ci/pytorch-dist-mnist_test:1.0",
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: field.ErrorList{
				field.Required(pytorchReplicaSpecPath.
					Key(string(trainingoperator.PyTorchJobReplicaTypeWorker)).
					Child("template").
					Child("spec").
					Child("containers"), ""),
			},
		},
		"the number of replicas in masterReplica is other than 1": {
			pytorchJob: &trainingoperator.PyTorchJob{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: trainingoperator.PyTorchJobSpec{
					PyTorchReplicaSpecs: map[trainingoperator.ReplicaType]*trainingoperator.ReplicaSpec{
						trainingoperator.PyTorchJobReplicaTypeMaster: {
							Replicas: ptr.To[int32](2),
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  "pytorch",
											Image: "gcr.io/kubeflow-ci/pytorch-dist-mnist_test:1.0",
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: field.ErrorList{
				field.Forbidden(pytorchReplicaSpecPath.Key(string(trainingoperator.PyTorchJobReplicaTypeMaster)).Child("replicas"), ""),
			},
		},
		"Spec.ElasticPolicy.NProcPerNode are set": {
			pytorchJob: &trainingoperator.PyTorchJob{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: trainingoperator.PyTorchJobSpec{
					ElasticPolicy: &trainingoperator.ElasticPolicy{
						NProcPerNode: ptr.To[int32](1),
					},
					PyTorchReplicaSpecs: map[trainingoperator.ReplicaType]*trainingoperator.ReplicaSpec{
						trainingoperator.PyTorchJobReplicaTypeMaster: {
							Replicas: ptr.To[int32](1),
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  "pytorch",
											Image: "gcr.io/kubeflow-ci/pytorch-dist-mnist_test:1.0",
										},
									},
								},
							},
						},
						trainingoperator.PyTorchJobReplicaTypeWorker: {
							Replicas: ptr.To[int32](1),
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  "pytorch",
											Image: "gcr.io/kubeflow-ci/pytorch-dist-mnist_test:1.0",
										},
									},
								},
							},
						},
					},
				},
			},
			wantWarnings: admission.Warnings{
				fmt.Sprintf("%s is deprecated, use %s instead",
					specPath.Child("elasticPolicy").Child("nProcPerNode"), specPath.Child("nprocPerNode")),
			},
		},
		"Spec.NprocPerNode and Spec.ElasticPolicy.NProcPerNode are set": {
			pytorchJob: &trainingoperator.PyTorchJob{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: trainingoperator.PyTorchJobSpec{
					NprocPerNode: ptr.To("1"),
					ElasticPolicy: &trainingoperator.ElasticPolicy{
						NProcPerNode: ptr.To[int32](1),
					},
					PyTorchReplicaSpecs: map[trainingoperator.ReplicaType]*trainingoperator.ReplicaSpec{
						trainingoperator.PyTorchJobReplicaTypeMaster: {
							Replicas: ptr.To[int32](1),
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  "pytorch",
											Image: "gcr.io/kubeflow-ci/pytorch-dist-mnist_test:1.0",
										},
									},
								},
							},
						},
						trainingoperator.PyTorchJobReplicaTypeWorker: {
							Replicas: ptr.To[int32](1),
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  "pytorch",
											Image: "gcr.io/kubeflow-ci/pytorch-dist-mnist_test:1.0",
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: field.ErrorList{
				field.Forbidden(specPath.Child("elasticPolicy").Child("nProcPerNode"), ""),
			},
			wantWarnings: admission.Warnings{
				fmt.Sprintf("%s is deprecated, use %s instead",
					specPath.Child("elasticPolicy").Child("nProcPerNode"), specPath.Child("nprocPerNode")),
			},
		},
		"attempt to set unsupported managedBy controller name gets rejected": {
			pytorchJob: &trainingoperator.PyTorchJob{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: trainingoperator.PyTorchJobSpec{
					RunPolicy: trainingoperator.RunPolicy{
						ManagedBy: ptr.To("other-job-controller"),
					},
					PyTorchReplicaSpecs: validPyTorchReplicaSpecs,
				},
			},
			wantErr: field.ErrorList{
				field.NotSupported(field.NewPath("spec", "runPolicy", "managedBy"), "", sets.List(sets.New(
					trainingoperator.MultiKueueController,
					trainingoperator.KubeflowJobsController))),
			},
		},
		"attempt to update the managedBy field gets rejected": {
			oldPytorchJob: &trainingoperator.PyTorchJob{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: trainingoperator.PyTorchJobSpec{
					RunPolicy: trainingoperator.RunPolicy{
						ManagedBy: ptr.To(trainingoperator.KubeflowJobsController),
					},
					PyTorchReplicaSpecs: validPyTorchReplicaSpecs,
				},
			},
			pytorchJob: &trainingoperator.PyTorchJob{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: trainingoperator.PyTorchJobSpec{
					RunPolicy: trainingoperator.RunPolicy{
						ManagedBy: ptr.To(trainingoperator.MultiKueueController),
					},
					PyTorchReplicaSpecs: validPyTorchReplicaSpecs,
				},
			},
			wantErr: field.ErrorList{
				field.Invalid(field.NewPath("spec", "runPolicy", "managedBy"), trainingoperator.MultiKueueController, apivalidation.FieldImmutableErrorMsg),
			},
		},
		"attempt to configure elasticPolicy when no worker is configured": {
			pytorchJob: &trainingoperator.PyTorchJob{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: trainingoperator.PyTorchJobSpec{
					ElasticPolicy: &trainingoperator.ElasticPolicy{},
					PyTorchReplicaSpecs: map[trainingoperator.ReplicaType]*trainingoperator.ReplicaSpec{
						trainingoperator.PyTorchJobReplicaTypeMaster: {
							Replicas: ptr.To[int32](1),
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  "pytorch",
											Image: "gcr.io/kubeflow-ci/pytorch-dist-mnist_test:1.0",
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: field.ErrorList{
				field.Required(pytorchReplicaSpecPath.Key(string(trainingoperator.PyTorchJobReplicaTypeWorker)), ""),
			},
		},
		"attempt to configure worker with 0 replicas": {
			pytorchJob: &trainingoperator.PyTorchJob{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
				},
				Spec: trainingoperator.PyTorchJobSpec{
					ElasticPolicy: &trainingoperator.ElasticPolicy{},
					PyTorchReplicaSpecs: map[trainingoperator.ReplicaType]*trainingoperator.ReplicaSpec{
						trainingoperator.PyTorchJobReplicaTypeWorker: {
							Replicas: ptr.To[int32](0),
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  "pytorch",
											Image: "gcr.io/kubeflow-ci/pytorch-dist-mnist_test:1.0",
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: field.ErrorList{
				field.Forbidden(pytorchReplicaSpecPath.Key(string(trainingoperator.PyTorchJobReplicaTypeWorker)).Child("replicas"), ""),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gotWarnings, gotError := validatePyTorchJob(tc.oldPytorchJob, tc.pytorchJob)
			if diff := cmp.Diff(tc.wantWarnings, gotWarnings, cmpopts.SortSlices(func(a, b string) bool { return a < b })); len(diff) != 0 {
				t.Errorf("Unexpected warnings (-want,+got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.wantErr, gotError, cmpopts.IgnoreFields(field.Error{}, "Detail", "BadValue")); len(diff) != 0 {
				t.Errorf("Unexpected errors (-want,+got):\n%s", diff)
			}
		})
	}
}
