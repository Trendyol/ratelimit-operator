package local

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/api/v1beta1"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/pkg/client/istio"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var rLocalEnvoyFilterSuffixName = "-local-ratelimit"

var localRateLimitEf = `{
    "apiVersion": "networking.istio.io/v1alpha3",
    "kind": "EnvoyFilter",
    "metadata": {
        "name": "%s",
        "namespace": "%s"
    },
    "spec": {
        "configPatches": [
            {
                "applyTo": "HTTP_FILTER",
                "listener": {
                    "filterChain": {
                        "filter": {
                            "name": "envoy.filters.network.http_connection_manager"
                        }
                    }
                },
                "patch": {
                    "operation": "INSERT_BEFORE",
                    "value": {
                        "name": "envoy.filters.http.local_ratelimit",
                        "typed_config": {
                            "@type": "type.googleapis.com/udpa.type.v1.TypedStruct",
                            "type_url": "type.googleapis.com/envoy.extensions.filters.http.local_ratelimit.v3.LocalRateLimit",
                            "value": {
                                "filter_enabled": {
                                    "default_value": {
                                        "denominator": "HUNDRED",
                                        "numerator": 100
                                    },
                                    "runtime_key": "local_rate_limit_enabled"
                                },
                                "filter_enforced": {
                                    "default_value": {
                                        "denominator": "HUNDRED",
                                        "numerator": 100
                                    },
                                    "runtime_key": "local_rate_limit_enforced"
                                },
                                "response_headers_to_add": [
                                    {
                                        "append": false,
                                        "header": {
                                            "key": "x-local-rate-limit",
                                            "value": "true"
                                        }
                                    }
                                ],
                                "stat_prefix": "http_local_rate_limiter",
                                "token_bucket": {
                                    "fill_interval": "%s",
                                    "max_tokens": %d,
                                    "tokens_per_fill": %d
                                }
                            }
                        }
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
}`

type localRateLimit struct {
	client client.Client
	istio  istio.Istio
}

type LocalRateLimit interface {
	DecommissionResources(ctx context.Context, name, namespace string) error
	PrepareUpdateEnvoyFilterObjects(ctx context.Context, global *v1beta1.LocalRateLimit, name, namespace string)
}

func NewLocalRateLimit(client client.Client, istioClient istio.Istio) LocalRateLimit {
	return &localRateLimit{
		client: client,
		istio:  istioClient,
	}
}

func (r *localRateLimit) DecommissionResources(ctx context.Context, name, namespace string) error {
	return r.istio.DeleteEnvoyFilter(ctx, namespace, getEnvoyFilterName(name))
}

func (r *localRateLimit) PrepareUpdateEnvoyFilterObjects(ctx context.Context, instance *v1beta1.LocalRateLimit, name, namespace string) {
	var err error

	patch, envoyFilter, _ := getLocalRateLimitEnvoyFilter(namespace, instance)

	_ , err = r.istio.GetEnvoyFilter(ctx, namespace, getEnvoyFilterName(name))

	if err != nil {
		klog.Infof("Envoyfilter %s is not found. Error %v", instance.Name, err)
		_, err = r.istio.CreateEnvoyFilter(ctx, namespace, envoyFilter)
		klog.Infof("Creating Envoyfilter %s", instance.Name)

		if err != nil {
			klog.Infof("Cannot create Ratelimit CR %s. Error %v", instance.Name, err)
		}
	} else {
		_, err := r.istio.PatchEnvoyFilter(ctx, patch, namespace, envoyFilter.Name)
		klog.Infof("Patching Envoyfilter %s", envoyFilter.Name)

		if err != nil {
			klog.Infof("Cannot path Ratelimit CR %s. Error %v", instance.Name, err)
		}
	}
}
func getLocalRateLimitEnvoyFilter(namespace string, limit *v1beta1.LocalRateLimit) ([]byte, *v1alpha3.EnvoyFilter, error) {
	tokenBucket := limit.Spec.TokenBucket
	envoyFilter := v1alpha3.EnvoyFilter{}
	printf := fmt.Sprintf(localRateLimitEf, getEnvoyFilterName(limit.Name), namespace, tokenBucket.FillInterval, tokenBucket.MaxTokens, tokenBucket.TokensPerFill, limit.Spec.Workload)
	byte := bytes.NewBufferString(printf).Bytes()
	err := json.Unmarshal(byte, &envoyFilter)
	if err != nil {
		return nil, nil, err
	}
	return byte, &envoyFilter, nil
}

func getEnvoyFilterName(crdName string) string {
	return crdName + rLocalEnvoyFilterSuffixName
}
