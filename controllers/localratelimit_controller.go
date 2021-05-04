/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	trendyolcomv1beta1 "gitlab.trendyol.com/platform/base/apps/ratelimit-operator/api/v1beta1"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/client/istio"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/pkg"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// LocalRateLimitReconciler reconciles a LocalRateLimit object
type LocalRateLimitReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	IstioClient istio.IstioClient
}

//+kubebuilder:rbac:groups=trendyol.com,resources=localratelimits,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=trendyol.com,resources=localratelimits/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=trendyol.com,resources=localratelimits/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LocalRateLimit object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *LocalRateLimitReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	namespace := req.Namespace
	_ = r.Log.WithValues("localratelimit", req.NamespacedName)
	localRateLimitInstance := &trendyolcomv1beta1.LocalRateLimit{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      req.Name,
	}, localRateLimitInstance)
	localEnvoyFilterName := localRateLimitInstance.Spec.Workload + "-local-ratelimit"

	if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonNotFound {
		r.IstioClient.DeleteEnvoyFilter(ctx, namespace, localEnvoyFilterName)
	}


	if err != nil {
		klog.Infof("Cannot get LocalRatelimit CR %s. Error %v", localRateLimitInstance.Name, err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	beingDeleted := localRateLimitInstance.GetDeletionTimestamp() != nil

	if beingDeleted {
		r.IstioClient.DeleteEnvoyFilter(ctx, namespace, localEnvoyFilterName)
	}

	err = pkg.Validate(localRateLimitInstance)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	byte, envoyFilter, err := istio.GetLocalRateLimitEnvoyFilter(namespace, localRateLimitInstance)

	if err != nil {
		//klog.Infof("Cannot get Ratelimit CR %s. Error %v", rateLimitInstance.Name, err)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	_, err = r.IstioClient.CreateEnvoyFilter(ctx, namespace, envoyFilter)
	if err != nil {
		_, err := r.IstioClient.PatchEnvoyFilter(ctx, byte, namespace, localEnvoyFilterName)
		return ctrl.Result{}, client.IgnoreNotFound(err)

	}
	// your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LocalRateLimitReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&trendyolcomv1beta1.LocalRateLimit{}).
		Complete(r)
}
