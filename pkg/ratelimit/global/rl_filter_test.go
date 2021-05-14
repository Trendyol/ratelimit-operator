package global

import (
	"github.com/stretchr/testify/assert"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/api/v1beta1"
	"testing"
)

func Test_ValidateEnvoyFilterJson(t *testing.T) {

	global := &v1beta1.GlobalRateLimit{
		Spec: v1beta1.GlobalRateLimitSpec{
			Domain:   "domain",
			Workload: "workload",
			Rate: []v1beta1.Rate{
				{
					Unit: "1m",
				},
			},
		},
	}
	_, envoyFilterObj, err := GetGlobalRateLimitEnvoyFilter("default", global)
	assert.Nil(t, err)

	workload := envoyFilterObj.Spec.WorkloadSelector.Labels["app"]
	assert.Equal(t, workload, "workload")

	domain := envoyFilterObj.Spec.ConfigPatches[0].Patch.Value.GetFields()["typed_config"].GetStructValue().GetFields()["domain"].GetStringValue()
	assert.Equal(t, domain, "domain")
}
