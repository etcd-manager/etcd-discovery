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

// MembersGetter has a method to return a MemberInterface.
// A group's client should implement this interface.
type MembersGetter interface {
	Members() MemberInterface
}

// MemberInterface has methods to work with Member resources.
type MemberInterface interface {
	Create(*v1alpha1.Member) (*v1alpha1.Member, error)
	MemberExpansion
}

// members implements MemberInterface
type members struct {
	client rest.Interface
}

// newMembers returns a Members
func newMembers(c *DiscoveryV1alpha1Client) *members {
	return &members{
		client: c.RESTClient(),
	}
}

// Create takes the representation of a member and creates it.  Returns the server's representation of the member, and an error, if there is any.
func (c *members) Create(member *v1alpha1.Member) (result *v1alpha1.Member, err error) {
	result = &v1alpha1.Member{}
	err = c.client.Post().
		Resource("members").
		Body(member).
		Do().
		Into(result)
	return
}
