package global

import (
	"context"
	"encoding/json"
	"fmt"
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

var RlConfigMapName = "ratelimit-configmap"
var RlConfigMapNameSpace = "ratelimit"
var EnvoyFilterNames = []string{
	"%s-ratelimit-action",
	"%s-ratelimit-filter",
}

type GlobalRateLimit struct {
	RlAction
	RlFilter
	client client.Client
	mutex  sync.RWMutex
	istio  istio.Istio
}

func NewGlobalRateLimit(client client.Client, istioClient istio.Istio) *GlobalRateLimit {
	return &GlobalRateLimit{
		RlAction: NewGlobalRateLimitAction(istioClient),
		RlFilter: NewGlobalRateLimitFilter(istioClient),
		client:   client,
		istio:    istioClient,
	}
}

func (r *GlobalRateLimit) DecommissionResources(ctx context.Context, instance *v1beta1.GlobalRateLimit) {
	for _, each := range EnvoyFilterNames {
		filterName := fmt.Sprintf(each, instance.Name)
		err := r.DeleteEnvoyFilterObjects(ctx, filterName, instance.Namespace)
		if err != nil {
			klog.Infof("Error when delete envoyfilter %s", filterName)
		}
	}

	found := v1.ConfigMap{}

	found, err := r.getConfigMap(RlConfigMapName, RlConfigMapNameSpace)
	if err != nil {
		klog.Infof("Configmap not found name: %s", RlConfigMapName)

	} else {
		configMapKey := "config." + instance.Name + ".yaml"

		if len(found.Data) > 0 {
			delete(found.Data, configMapKey)
		}

		applyOpts := []client.UpdateOption{client.FieldOwner("globalratelimit-controller")}

		err := r.client.Update(context.TODO(), &found, applyOpts...)

		if err != nil {
			klog.Infof("Error update configmap domain name  %s ", configMapKey)
		}
	}
}

func (r *GlobalRateLimit) CreateOrUpdateResources(ctx context.Context, global *v1beta1.GlobalRateLimit) {
	r.PrepareUpdateEnvoyFilterObjects(ctx, global)
	err := r.CreateOrUpdateConfigMap(global)
	if err != nil {
		return
	}
}

func (r *GlobalRateLimit) PrepareUpdateEnvoyFilterObjects(ctx context.Context, global *v1beta1.GlobalRateLimit) {
	r.PrepareUpdateEnvoyFilterActionObjects(ctx, global)
	r.PrepareUpdateEnvoyFilterExternalObjects(ctx, global)
}

func (r *GlobalRateLimit) DeleteEnvoyFilterObjects(ctx context.Context, name, namespace string) error {
	return r.istio.DeleteEnvoyFilter(ctx, namespace, name)
}

func (r *GlobalRateLimit) CreateOrUpdateConfigMap(global *v1beta1.GlobalRateLimit) error {
	var err error
	name := global.Name
	cmData, err := prepareConfigMapData(name, global)
	if err != nil {
		klog.Infof("Cannot generate %v, Error: %v", cmData, err)
		return err
	}

	found := v1.ConfigMap{}
	var configMapData = make(map[string]string, 0)
	found, err = r.getConfigMap(RlConfigMapName, RlConfigMapNameSpace)

	if err == nil {
		if found.Data == nil {
			found.Data = configMapData
		}
		configMapKey := "config." + name + ".yaml"
		found.Data[configMapKey] = cmData[configMapKey]

		applyOpts := []client.UpdateOption{client.FieldOwner("ratelimit-controller")}

		r.client.Update(context.TODO(), &found, applyOpts...)
		r.mutex.Lock()
		defer r.mutex.Unlock()

	}

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

func (r *GlobalRateLimit) InitResources() v1.ConfigMap {
	name := RlConfigMapName
	namespace := RlConfigMapNameSpace
	cm := v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    map[string]string{"generator": "ratelimit-operator"},
		},
		Data: make(map[string]string, 0),
	}

	//TODO: os.GetEnv || default name,namespace
	foundCm, err := r.getConfigMap(name, namespace)
	if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonNotFound {
		err = r.client.Create(context.TODO(), &cm)
		return cm
	} else {
		return foundCm
	}
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
			if eachDimension.RequestHeader != nil {
				cfDescriptor.Key = eachDimension.RequestHeader.DescriptorKey
				cfDescriptor.Value = eachDimension.RequestHeader.Value
			}
			if eachDimension.HeaderValueMatch != nil {
				cfDescriptor.Key = "header_match"
				cfDescriptor.Value = eachDimension.HeaderValueMatch.DescriptorValue
			}

			if eachDimension.RemoteAddress != nil {
				cfDescriptor.Key = "remote_address"
			}

			if eachDimension.GenericKey != nil {
				cfDescriptor.Key = eachDimension.GenericKey.DescriptorKey
				cfDescriptor.Value = eachDimension.GenericKey.DescriptorValue
			}

			if eachDimension.SourceCluster != nil {
				cfDescriptor.Key = "source_cluster"
			}

			if eachDimension.DestinationCluster != nil {
				cfDescriptor.Key = "destination_cluster"
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

//TODO:Nested descriptor
type configMapValue struct {
	Domain      string                `json:"domain"`
	Descriptors []configMapDescriptor `json:"descriptors"`
}

type configMapDescriptor struct {
	Key       string                 `json:"key"`
	Value     string                 `json:"value,omitempty"`
	RateLimit configMapRatelimitUnit `json:"rate_limit"`
}

type configMapRatelimitUnit struct {
	Unit            string `json:"unit"`
	RequestsPerUnit int64  `json:"requests_per_unit"`
}
