// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VolumeGroupReplicationClassSpec defines the desired state
type VolumeGroupReplicationClassSpec struct {
	// Storage provisioner that this replication class supports
	// Example: rook-ceph.cephfs.csi.ceph.com
	// +kubebuilder:validation:Optional
	Provisioner string `json:"provisioner,omitempty"`

	// Driver specific parameters
	// +kubebuilder:validation:Optional
	Parameters map[string]string `json:"parameters,omitempty"`
}

// VolumeGroupReplicationClassStatus defines the observed state
type VolumeGroupReplicationClassStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,shortName=vgrc
// +kubebuilder:printcolumn:JSONPath=".spec.provisioner",name=Provisioner,type=string
// +kubebuilder:printcolumn:JSONPath=".metadata.creationTimestamp",name=Age,type=date

// VolumeGroupReplicationClass is the Schema for the volumegroupreplicationclasses API
type VolumeGroupReplicationClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VolumeGroupReplicationClassSpec   `json:"spec,omitempty"`
	Status VolumeGroupReplicationClassStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VolumeGroupReplicationClassList contains a list of VolumeGroupReplicationClass
type VolumeGroupReplicationClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VolumeGroupReplicationClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VolumeGroupReplicationClass{}, &VolumeGroupReplicationClassList{})
}
