package ping

import (
	"net"

	api "github.com/etcd-manager/etcd-discovery/apis/discovery/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
)

type REST struct {
	id   api.PeerID
	host net.IP
}

var _ rest.Creater = &REST{}
var _ rest.GroupVersionKindProvider = &REST{}

func NewREST(id api.PeerID, host net.IP) *REST {
	return &REST{id, host}
}

func (r *REST) New() runtime.Object {
	return &api.Ping{}
}

func (r *REST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return api.SchemeGroupVersion.WithKind(api.ResourceKindPing)
}

func (r *REST) Create(ctx apirequest.Context, obj runtime.Object, _ rest.ValidateObjectFunc, _ bool) (runtime.Object, error) {
	req := obj.(*api.Ping)

	req.Response = &api.PingResponse{
		Info: &api.PeerInfo{
			ID:    string(r.id),
			Hosts: []string{r.host.String()},
		},
	}
	return req, nil
}
