// +build !ignore_autogenerated

/*
Copyright 2021.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretMapping) DeepCopyInto(out *SecretMapping) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretMapping.
func (in *SecretMapping) DeepCopy() *SecretMapping {
	if in == nil {
		return nil
	}
	out := new(SecretMapping)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecretMapping) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretMappingItem) DeepCopyInto(out *SecretMappingItem) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretMappingItem.
func (in *SecretMappingItem) DeepCopy() *SecretMappingItem {
	if in == nil {
		return nil
	}
	out := new(SecretMappingItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretMappingList) DeepCopyInto(out *SecretMappingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SecretMapping, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretMappingList.
func (in *SecretMappingList) DeepCopy() *SecretMappingList {
	if in == nil {
		return nil
	}
	out := new(SecretMappingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecretMappingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretMappingSpec) DeepCopyInto(out *SecretMappingSpec) {
	*out = *in
	out.Source = in.Source
	if in.Mappings != nil {
		in, out := &in.Mappings, &out.Mappings
		*out = make([]SecretMappingItem, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretMappingSpec.
func (in *SecretMappingSpec) DeepCopy() *SecretMappingSpec {
	if in == nil {
		return nil
	}
	out := new(SecretMappingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretMappingStatus) DeepCopyInto(out *SecretMappingStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretMappingStatus.
func (in *SecretMappingStatus) DeepCopy() *SecretMappingStatus {
	if in == nil {
		return nil
	}
	out := new(SecretMappingStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretSource) DeepCopyInto(out *SecretSource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretSource.
func (in *SecretSource) DeepCopy() *SecretSource {
	if in == nil {
		return nil
	}
	out := new(SecretSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecretSource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretSourceGitRepositoryRef) DeepCopyInto(out *SecretSourceGitRepositoryRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretSourceGitRepositoryRef.
func (in *SecretSourceGitRepositoryRef) DeepCopy() *SecretSourceGitRepositoryRef {
	if in == nil {
		return nil
	}
	out := new(SecretSourceGitRepositoryRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretSourceList) DeepCopyInto(out *SecretSourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SecretSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretSourceList.
func (in *SecretSourceList) DeepCopy() *SecretSourceList {
	if in == nil {
		return nil
	}
	out := new(SecretSourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SecretSourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretSourcePrivateKeyRef) DeepCopyInto(out *SecretSourcePrivateKeyRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretSourcePrivateKeyRef.
func (in *SecretSourcePrivateKeyRef) DeepCopy() *SecretSourcePrivateKeyRef {
	if in == nil {
		return nil
	}
	out := new(SecretSourcePrivateKeyRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretSourceRef) DeepCopyInto(out *SecretSourceRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretSourceRef.
func (in *SecretSourceRef) DeepCopy() *SecretSourceRef {
	if in == nil {
		return nil
	}
	out := new(SecretSourceRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretSourceSpec) DeepCopyInto(out *SecretSourceSpec) {
	*out = *in
	out.GitRepository = in.GitRepository
	out.PrivateKey = in.PrivateKey
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretSourceSpec.
func (in *SecretSourceSpec) DeepCopy() *SecretSourceSpec {
	if in == nil {
		return nil
	}
	out := new(SecretSourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretSourceStatus) DeepCopyInto(out *SecretSourceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretSourceStatus.
func (in *SecretSourceStatus) DeepCopy() *SecretSourceStatus {
	if in == nil {
		return nil
	}
	out := new(SecretSourceStatus)
	in.DeepCopyInto(out)
	return out
}
