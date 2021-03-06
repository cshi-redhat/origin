package etcd

import (
	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	etcdgeneric "github.com/GoogleCloudPlatform/kubernetes/pkg/registry/generic/etcd"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/tools"

	"github.com/openshift/origin/pkg/template/api"
	"github.com/openshift/origin/pkg/template/registry"
)

// REST implements a RESTStorage for templates against etcd
type REST struct {
	*etcdgeneric.Etcd
}

// NewREST returns a RESTStorage object that will work against templates.
func NewREST(h tools.EtcdHelper) *REST {
	store := &etcdgeneric.Etcd{
		NewFunc:     func() runtime.Object { return &api.Template{} },
		NewListFunc: func() runtime.Object { return &api.TemplateList{} },
		KeyRootFunc: func(ctx kapi.Context) string {
			return etcdgeneric.NamespaceKeyRootFunc(ctx, "/registry/templates")
		},
		KeyFunc: func(ctx kapi.Context, name string) (string, error) {
			return etcdgeneric.NamespaceKeyFunc(ctx, "/registry/templates", name)
		},
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*api.Template).Name, nil
		},
		EndpointName: "templates",

		CreateStrategy: registry.Strategy,
		UpdateStrategy: registry.Strategy,

		ReturnDeletedObject: true,

		Helper: h,
	}
	return &REST{store}
}

// New returns a new object
func (r *REST) New() runtime.Object {
	return r.NewFunc()
}

// NewList returns a new list object
func (r *REST) NewList() runtime.Object {
	return r.NewListFunc()
}

// List obtains a list of templates with labels that match selector.
func (r *REST) List(ctx kapi.Context, label labels.Selector, field fields.Selector) (runtime.Object, error) {
	return r.Etcd.ListPredicate(ctx, registry.MatchTemplate(label, field))
}
