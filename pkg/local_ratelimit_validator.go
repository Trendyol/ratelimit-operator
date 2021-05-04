package pkg

import (
	"errors"
	trendyolcomv1beta1 "gitlab.trendyol.com/platform/base/apps/ratelimit-operator/api/v1beta1"
)

var (
	workloadEmptyError = "Workload can not be empty"
	tokenBucketEmptyError = "TokenBucket can not be empty"
)

func Validate(localRateLimit *trendyolcomv1beta1.LocalRateLimit) error {

	if len(localRateLimit.Spec.Workload) == 0 {
		return errors.New(workloadEmptyError)
	}

	if localRateLimit.Spec.TokenBucket == (trendyolcomv1beta1.TokenBucket{}) {
		return errors.New(tokenBucketEmptyError)
	}

	return nil
}
