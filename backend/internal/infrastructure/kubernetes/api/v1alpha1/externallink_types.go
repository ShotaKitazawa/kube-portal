/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExternalLinkSpec defines the desired state of ExternalLink
type ExternalLinkSpec struct {
	Title     string   `json:"title"`
	Href      string   `json:"href"`
	IconURL   string   `json:"iconURL,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	IsPrivate bool     `json:"isPrivate,omitempty"`
}

// ExternalLinkStatus defines the observed state of ExternalLink
type ExternalLinkStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ExternalLink is the Schema for the externallinks API
type ExternalLink struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExternalLinkSpec   `json:"spec,omitempty"`
	Status ExternalLinkStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ExternalLinkList contains a list of ExternalLink
type ExternalLinkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ExternalLink `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ExternalLink{}, &ExternalLinkList{})
}
