// SPDX-FileCopyrightText: The RamenDR authors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ReplicationState represents the replication operations to be performed on the volume
type ReplicationState string

const (
	// Promote the protected PVCs to primary
	Primary ReplicationState = "primary"

	// Demote the protected PVCs to secondary
	Secondary ReplicationState = "secondary"
)

// VolumeGroupReplicationSpec defines the desired state
type VolumeGroupReplicationSpec struct {
	// ReplicationState indicates whether this cluster is primary or secondary
	ReplicationState ReplicationState `json:"replicationState"`

	// VolumeGroupReplicationClassName references the VGRC used for replication
	VolumeGroupReplicationClassName string `json:"volumeGroupReplicationClassName"`

	// External indicates the replication is handled by external vendor controller
	External bool `json:"external,omitempty"`

	// Source selector identifies PVCs that belong to the group
	Source VolumeGroupReplicationSource `json:"source"`
}

type VolumeGroupReplicationSource struct {
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

// VolumeGroupReplicationStatus defines the observed state
type VolumeGroupReplicationStatus struct {
	// State represents the current replication state
	State ReplicationState `json:"state,omitempty"`

	// Conditions represent the latest available observations
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=vgr
// +kubebuilder:printcolumn:JSONPath=".spec.replicationState",name=DesiredState,type=string
// +kubebuilder:printcolumn:JSONPath=".status.state",name=CurrentState,type=string
// +kubebuilder:printcolumn:JSONPath=".metadata.creationTimestamp",name=Age,type=date

// VolumeGroupReplication is the Schema for the volumegroupreplications API
type VolumeGroupReplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VolumeGroupReplicationSpec   `json:"spec,omitempty"`
	Status VolumeGroupReplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VolumeGroupReplicationList contains a list of VolumeGroupReplication
type VolumeGroupReplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VolumeGroupReplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VolumeGroupReplication{}, &VolumeGroupReplicationList{})
}
