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
	v1alpha1 "github.com/etcd-manager/etcd-discovery/client/clientset/versioned/typed/discovery/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeDiscoveryV1alpha1 struct {
	*testing.Fake
}

func (c *FakeDiscoveryV1alpha1) Members() v1alpha1.MemberInterface {
	return &FakeMembers{c}
}

func (c *FakeDiscoveryV1alpha1) Pings() v1alpha1.PingInterface {
	return &FakePings{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeDiscoveryV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
