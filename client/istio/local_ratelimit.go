package istio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/api/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
)

var localRateLimit = `{
    "apiVersion": "networking.istio.io/v1alpha3",
    "kind": "EnvoyFilter",
    "metadata": {
        "name": "%s-local-ratelimit-svc",
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

func GetLocalRateLimitEnvoyFilter(namespace string, limit *v1beta1.LocalRateLimit) (*v1alpha3.EnvoyFilter, error) {
	tokenBucket := limit.Spec.TokenBucket
	envoyFilter := v1alpha3.EnvoyFilter{}

	printf := fmt.Sprintf(localRateLimit, limit.Spec.Workload, namespace, tokenBucket.FillInterval, tokenBucket.MaxToken, tokenBucket.TokenPerFill, limit.Spec.Workload)
	err := json.Unmarshal(bytes.NewBufferString(printf).Bytes(), &envoyFilter)
	if err != nil {
		return nil, err
	}
	return &envoyFilter, nil
}
