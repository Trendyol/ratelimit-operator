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

type GlobalRateLimitAction struct {
	istio istio.IstioClient
}

func NewGlobalRateLimitAction(istio istio.IstioClient) *GlobalRateLimitAction {
	return &GlobalRateLimitAction{istio: istio}
}

var globalEnvoyFilterAction = `{
  "kind": "EnvoyFilter",
  "apiVersion": "networking.istio.io/v1alpha3",
  "metadata": {
    "name": "%s-ratelimit-actions",
    "namespace": "%s"
  },
  "spec": {
    "configPatches": [
      {
        "applyTo": "VIRTUAL_HOST",
        "match": {
          "context": "SIDECAR_INBOUND"
        },
        "patch": {
          "operation": "MERGE",
          "value": 
           %s
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


func (r *GlobalRateLimitAction) PrepareUpdateEnvoyFilterActionObjects(ctx context.Context, global *v1beta1.GlobalRateLimit, namespace, name string) {
	rlAction := &v1beta1.RateLimitAction{}
	var err error
	var rateLimits []v1beta1.RateLimits

	for _, eachRate := range global.Spec.Rate {
		var actions []v1beta1.Actions
		var rateLimit v1beta1.RateLimits
		for _, dimension := range eachRate.Dimensions {
			action := v1beta1.Actions{}
			//TODO: null or empty check
			action.RequestHeader = dimension.RequestHeader
			action.HeaderValueMatch = dimension.HeaderValueMatch
			action.RemoteAddress = dimension.RemoteAddress
			actions = append(actions, action)
		}
		rateLimit.Actions = actions

		rateLimits = append(rateLimits, rateLimit)
	}
	rlAction.RateLimits = rateLimits
	strRlAction, _ := json.Marshal(&rlAction)
	pretty, _ := prettyprint(strRlAction)
	patchValue, envoyFilterObj, err := getGlobalEnvoyFilterAction(namespace, string(pretty), global)

	_, err = r.istio.GetEnvoyFilter(ctx, namespace, name)
	if err != nil {
		klog.Infof("Envoyfilter %s is not found. Error %v", name, err)
		_, err = r.istio.CreateEnvoyFilter(ctx, namespace, envoyFilterObj)
		klog.Infof("Creating Envoyfilter %s", name)

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
func getGlobalEnvoyFilterAction(namespace string, action string, global *v1beta1.GlobalRateLimit) ([]byte, *v1alpha3.EnvoyFilter, error) {
	envoyFilter := v1alpha3.EnvoyFilter{}

	result := fmt.Sprintf(globalEnvoyFilterAction, global.Name, namespace, action, global.Spec.Workload)
	byte := bytes.NewBufferString(result).Bytes()
	err := json.Unmarshal(byte, &envoyFilter)
	if err != nil {
		return nil, nil, err
	}
	return byte, &envoyFilter, nil
}

func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
