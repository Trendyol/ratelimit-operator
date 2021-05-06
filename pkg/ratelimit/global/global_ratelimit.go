package global

import (
	"context"
	"encoding/json"
	"github.com/ghodss/yaml"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/api/v1beta1"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/pkg/client/istio"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

type GlobalRateLimit struct {
	*GlobalRateLimitAction
	*GlobalRateLimitFilter
	client client.Client
	mutex  sync.RWMutex
	istio  istio.IstioClient
}

func NewGlobalRateLimit(client client.Client, istioClient istio.IstioClient) *GlobalRateLimit {
	return &GlobalRateLimit{
		GlobalRateLimitAction: NewGlobalRateLimitAction(istioClient),
		GlobalRateLimitFilter: NewGlobalRateLimitFilter(istioClient),
		client:                client,
		istio:                 istioClient,
	}
}

func (r *GlobalRateLimit) DecommissionResources(ctx context.Context, name, namespace string) {

	//TODO: Delete Envoy filters
	// Patch Configmap
}

func (r *GlobalRateLimit) CreateOrUpdateResources(ctx context.Context, global *v1beta1.GlobalRateLimit, name, namespace string) {
	r.PrepareUpdateEnvoyFilterObjects(ctx, global, name, namespace)
	err := r.CreateOrUpdateConfigMap(global, name, namespace)
	if err != nil {
		return
	}
}

func (r *GlobalRateLimit) PrepareUpdateEnvoyFilterObjects(ctx context.Context, global *v1beta1.GlobalRateLimit, name, namespace string) {
	r.PrepareUpdateEnvoyFilterActionObjects(ctx, global, namespace, name)
	r.PrepareUpdateEnvoyFilterExternalObjects(ctx, global, namespace, name)
}

func getEmptyConfigMap() v1.ConfigMap {
	//TODO
	return v1.ConfigMap{}
}

func (r *GlobalRateLimit) CreateOrUpdateConfigMap(global *v1beta1.GlobalRateLimit, name, namespace string) error {
	var err error
	cmData, err := prepareConfigMapData(name, global)
	if err != nil {
		klog.Infof("Cannot generate %v, Error: %v", cmData, err)
		return err
	}

	found := v1.ConfigMap{}

	found, err = r.getConfigMap("ratelimit-configmap2", namespace)
	cm := v1.ConfigMap{}
	if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonNotFound {
		err = r.client.Create(context.TODO(), &cm)
		return nil
	} else {
		configMapKey := "config." + name + ".yaml"
		found.Data[configMapKey] = cmData[configMapKey]
	}

	applyOpts := []client.UpdateOption{client.FieldOwner("globalratelimit-controller")}

	r.client.Update(context.TODO(), &found, applyOpts...)
	r.mutex.Lock()
	defer r.mutex.Unlock()
	//if err != nil {
	//	err = r.client.Create(context.TODO(), &cm)
	//	if err != nil {
	//		//klog.Infof("Cannot create %v, Error: %v", cm, err.Error())
	//	}
	//} else {
	//	applyOpts := []client.PatchOption{client.ForceOwnership, client.FieldOwner("globalratelimit-controller")}
	//
	//	r.client.Update(context.TODO(), &cm, applyOpts...)
	//
	//	klog.Infof("the 2 resources are not the same, we should patch the deployment")
	//
	//	err = r.client.Patch(context.TODO(), &cm, client.Apply, applyOpts...)
	//	if err != nil {
	//		klog.Infof("Cannot patch cm. Error: %v", err)
	//		return err
	//	}
	//
	//	r.mutex.Lock()
	//	defer r.mutex.Unlock()
	//}

	return nil
}

func (r *GlobalRateLimit) getConfigMap(name string, namespace string) (v1.ConfigMap, error) {

	found := v1.ConfigMap{}

	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, &found)
	if err != nil {
		//klog.Infof("Cannot Found configMap %s. Error %v", found.Name, err)
		return found, err
	}

	return found, nil

}

func prepareConfigMapData(name string, global *v1beta1.GlobalRateLimit) (map[string]string, error) {
	configMapData := make(map[string]string)
	var err error

	domain := global.Spec.Domain

	var descriptors []configMapDescriptor
	value := configMapValue{}
	value.Domain = domain
	for _, eachRate := range global.Spec.Rate {
		cfDescriptor := configMapDescriptor{}
		cfDescriptor.RateLimit.RequestsPerUnit = eachRate.RequestPerUnit
		cfDescriptor.RateLimit.Unit = eachRate.Unit
		for _, eachDimension := range eachRate.Dimensions {
			if len(eachDimension.RequestHeader.DescriptorKey) > 0 {
				cfDescriptor.Key = eachDimension.RequestHeader.DescriptorKey
				cfDescriptor.Value = eachDimension.RequestHeader.Value
			}
			if len(eachDimension.HeaderValueMatch.Headers) > 0 {
				cfDescriptor.Key = "header_match"
				cfDescriptor.Value = eachDimension.HeaderValueMatch.DescriptorValue
			}
			descriptors = append(descriptors, cfDescriptor)
		}
	}
	value.Descriptors = descriptors
	configMapKey := "config." + name + ".yaml"
	var output []byte

	output, err = json.Marshal(value)
	if err != nil {
		klog.Infof("Cannot generate configmap as we cannot marshal descriptor value")
		return configMapData, err

	}
	y, err := yaml.JSONToYAML(output)
	if err != nil {
		return configMapData, err
	}
	configMapData[configMapKey] = string(y)

	return configMapData, nil
}

type configMapValue struct {
	Domain      string                `json:"domain"`
	Descriptors []configMapDescriptor `json:"descriptors"`
}

type configMapDescriptor struct {
	Key       string                 `json:"key"`
	Value     string                 `json:"value"`
	RateLimit configMapRatelimitUnit `json:"rate_limit"`
}

type configMapRatelimitUnit struct {
	Unit            string `json:"unit"`
	RequestsPerUnit int64  `json:"requests_per_unit"`
}
