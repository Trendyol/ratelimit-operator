package istio

import (
	"context"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

//
import (
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
)

type Istio interface {
    GetEnvoyFilter(ctx context.Context, namespace string, name string) (*v1alpha3.EnvoyFilter, error)
	CreateEnvoyFilter(ctx context.Context, namespace string, envoyFilter *v1alpha3.EnvoyFilter) (*v1alpha3.EnvoyFilter, error)
	DeleteEnvoyFilter(ctx context.Context, namespace, name string) error
	PatchEnvoyFilter(ctx context.Context, data []byte, namespace, name string) (*v1alpha3.EnvoyFilter, error)
}
type istioClient struct {
	cfg    *rest.Config
	client versionedclient.Interface
}

func NewIstioClient(cfg *rest.Config) Istio {
	clientSet := versionedclient.NewForConfigOrDie(cfg)
	return &istioClient{cfg: cfg, client: clientSet}
}
func (r *istioClient) DeleteEnvoyFilter(ctx context.Context, namespace, name string) error {
	return r.client.NetworkingV1alpha3().EnvoyFilters(namespace).Delete(ctx, name, v1.DeleteOptions{})
}

func (r *istioClient) GetEnvoyFilter(ctx context.Context, namespace string, name string) (*v1alpha3.EnvoyFilter, error) {
    return r.client.NetworkingV1alpha3().EnvoyFilters(namespace).Get(ctx, name, v1.GetOptions{})
}

func (r *istioClient) PatchEnvoyFilter(ctx context.Context, data []byte, namespace, name string) (*v1alpha3.EnvoyFilter, error) {
	return r.client.NetworkingV1alpha3().EnvoyFilters(namespace).Patch(ctx, name, types.MergePatchType, data, v1.PatchOptions{})
}
func (r *istioClient) CreateEnvoyFilter(ctx context.Context, namespace string, envoyFilter *v1alpha3.EnvoyFilter) (*v1alpha3.EnvoyFilter, error) {
	return r.client.NetworkingV1alpha3().EnvoyFilters(namespace).Create(ctx, envoyFilter, v1.CreateOptions{})
}
