package global

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/api/v1beta1"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/pkg/client/istio"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/klog"
)

type GlobalRateLimitFilter struct {
	istio istio.IstioClient
}

func NewGlobalRateLimitFilter(istio istio.IstioClient) *GlobalRateLimitFilter {
	return &GlobalRateLimitFilter{istio: istio}
}

var globalRateLimitEnvoyFilter = `
kind: EnvoyFilter
apiVersion: networking.istio.io/v1alpha3
metadata:
  labels:
    generator: ratelimit-operator
  name: details-app-demo-ratelimit
  namespace: default
spec:
  configPatches:
    - applyTo: HTTP_FILTER
      match:
        context: SIDECAR_INBOUND
        listener:
          filterChain:
            filter:
              name: envoy.http_connection_manager
              subFilter:
                name: envoy.router
      patch:
        operation: INSERT_BEFORE
        value:
          name: envoy.filters.ratelimit
          typed_config:
            '@type': >-
              type.googleapis.com/envoy.extensions.filters.http.ratelimit.v3.RateLimit
            domain: productpage600-ratelimit
            failure_mode_deny: false
            rate_limit_service:
              grpc_service:
                envoy_grpc:
                  cluster_name: rate_limit_cluster
                timeout: 10s
              transport_api_version: V3
    - applyTo: CLUSTER
      match:
        cluster:
          service: ratelimit.default
      patch:
        operation: ADD
        value:
          connect_timeout: 10s
          http2_protocol_options: {}
          lb_policy: ROUND_ROBIN
          load_assignment:
            cluster_name: rate_limit_cluster
            endpoints:
              - lb_endpoints:
                  - endpoint:
                      address:
                        socket_address:
                          address: ratelimit.default
                          port_value: 8081
          name: rate_limit_cluster
          type: STRICT_DNS
  workloadSelector:
    labels:
      app: details
`

func (r *GlobalRateLimitFilter) PrepareUpdateEnvoyFilterExternalObjects(ctx context.Context, global *v1beta1.GlobalRateLimit, namespace, name string) {
	var err error

	patchValue, envoyFilterObj, err := getGlobalRateLimitEnvoyFilter(namespace, global)

	_, err = r.istio.GetEnvoyFilter(ctx, namespace, name)
	if err != nil {
		klog.Infof("Envoyfilter %s is not found. Error %v", name, err)
		klog.Infof("Creating Envoyfilter %s", name)
		_, err = r.istio.CreateEnvoyFilter(ctx, namespace, envoyFilterObj)

		if err != nil {
			klog.Infof("Cannot get Ratelimit CR %s. Error %v", global.Name, err)
		}
	} else {
		_, err := r.istio.PatchEnvoyFilter(ctx, patchValue, namespace, name)
		klog.Infof("Patching Envoyfilter %s", name)

		if err != nil {
			klog.Infof("Cannot path Ratelimit CR %s. Error %v", global.Name, err)
		}
	}

}

func getGlobalRateLimitEnvoyFilter(namespace string, limit *v1beta1.GlobalRateLimit) ([]byte, *v1alpha3.EnvoyFilter, error) {
	domain := limit.Spec.Domain
	envoyFilter := v1alpha3.EnvoyFilter{}

	printf := fmt.Sprintf(globalRateLimitEnvoyFilter, limit.Name, namespace, domain, limit.Spec.Workload)
	byte := bytes.NewBufferString(printf).Bytes()
	err := json.Unmarshal(byte, &envoyFilter)
	if err != nil {
		return nil, nil, err
	}
	return byte, &envoyFilter, nil
}
