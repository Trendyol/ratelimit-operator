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

var rLocalEnvoyFilterSuffixName = "-ratelimit-action"

type RlAction interface {
	PrepareUpdateEnvoyFilterActionObjects(ctx context.Context, global *v1beta1.GlobalRateLimit)
}
type globalRateLimitAction struct {
	istio istio.Istio
}

func NewGlobalRateLimitAction(istio istio.Istio) RlAction {
	return &globalRateLimitAction{istio: istio}
}

var globalEnvoyFilterAction = `{
  "kind": "EnvoyFilter",
  "apiVersion": "networking.istio.io/v1alpha3",
  "metadata": {
    "name": "%s-ratelimit-action",
    "namespace": "%s",
    "labels":{
     "generator": "ratelimit-operator"
   }
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

func (r *globalRateLimitAction) PrepareUpdateEnvoyFilterActionObjects(ctx context.Context, instance *v1beta1.GlobalRateLimit) {
	rlAction := &v1beta1.RateLimitAction{}
	var err error
	var rateLimits []v1beta1.RateLimits

	for _, eachRate := range instance.Spec.Rate {
		var actions []v1beta1.Actions
		var rateLimit v1beta1.RateLimits
		for _, dimension := range eachRate.Dimensions {
			action := v1beta1.Actions{}
			//TODO: refactor here
			if dimension.HeaderValueMatch != nil {
				action.HeaderValueMatch = dimension.HeaderValueMatch
			}
			if dimension.RequestHeader != nil {
				action.RequestHeader = dimension.RequestHeader
			}
			if dimension.SourceCluster != nil {
				action.SourceCluster = dimension.SourceCluster
			}
			if dimension.DestinationCluster != nil {
				action.DestinationCluster = dimension.DestinationCluster
			}
			if dimension.GenericKey != nil {
				action.GenericKey = dimension.GenericKey
			}
			action.RemoteAddress = dimension.RemoteAddress
			actions = append(actions, action)
		}
		rateLimit.Actions = actions

		rateLimits = append(rateLimits, rateLimit)
	}

	rlAction.RateLimits = rateLimits
	strRlAction, _ := json.Marshal(&rlAction)
	pretty, _ := prettyPrint(strRlAction)

	namespace := instance.Namespace
	name := getActionEnvoyFilterName(instance.Name)
	patchValue, envoyFilterObj, err := getGlobalEnvoyFilterAction(namespace, string(pretty), instance)

	_, err = r.istio.GetEnvoyFilter(ctx, namespace, name)
	if err != nil {
		_, err = r.istio.CreateEnvoyFilter(ctx, namespace, envoyFilterObj)
		klog.Infof("Creating Envoyfilter %s", name)

		if err != nil {
			klog.Infof("Cannot get Ratelimit CR %s. Error %v", instance.Name, err)
		}
	} else {
		_, err := r.istio.PatchEnvoyFilter(ctx, patchValue, namespace, name)
		klog.Infof("Patching Envoyfilter %s", name)

		if err != nil {
			klog.Infof("Cannot path Ratelimit CR %s. Error %v", instance.Name, err)
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

func prettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

func getActionEnvoyFilterName(crdName string) string {
	return crdName + rLocalEnvoyFilterSuffixName
}
