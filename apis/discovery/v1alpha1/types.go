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

type PeerId string

type PeerInfo struct {
	ID        string   `json:"id,omitempty"`
	Addresses []string `json:"addresses,omitempty"`
}

type PingRequest struct {
	Info *PeerInfo `json:"info,omitempty"`
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
	ResourceKindJoinCluster     = "JoinCluster"
	ResourcePluralJoinCluster   = "joinclusters"
	ResourceSingularJoinCluster = "joincluster"
)

type EtcdNode struct {
	Name                  string   `json:"name,omitempty"`
	PeerURLs              []string `json:"peerURLs,omitempty"`
	ClientURLs            []string `json:"clientURLs,omitempty"`
	QuarantinedClientURLs []string `json:"quarantinedClientURLs,omitempty"`
}

type JoinClusterRequest struct {
	LeadershipToken string      `json:"leadershipToken,omitempty"`
	ClusterName     string      `json:"clusterName,omitempty"`
	ClusterToken    string      `json:"clusterToken,omitempty"`
	Nodes           []*EtcdNode `json:"nodes,omitempty"`
	AddNode         *EtcdNode   `json:"addNode,omitempty"`
	EtcdVersion     string      `json:"etcdVersion,omitempty"`
}

type JoinClusterResponse struct {
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=get,list,update,patch,delete,deleteCollection,watch
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type JoinCluster struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	Request *JoinClusterRequest `json:"request,omitempty"`
	// +optional
	Response *JoinClusterResponse `json:"response,omitempty"`
}
