//go:build !ignore_autogenerated



// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/negashev/hf-provisioner-digitalenergy/pkg/retries"
	"github.com/rancher/wrangler/pkg/genericcondition"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Instance) DeepCopyInto(out *Instance) {
*out = *in
out.TypeMeta = in.TypeMeta
in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
in.Spec.DeepCopyInto(&out.Spec)
in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Instance.
func (in *Instance) DeepCopy() *Instance {
	if in == nil { return nil }
	out := new(Instance)
	in.DeepCopyInto(out)
	return out
}


// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Instance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceList) DeepCopyInto(out *InstanceList) {
*out = *in
out.TypeMeta = in.TypeMeta
in.ListMeta.DeepCopyInto(&out.ListMeta)
if in.Items != nil {
in, out := &in.Items, &out.Items
*out = make([]Instance, len(*in))
for i := range *in {
(*in)[i].DeepCopyInto(&(*out)[i])
}
}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceList.
func (in *InstanceList) DeepCopy() *InstanceList {
	if in == nil { return nil }
	out := new(InstanceList)
	in.DeepCopyInto(out)
	return out
}


// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceSpec) DeepCopyInto(out *InstanceSpec) {
*out = *in
in.Instance.DeepCopyInto(&out.Instance)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceSpec.
func (in *InstanceSpec) DeepCopy() *InstanceSpec {
	if in == nil { return nil }
	out := new(InstanceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceStatus) DeepCopyInto(out *InstanceStatus) {
*out = *in
in.Instance.DeepCopyInto(&out.Instance)
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]genericcondition.GenericCondition, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceStatus.
func (in *InstanceStatus) DeepCopy() *InstanceStatus {
	if in == nil { return nil }
	out := new(InstanceStatus)
	in.DeepCopyInto(out)
	return out
}


// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Key) DeepCopyInto(out *Key) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Key.
func (in *Key) DeepCopy() *Key {
	if in == nil {
		return nil
	}
	out := new(Key)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Key) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeyList) DeepCopyInto(out *KeyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Key, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyList.
func (in *KeyList) DeepCopy() *KeyList {
	if in == nil {
		return nil
	}
	out := new(KeyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KeyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeySpec) DeepCopyInto(out *KeySpec) {
	*out = *in
	in.Key.DeepCopyInto(&out.Key)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeySpec.
func (in *KeySpec) DeepCopy() *KeySpec {
	if in == nil {
		return nil
	}
	out := new(KeySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeyStatus) DeepCopyInto(out *KeyStatus) {
	*out = *in
	in.Key.DeepCopyInto(&out.Key)
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]genericcondition.GenericCondition, len(*in))
		copy(*out, *in)
	}
	if in.Retries != nil {
		in, out := &in.Retries, &out.Retries
		*out = make([]retries.GenericRetry, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeyStatus.
func (in *KeyStatus) DeepCopy() *KeyStatus {
	if in == nil {
		return nil
	}
	out := new(KeyStatus)
	in.DeepCopyInto(out)
	return out
}
