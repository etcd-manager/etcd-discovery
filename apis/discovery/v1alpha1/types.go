/*
Copyright 2018 The Pharmer Authors.

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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type PeerID string

type PeerInfo struct {
	ID    string   `json:"id,omitempty"`
	Hosts []string `json:"hosts,omitempty"`
}

type PingRequest struct {
}

type PingResponse struct {
	Info *PeerInfo `json:"info,omitempty"`
}

const (
	ResourceKindPing     = "Ping"
	ResourcePluralPing   = "pings"
	ResourceSingularPing = "ping"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=get,list,update,patch,delete,deleteCollection,watch
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Ping describes a peer ping request/response.
type Ping struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	Request *PingRequest `json:"request,omitempty"`
	// +optional
	Response *PingResponse `json:"response,omitempty"`
}

const (
	ResourceKindMember     = "Member"
	ResourcePluralMember   = "members"
	ResourceSingularMember = "member"
)

type MemberRequest struct {
	PeerURL string `json:"peerURL,omitempty"`
}

type MemberResponse struct {
	ClusterName  string   `json:"clusterName,omitempty"`
	ClusterToken string   `json:"clusterToken,omitempty"`
	PeerURLs     []string `json:"peerURLs,omitempty"`
	EtcdVersion  string   `json:"etcdVersion,omitempty"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=get,list,update,patch,delete,deleteCollection,watch
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Member struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	Request *MemberRequest `json:"request,omitempty"`
	// +optional
	Response *MemberResponse `json:"response,omitempty"`
}
