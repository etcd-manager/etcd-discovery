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
package fake

import (
	v1alpha1 "github.com/etcd-manager/etcd-discovery/apis/discovery/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
)

// FakeMembers implements MemberInterface
type FakeMembers struct {
	Fake *FakeDiscoveryV1alpha1
}

var membersResource = schema.GroupVersionResource{Group: "discovery.etcd-manager.com", Version: "v1alpha1", Resource: "members"}

var membersKind = schema.GroupVersionKind{Group: "discovery.etcd-manager.com", Version: "v1alpha1", Kind: "Member"}

// Create takes the representation of a member and creates it.  Returns the server's representation of the member, and an error, if there is any.
func (c *FakeMembers) Create(member *v1alpha1.Member) (result *v1alpha1.Member, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(membersResource, member), &v1alpha1.Member{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Member), err
}
