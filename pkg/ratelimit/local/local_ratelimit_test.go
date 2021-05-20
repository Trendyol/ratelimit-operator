package local

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/api/v1beta1"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/pkg/client/istio"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func Test_ValidateLocalEnvoyFilterJson(t *testing.T) {

	local := &v1beta1.LocalRateLimit{
		Spec: v1beta1.LocalRateLimitSpec{
			Workload: "workload",
			TokenBucket: v1beta1.TokenBucket{
				MaxTokens:     100,
				TokensPerFill: 50,
				FillInterval:  "60s",
			},
		},
	}
	_, envoyFilterObj, err := getLocalRateLimitEnvoyFilter("default", local)
	assert.Nil(t, err)

	workload := envoyFilterObj.Spec.WorkloadSelector.Labels["app"]
	assert.Equal(t, workload, "workload")

	value := envoyFilterObj.Spec.ConfigPatches[0].Patch.Value.GetFields()["typed_config"].GetStructValue().GetFields()["value"].GetStructValue().GetFields()["token_bucket"].GetStructValue()
	fillIntervalValue := value.GetFields()["fill_interval"].GetStringValue()
	maxTokens := value.GetFields()["max_tokens"].GetNumberValue()
	TokenPerFill := value.GetFields()["tokens_per_fill"].GetNumberValue()
	assert.Equal(t, fillIntervalValue, "60s")
	assert.Equal(t, maxTokens, float64(100))
	assert.Equal(t, TokenPerFill, float64(50))
}

func Test_CreateEnvoyFilterObjectWhenNoExistingEnvoyFilter(t *testing.T) {

	localRateLimit := &localRateLimit{istio: FakeClient()}

	local := &v1beta1.LocalRateLimit{
		Spec: v1beta1.LocalRateLimitSpec{
			Workload: "workload",
			TokenBucket: v1beta1.TokenBucket{
				MaxTokens:     100,
				TokensPerFill: 50,
				FillInterval:  "60s",
			},
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "localrate",
			Namespace: "default",
		},
	}

	localRateLimit.PrepareUpdateEnvoyFilterObjects(context.Background(), local, "localrate", "default")

	filter, _ := localRateLimit.istio.GetEnvoyFilter(context.Background(), "default", "localrate-local-ratelimit")

	assert.Equal(t, "localrate-local-ratelimit", filter.Name)
}

func FakeClient() istio.Istio {
	return istio.FakeClient()
}
