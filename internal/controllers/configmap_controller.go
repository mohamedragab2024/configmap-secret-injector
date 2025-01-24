package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigMapReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Logger *zerolog.Logger
}

func (r *ConfigMapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	cm := &corev1.ConfigMap{}

	err := r.Get(ctx, req.NamespacedName, cm)
	if err != nil {
		return ctrl.Result{}, err
	}

	return r.processConfigMapInjection(ctx, cm)
}

func (r *ConfigMapReconciler) processConfigMapInjection(ctx context.Context, cm *corev1.ConfigMap) (ctrl.Result, error) {
	if cm.ObjectMeta.Annotations["secret-injector/enabled"] != "true" || cm.ObjectMeta.Annotations["secret-injector/secret-name"] == "" {
		return ctrl.Result{}, nil
	}

	r.Logger.Info().Msgf("ConfigMap %s/%s is enabled for secret injection", cm.Namespace, cm.Name)
	s := &corev1.Secret{}
	err := r.Get(ctx, client.ObjectKey{
		Namespace: cm.Namespace,
		Name:      cm.ObjectMeta.Annotations["secret-injector/secret-name"],
	}, s)
	if err != nil {
		return ctrl.Result{}, err
	}
	injected := false
	cm.Data, injected = configMapSubstitute(cm.Data, s.Data)
	if injected {
		cm.ObjectMeta.Annotations["secret-injector/last-inject-date"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")
		if err := r.Update(ctx, cm); err != nil {
			return ctrl.Result{}, err
		}
		r.Logger.Info().Msgf("ConfigMap %s/%s has been updated", cm.Namespace, cm.Name)
	}
	return ctrl.Result{}, nil

}
func (r *ConfigMapReconciler) New(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		Complete(r)
}

func configMapSubstitute(cm map[string]string, s map[string][]byte) (map[string]string, bool) {
	injected := false
	cmdata := make(map[string]string)
	for k, v := range cm {
		for sk, sv := range s {
			if strings.Contains(v, fmt.Sprintf("${%s}", sk)) {
				cmdata[k] = strings.ReplaceAll(v, fmt.Sprintf("${%s}", sk), strings.TrimSpace(string(sv)))
				injected = true
			} else {
				cmdata[k] = v
			}
		}
	}
	return cmdata, injected
}
