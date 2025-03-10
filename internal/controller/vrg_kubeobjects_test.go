// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

// white box testing desired for Recipe/KubeObject conversions
package controllers //nolint:testpackage

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ramendr/ramen/internal/controller/kubeobjects"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	Recipe "github.com/ramendr/recipe/api/v1alpha1"
)

var _ = Describe("VRG_KubeObjectProtection", func() {
	const namespaceName = "my-ns"

	var hook *Recipe.Hook
	var group *Recipe.Group

	BeforeEach(func() {
		duration := 30

		hook = &Recipe.Hook{
			Namespace: namespaceName,
			Name:      "hook-single",
			Type:      "exec",
			LabelSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"myapp": "testapp",
				},
			},
			SinglePodOnly: false,
			Ops: []*Recipe.Operation{
				{
					Name:      "checkpoint",
					Container: "main",
					Timeout:   duration,
					Command:   "bash /scripts/checkpoint.sh",
				},
			},
			Chks:      []*Recipe.Check{},
			Essential: new(bool),
		}

		group = &Recipe.Group{
			Name:                  "test-group",
			BackupRef:             "test-backup-ref",
			Type:                  "resource",
			IncludedNamespaces:    []string{namespaceName},
			IncludedResourceTypes: []string{"deployment", "replicaset"},
			ExcludedResourceTypes: nil,
			LabelSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "test",
						Operator: metav1.LabelSelectorOpNotIn,
						Values:   []string{"empty-on-backup notin", "ignore-on-backup"},
					},
				},
			},
		}
	})

	Context("Conversion", func() {
		It("Hook to CaptureSpec", func() {
			targetCaptureSpec := &kubeobjects.CaptureSpec{
				Name: hook.Name + "-" + hook.Ops[0].Name,
				Spec: kubeobjects.Spec{
					KubeResourcesSpec: kubeobjects.KubeResourcesSpec{
						IncludedNamespaces: []string{namespaceName},
						IncludedResources:  []string{"pod"},
						ExcludedResources:  []string{},
						Hook: kubeobjects.HookSpec{
							Name:          hook.Name,
							Namespace:     namespaceName,
							Type:          hook.Type,
							LabelSelector: hook.LabelSelector,
							Op: kubeobjects.Operation{
								Name:      hook.Ops[0].Name,
								Command:   hook.Ops[0].Command,
								Container: hook.Ops[0].Container,
							},
						},
						IsHook: true,
					},
					LabelSelector:           hook.LabelSelector,
					IncludeClusterResources: new(bool),
				},
			}
			converted, err := convertRecipeHookToCaptureSpec(*hook, hook.Ops[0].Name)

			Expect(err).To(BeNil())
			Expect(converted).To(Equal(targetCaptureSpec))
		})

		It("Hook to RecoverSpec", func() {
			targetRecoverSpec := &kubeobjects.RecoverSpec{
				Spec: kubeobjects.Spec{
					KubeResourcesSpec: kubeobjects.KubeResourcesSpec{
						IncludedNamespaces: []string{namespaceName},
						IncludedResources:  []string{"pod"},
						ExcludedResources:  []string{},
						Hook: kubeobjects.HookSpec{
							Name:      hook.Name,
							Type:      hook.Type,
							Namespace: namespaceName,
							// Timeout:       hook.Ops[0].Timeout,
							LabelSelector: hook.LabelSelector,
							Op: kubeobjects.Operation{
								Name:      hook.Ops[0].Name,
								Command:   hook.Ops[0].Command,
								Container: hook.Ops[0].Container,
							},
						},
						IsHook: true,
					},
					LabelSelector:           hook.LabelSelector,
					IncludeClusterResources: new(bool),
				},
			}
			converted, err := convertRecipeHookToRecoverSpec(*hook, hook.Ops[0].Name)

			Expect(err).To(BeNil())
			Expect(converted).To(Equal(targetRecoverSpec))
		})

		It("Group to CaptureSpec", func() {
			targetCaptureSpec := &kubeobjects.CaptureSpec{
				Name: group.Name,
				Spec: kubeobjects.Spec{
					KubeResourcesSpec: kubeobjects.KubeResourcesSpec{
						IncludedNamespaces: group.IncludedNamespaces,
						IncludedResources:  group.IncludedResourceTypes,
						ExcludedResources:  group.ExcludedResourceTypes,
					},
					LabelSelector:           group.LabelSelector,
					IncludeClusterResources: group.IncludeClusterResources,
					OrLabelSelectors:        []*metav1.LabelSelector{},
				},
			}
			converted, err := convertRecipeGroupToCaptureSpec(*group)

			Expect(err).To(BeNil())
			Expect(converted).To(Equal(targetCaptureSpec))
		})

		It("Group to RecoverSpec", func() {
			targetRecoverSpec := &kubeobjects.RecoverSpec{
				BackupName: group.BackupRef,
				Spec: kubeobjects.Spec{
					KubeResourcesSpec: kubeobjects.KubeResourcesSpec{
						IncludedNamespaces: group.IncludedNamespaces,
						IncludedResources:  group.IncludedResourceTypes,
						ExcludedResources:  group.ExcludedResourceTypes,
					},
					LabelSelector:           group.LabelSelector,
					IncludeClusterResources: group.IncludeClusterResources,
					OrLabelSelectors:        []*metav1.LabelSelector{},
				},
			}
			converted, err := convertRecipeGroupToRecoverSpec(*group)

			Expect(err).To(BeNil())
			Expect(converted).To(Equal(targetRecoverSpec))
		})
	})
})
