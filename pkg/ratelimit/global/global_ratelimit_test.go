package global

import (
	"github.com/stretchr/testify/assert"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/api/v1beta1"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/pkg/client/istio"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func Test_CreateOrUpdateConfigMap_When_GenericKey(t *testing.T) {
	globalRateLimit := NewGlobalRateLimit(fake.NewFakeClient(), istio.FakeClient())
	globalRateLimit.InitResources()
	global := &v1beta1.GlobalRateLimit{
		Spec: v1beta1.GlobalRateLimitSpec{
			Domain:   "domain",
			Workload: "workload",
			Rate: []v1beta1.Rate{
				{
					Unit:           "1m",
					RequestPerUnit: 10,
					Dimensions: []v1beta1.Dimensions{
						{
							GenericKey: &v1beta1.GenericKey{
								DescriptorValue: "value",
								DescriptorKey:   "key",
							},
						},
					},
				},
			},
		},
	}

	_ = globalRateLimit.createOrUpdateConfigMap(global)
	configMap, _ := globalRateLimit.getConfigMap(RlConfigMapName, RlConfigMapNameSpace)

	assert.Equal(t, "descriptors:\n- key: key\n  rate_limit:\n    requests_per_unit: 10\n    unit: 1m\n  value: value\ndomain: domain\n", configMap.Data["config..yaml"])
}
