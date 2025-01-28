package controllers

import (
	"context"
	"fmt"
	"strings"
	"sync"
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
	mu     sync.Mutex
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
		r.mu.Lock()
		defer r.mu.Unlock()
		cm.ObjectMeta.Annotations["secret-injector/last-inject-date"] = time.Now().UTC().Format("0000-00-00T00:00:00Z")
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
		originalValue := v
		for sk, sv := range s {
			placeholder := fmt.Sprintf("${%s}", sk)
			if strings.Contains(v, placeholder) {
				v = strings.ReplaceAll(v, placeholder, strings.TrimSpace(string(sv)))
			}
		}
		cmdata[k] = v
		if v != originalValue {
			injected = true
		}
	}
	return cmdata, injected
}
