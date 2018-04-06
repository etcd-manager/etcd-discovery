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

import (
	v1alpha1 "github.com/etcd-manager/etcd-discovery/apis/discovery/v1alpha1"
	rest "k8s.io/client-go/rest"
)

// JoinClustersGetter has a method to return a JoinClusterInterface.
// A group's client should implement this interface.
type JoinClustersGetter interface {
	JoinClusters() JoinClusterInterface
}

// JoinClusterInterface has methods to work with JoinCluster resources.
type JoinClusterInterface interface {
	Create(*v1alpha1.JoinCluster) (*v1alpha1.JoinCluster, error)
	JoinClusterExpansion
}

// joinClusters implements JoinClusterInterface
type joinClusters struct {
	client rest.Interface
}

// newJoinClusters returns a JoinClusters
func newJoinClusters(c *DiscoveryV1alpha1Client) *joinClusters {
	return &joinClusters{
		client: c.RESTClient(),
	}
}

// Create takes the representation of a joinCluster and creates it.  Returns the server's representation of the joinCluster, and an error, if there is any.
func (c *joinClusters) Create(joinCluster *v1alpha1.JoinCluster) (result *v1alpha1.JoinCluster, err error) {
	result = &v1alpha1.JoinCluster{}
	err = c.client.Post().
		Resource("joinclusters").
		Body(joinCluster).
		Do().
		Into(result)
	return
}
