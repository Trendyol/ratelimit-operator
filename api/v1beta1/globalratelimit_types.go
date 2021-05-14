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

// GlobalRateLimitSpec defines the desired state of GlobalRateLimit
type GlobalRateLimitSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Domain   string `json:"domain,required"`
	Workload string `json:"workload,required"`
	Rate     []Rate `json:"rate,required"`
}

type Rate struct {
	Unit           string       `json:"unit,required"`
	RequestPerUnit int64        `json:"request_per_unit,required"`
	Dimensions     []Dimensions `json:"dimensions,required"`
}
type Dimensions struct {
	RequestHeader      *RequestHeader      `json:"request_headers,omitempty"`
	HeaderValueMatch   *HeaderValueMatch   `json:"header_value_match,omitempty"`
	GenericKey         *GenericKey         `json:"generic_key,omitempty"`
	SourceCluster      *SourceCluster      `json:"source_cluster,omitempty"`
	DestinationCluster *DestinationCluster `json:"destination_cluster,omitempty"`
	RemoteAddress      *RemoteAddress              `json:"remote_address,omitempty"`
}

type HeaderValueMatch struct {
	DescriptorValue string          `json:"descriptor_value,required"`
	HeaderMatcher   []HeaderMatcher `json:"headers,required"`
	//QueryMatcher    []QueryMatcher  `json:"query_matcher,required"`
}

type HeaderMatcher struct {
	Name           string          `json:"name,required"`
	ExactMatch     string          `json:"exact_match,omitempty"`
	SafeRegexMatch *SafeRegexMatch `json:"safe_regex_match,omitempty"`
	RangeMatch     int64           `json:"range_match,omitempty"`
	PresentMatch   int64           `json:"present_match,omitempty"`
	PrefixMatch    string          `json:"prefix_match,omitempty"`
	SuffixMatch    string          `json:"suffix_match,omitempty"`
	ContainsMatch  bool            `json:"contains_match,omitempty"`
	InvertMatch    int64           `json:"invert_match,omitempty"`
}

//TODO: How to impl querymatcher
//type QueryMatcher struct {
//	Name         string `json:"name,required"`
//	StringMatch  string `json:"string_match,omitempty"`
//	PresentMatch string `json:"present_match,omitempty"`
//}

type GoogleRE2 struct{}
type SafeRegexMatch struct {
	GoogleRe2 GoogleRE2 `json:"google_re2,required"`
	Regex     string    `json:"regex,required"`
}

type RequestHeader struct {
	DescriptorKey string `json:"descriptor_key,required"`
	HeaderName    string `json:"header_name,required"`
	Value         string `json:"value,omitempty"`
	SkipIfAbsent         string `json:"skip_if_absent,omitempty"`
}

type GenericKey struct {
	DescriptorValue string `json:"descriptor_value,required"`
	DescriptorKey   string `json:"descriptor_key,omitempty"`
}

type RemoteAddress struct{}
type SourceCluster struct{}
type DestinationCluster struct{}

type RateLimitAction struct {
	RateLimits []RateLimits `json:"rate_limits"`
}

type RateLimits struct {
	Actions []Actions `json:"actions"`
}
type Actions struct {
	RequestHeader      *RequestHeader      `json:"request_headers,omitempty"`
	HeaderValueMatch   *HeaderValueMatch   `json:"header_value_match,omitempty"`
	GenericKey         *GenericKey         `json:"generic_key,omitempty"`
	SourceCluster      *SourceCluster      `json:"source_cluster,omitempty"`
	DestinationCluster *DestinationCluster `json:"destination_cluster,omitempty"`
	RemoteAddress      *RemoteAddress              `json:"remote_address,omitempty"`
}

// GlobalRateLimitStatus defines the observed state of GlobalRateLimit
type GlobalRateLimitStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GlobalRateLimit is the Schema for the globalratelimits API
type GlobalRateLimit struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GlobalRateLimitSpec   `json:"spec,omitempty"`
	Status GlobalRateLimitStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GlobalRateLimitList contains a list of GlobalRateLimit
type GlobalRateLimitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlobalRateLimit `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GlobalRateLimit{}, &GlobalRateLimitList{})
}
