//go:build dev_features

package ratelimit

import (
	ratelimitv1alpha1 "github.com/kyma-project/api-gateway/apis/gateway/ratelimit/v1alpha1"
	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"time"
)

const (
	RateLimiterEnabled = true

	defaultReconciliationPeriod = 3 * time.Minute
)

func Setup(mgr manager.Manager, scheme *runtime.Scheme) error {
	utilruntime.Must(ratelimitv1alpha1.AddToScheme(scheme))
	utilruntime.Must(networkingv1alpha3.AddToScheme(scheme))
	return (&RateLimitReconciler{
		Client:          mgr.GetClient(),
		Scheme:          mgr.GetScheme(),
		ReconcilePeriod: defaultReconciliationPeriod,
	}).SetupWithManager(mgr)
}
