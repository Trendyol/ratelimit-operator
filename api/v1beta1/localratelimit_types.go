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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LocalRateLimitSpec defines the desired state of LocalRateLimit
type LocalRateLimitSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of LocalRateLimit. Edit localratelimit_types.go to remove/update
	Foo         string      `json:"foo,omitempty"`
	Workload    string      `json:"workload,omitempty"`
	TokenBucket TokenBucket `json:"token_bucket,omitempty"`
}

// LocalRateLimitStatus defines the observed state of LocalRateLimit
type LocalRateLimitStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LocalRateLimit is the Schema for the localratelimits API
type LocalRateLimit struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LocalRateLimitSpec   `json:"spec,omitempty"`
	Status LocalRateLimitStatus `json:"status,omitempty"`
}

type TokenBucket struct {
	MaxTokens     int16  `json:"max_tokens,required"`
	TokensPerFill int16  `json:"tokens_per_fill,required"`
	FillInterval  string `json:"fill_interval,required"`
}

//+kubebuilder:object:root=true

// LocalRateLimitList contains a list of LocalRateLimit
type LocalRateLimitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LocalRateLimit `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LocalRateLimit{}, &LocalRateLimitList{})
}
