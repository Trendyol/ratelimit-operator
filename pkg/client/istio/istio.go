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

const EnvoyFilterNamespace = "istio-system"

type IstioClient interface {
	CreateEnvoyFilter(ctx context.Context, namespace string, envoyFilter *v1alpha3.EnvoyFilter) (*v1alpha3.EnvoyFilter, error)
	DeleteEnvoyFilter(ctx context.Context, namespace, name string) error
	PatchEnvoyFilter(ctx context.Context, data []byte, namespace, name string) (*v1alpha3.EnvoyFilter, error)
}
type istioClient struct {
	cfg    *rest.Config
	client versionedclient.Interface
}

func NewIstioClient(cfg *rest.Config) IstioClient {
	clientSet := versionedclient.NewForConfigOrDie(cfg)
	return &istioClient{cfg: cfg, client: clientSet}
}
func (r *istioClient) DeleteEnvoyFilter(ctx context.Context, namespace, name string) error {
	return r.client.NetworkingV1alpha3().EnvoyFilters(namespace).Delete(ctx, name, v1.DeleteOptions{})
}

func (r *istioClient) PatchEnvoyFilter(ctx context.Context, data []byte, namespace, name string) (*v1alpha3.EnvoyFilter, error) {
	return r.client.NetworkingV1alpha3().EnvoyFilters(namespace).Patch(ctx, name, types.MergePatchType, data, v1.PatchOptions{})
}
func (r *istioClient) CreateEnvoyFilter(ctx context.Context, namespace string, envoyFilter *v1alpha3.EnvoyFilter) (*v1alpha3.EnvoyFilter, error) {
	return r.client.NetworkingV1alpha3().EnvoyFilters(namespace).Create(ctx, envoyFilter, v1.CreateOptions{})
	//a := nv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch{
	//	ApplyTo: nv1alpha3.EnvoyFilter_HTTP_FILTER,
	//	Patch: &nv1alpha3.EnvoyFilter_Patch{
	//		Operation: nv1alpha3.EnvoyFilter_Patch_INSERT_BEFORE,
	//		Value: &types.Struct{
	//			Fields: map[string]*types.Value{
	//				"name": {
	//					Kind: &types.Value_StringValue{
	//						StringValue: "envoy.filters.http.local_ratelimit",
	//					},
	//				},
	//				"typed_config": {
	//					Kind: &types.Value_StructValue{
	//						StructValue:types.str
	//					},
	//				},
	//			},
	//		},
	//	},
	//	//Match: nv1alpha3.EnvoyFilter_ListenerMatch{}
	//}
	//
	//var configPacthes []*nv1alpha3.EnvoyFilter_EnvoyConfigObjectPatch
	//e := &v1alpha3.EnvoyFilter{
	//	Spec: nv1alpha3.EnvoyFilter{
	//		WorkloadSelector: &nv1alpha3.WorkloadSelector{
	//			Labels: map[string]string{
	//				"app": localRateLimit.Spec.Workload,
	//			},
	//		},
	//		ConfigPatches: configPacthes,
	//	},
	//}

}