package local

import (
	"github.com/stretchr/testify/assert"
	trendyolcomv1beta1 "gitlab.trendyol.com/platform/base/apps/ratelimit-operator/api/v1beta1"
	"testing"
)

func Test_ValidateIsWorkloadEmpty(t *testing.T) {
	 err := Validate(&trendyolcomv1beta1.LocalRateLimit{
	 	Spec: trendyolcomv1beta1.LocalRateLimitSpec{
	 	Workload: "app",
	 	TokenBucket: trendyolcomv1beta1.TokenBucket{
	 	TokensPerFill: 100,
		},
		},
	 })

	if err != nil {
		assert.Error(t,err)
	}
}
