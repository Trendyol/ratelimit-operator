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
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/pkg/ratelimit/global"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	trendyolcomv1beta1 "gitlab.trendyol.com/platform/base/apps/ratelimit-operator/api/v1beta1"
)

// GlobalRateLimitReconciler reconciles a GlobalRateLimit object
type GlobalRateLimitReconciler struct {
	client.Client
	Log             logr.Logger
	Scheme          *runtime.Scheme
	GlobalRateLimit *global.GlobalRateLimit
}

//+kubebuilder:rbac:groups=trendyol.com,resources=globalratelimits,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=trendyol.com,resources=globalratelimits/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=trendyol.com,resources=globalratelimits/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GlobalRateLimit object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *GlobalRateLimitReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("globalratelimit", req.NamespacedName)

	name := req.Name
	namespace := req.Namespace
	globalRateLimitInstance := &trendyolcomv1beta1.GlobalRateLimit{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      req.Name,
	}, globalRateLimitInstance)
	//If resource deleted set from request
	if len(globalRateLimitInstance.Name) == 0 {
		globalRateLimitInstance.Name = name
		globalRateLimitInstance.Namespace = namespace
	}

	//Init k8s resources
	r.GlobalRateLimit.InitResources()

	if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonNotFound {
		r.GlobalRateLimit.DecommissionResources(ctx, globalRateLimitInstance)
		return ctrl.Result{}, nil
	}

	r.GlobalRateLimit.CreateOrUpdateResources(ctx, globalRateLimitInstance)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GlobalRateLimitReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&trendyolcomv1beta1.GlobalRateLimit{}).
		Complete(r)
}
