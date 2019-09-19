// +build !ignore_autogenerated

// Code generated by operator-sdk-v0.10.0-x86_64-linux-gnu. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlightRecorder) DeepCopyInto(out *FlightRecorder) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlightRecorder.
func (in *FlightRecorder) DeepCopy() *FlightRecorder {
	if in == nil {
		return nil
	}
	out := new(FlightRecorder)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FlightRecorder) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlightRecorderList) DeepCopyInto(out *FlightRecorderList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FlightRecorder, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlightRecorderList.
func (in *FlightRecorderList) DeepCopy() *FlightRecorderList {
	if in == nil {
		return nil
	}
	out := new(FlightRecorderList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FlightRecorderList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlightRecorderSpec) DeepCopyInto(out *FlightRecorderSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlightRecorderSpec.
func (in *FlightRecorderSpec) DeepCopy() *FlightRecorderSpec {
	if in == nil {
		return nil
	}
	out := new(FlightRecorderSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlightRecorderStatus) DeepCopyInto(out *FlightRecorderStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlightRecorderStatus.
func (in *FlightRecorderStatus) DeepCopy() *FlightRecorderStatus {
	if in == nil {
		return nil
	}
	out := new(FlightRecorderStatus)
	in.DeepCopyInto(out)
	return out
}