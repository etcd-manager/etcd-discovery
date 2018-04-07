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

package discovery

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type PeerInfo struct {
	ID    string
	Hosts []string
}

type PingRequest struct {
}

type PingResponse struct {
	Info *PeerInfo
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=get,list,update,patch,delete,deleteCollection,watch
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Ping describes a peer ping request/response.
type Ping struct {
	metav1.TypeMeta
	// +optional
	Request *PingRequest
	// +optional
	Response *PingResponse
}

type MemberRequest struct {
	PeerURL string
}

type MemberResponse struct {
	ClusterName  string
	ClusterToken string
	PeerURLs     []string
	EtcdVersion  string
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=get,list,update,patch,delete,deleteCollection,watch
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Member struct {
	metav1.TypeMeta
	// +optional
	Request *MemberRequest
	// +optional
	Response *MemberResponse
}
