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

var RlEnvoyFilterSuffixName = "-ratelimit-filter"

type RlFilter interface {
	PrepareUpdateEnvoyFilterExternalObjects(ctx context.Context, global *v1beta1.GlobalRateLimit)
}

type globalRateLimitFilter struct {
	istio istio.Istio
}

func NewGlobalRateLimitFilter(istio istio.Istio) RlFilter {
	return &globalRateLimitFilter{istio: istio}
}

var globalRateLimitEnvoyFilter = `
{
  "kind": "EnvoyFilter",
  "apiVersion": "networking.istio.io/v1alpha3",
  "metadata": {
    "labels": {
      "generator": "ratelimit-operator"
    },
    "name": "%s",
    "namespace": "%s"
  },
  "spec": {
    "configPatches": [
      {
        "applyTo": "HTTP_FILTER",
        "match": {
          "context": "SIDECAR_INBOUND",
          "listener": {
            "filterChain": {
              "filter": {
                "name": "envoy.http_connection_manager",
                "subFilter": {
                  "name": "envoy.router"
                }
              }
            }
          }
        },
        "patch": {
          "operation": "INSERT_BEFORE",
          "value": {
            "name": "envoy.filters.ratelimit",
            "typed_config": {
              "@type": "type.googleapis.com/envoy.extensions.filters.http.ratelimit.v3.RateLimit",
              "domain": "%s",
              "failure_mode_deny": false,
              "rate_limit_service": {
                "grpc_service": {
                  "envoy_grpc": {
                    "cluster_name": "rate_limit_cluster"
                  }
                },
                "transport_api_version": "V3"
              }
            }
          }
        }
      },
      {
        "applyTo": "CLUSTER",
        "match": {
          "cluster": {
            "service": "ratelimit.default"
          }
        },
        "patch": {
          "operation": "ADD",
          "value": {
            "connect_timeout": "10s",
            "http2_protocol_options": {},
            "lb_policy": "ROUND_ROBIN",
            "load_assignment": {
              "cluster_name": "rate_limit_cluster",
              "endpoints": [
                {
                  "lb_endpoints": [
                    {
                      "endpoint": {
                        "address": {
                          "socket_address": {
                            "address": "ratelimit.default",
                            "port_value": 8081
                          }
                        }
                      }
                    }
                  ]
                }
              ]
            },
            "name": "rate_limit_cluster",
            "type": "STRICT_DNS"
          }
        }
      }
    ],
    "workloadSelector": {
      "labels": {
        "app": "%s"
      }
    }
  }
}
`

func (r *globalRateLimitFilter) PrepareUpdateEnvoyFilterExternalObjects(ctx context.Context, instance *v1beta1.GlobalRateLimit) {
	var err error

	name := instance.Name
	namespace := instance.Namespace

	patchValue, envoyFilterObj, err := getGlobalRateLimitEnvoyFilter(namespace, instance)
	efCustomName := name + RlEnvoyFilterSuffixName

	_, err = r.istio.GetEnvoyFilter(ctx, namespace, efCustomName)
	if err != nil {
		klog.Infof("Envoyfilter %s is not found. Error %v", efCustomName)
		_, err = r.istio.CreateEnvoyFilter(ctx, namespace, envoyFilterObj)

		if err != nil {
			klog.Infof("EnvoyFilter created %s", efCustomName)

		}
	} else {
		_, err := r.istio.PatchEnvoyFilter(ctx, patchValue, namespace, efCustomName)
		klog.Infof("Patching Envoyfilter %s", efCustomName)

		if err != nil {
			klog.Infof("Cannot path Ratelimit CR %s. Error %v", efCustomName, err)
		}
	}
}

func getGlobalRateLimitEnvoyFilter(namespace string, limit *v1beta1.GlobalRateLimit) ([]byte, *v1alpha3.EnvoyFilter, error) {
	domain := limit.Spec.Domain
	envoyFilter := v1alpha3.EnvoyFilter{}
	printf := fmt.Sprintf(globalRateLimitEnvoyFilter, getEnvoyFilterName(limit.Name), namespace, domain, limit.Spec.Workload)
	byte := bytes.NewBufferString(printf).Bytes()
	err := json.Unmarshal(byte, &envoyFilter)
	if err != nil {
		return nil, nil, err
	}
	return byte, &envoyFilter, nil
}

func getEnvoyFilterName(crdName string) string {
	return crdName + RlEnvoyFilterSuffixName
}
