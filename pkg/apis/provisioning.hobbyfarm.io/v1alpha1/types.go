package v1alpha1

import (
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/retries"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/genericcondition"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"repository.basistech.ru/BASIS/decort-golang-sdk/pkg/cloudapi/compute"
)

var (
	ConditionInstanceExists  = condition.Cond("InstanceExists")
	ConditionInstanceRunning = condition.Cond("InstanceRunning")
	ConditionInstanceUpdated = condition.Cond("InstanceUpdated")

	ConditionConnectionReady = condition.Cond("InstanceConnectionReady")

	ConditionKeyExists  = condition.Cond("KeyExists")
	ConditionKeyCreated = condition.Cond("DigitaEnergyKeyCreated")

	RetryDeleteKey     = retries.NewRetry("DeleteKey", 0)
	RetryDeleteDroplet = retries.NewRetry("DeleteInstance", 0)
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,inline"`

	Spec   InstanceSpec   `json:"spec"`
	Status InstanceStatus `json:"status"`
}

// +k8s:deepcopy-gen=true

type InstanceCreateRequest struct {
	compute.RecordCompute
}

// +k8s:deepcopy-gen=true

type InstanceSpec struct {
	Machine  string  `json:"machine"`
	Instance v1.JSON `json:"instance"`
}

// +k8s:deepcopy-gen=true

type InstanceStatus struct {
	Instance   v1.JSON                             `json:"instance"`
	Conditions []genericcondition.GenericCondition `json:"conditions"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Instance `json:"items,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Key struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,inline"`

	Spec   KeySpec   `json:"spec"`
	Status KeyStatus `json:"status"`
}

// +k8s:deepcopy-gen=true

type KeySpec struct {
	// HF machine with which this key is associated
	Machine string `json:"machine"`
	Secret  string `json:"secret"`

	Key v1.JSON `json:"key"`
}

// +k8s:deepcopy-gen=true

type KeyStatus struct {
	Key        v1.JSON                             `json:"key"`
	Conditions []genericcondition.GenericCondition `json:"conditions"`
	Retries    []retries.GenericRetry              `json:"retries"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type KeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Key `json:"items,omitempty"`
}
