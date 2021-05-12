package controllers

import (
	"context"
	"github.com/go-logr/logr"
	trendyolcomv1beta1 "gitlab.trendyol.com/platform/base/apps/ratelimit-operator/api/v1beta1"
	"gitlab.trendyol.com/platform/base/apps/ratelimit-operator/pkg/ratelimit/local"
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
	Log    logr.Logger
	Scheme *runtime.Scheme
	Local  local.LocalRateLimit
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
	name := req.Name
	_ = r.Log.WithValues("localratelimit", req.NamespacedName)
	localRateLimitInstance := &trendyolcomv1beta1.LocalRateLimit{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      req.Name,
	}, localRateLimitInstance)

	if statusError, isStatus := err.(*errors.StatusError); isStatus && statusError.Status().Reason == metav1.StatusReasonNotFound {
		err := r.Local.DecommissionResources(ctx, name, namespace)
		if err != nil {
			klog.Infof("Cannot delete localRatelimit CR %s. Error %v", localRateLimitInstance.Name, err)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		return ctrl.Result{}, nil
	}

	r.Local.PrepareUpdateEnvoyFilterObjects(ctx, localRateLimitInstance, name, namespace)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LocalRateLimitReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&trendyolcomv1beta1.LocalRateLimit{}).
		Complete(r)
}
