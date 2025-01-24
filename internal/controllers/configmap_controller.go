package controllers

import (
	"context"
	"fmt"
	"strings"

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

	if cm.ObjectMeta.Annotations["secret-injector/enabled"] == "true" && cm.ObjectMeta.Annotations["secret-injector/secret-name"] != "" {
		r.Logger.Info().Msgf("ConfigMap %s/%s is enabled for secret injection", cm.Namespace, cm.Name)
		sr := &corev1.Secret{}
		err := r.Get(ctx, client.ObjectKey{Namespace: cm.Namespace, Name: cm.ObjectMeta.Annotations["secret-injector/secret-name"]}, sr)
		if err != nil {
			return ctrl.Result{}, err
		}
		cmdata := configMapSubstitute(cm.Data, sr.Data)
		cm.Data = cmdata
		uerr := r.Update(ctx, cm)
		if uerr != nil {
			return ctrl.Result{}, uerr
		}

		r.Logger.Info().Msgf("ConfigMap %s/%s has been updated", cm.Namespace, cm.Name)
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil

}
func (r *ConfigMapReconciler) New(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		Complete(r)
}

func configMapSubstitute(cm map[string]string, sr map[string][]byte) map[string]string {
	cmdata := make(map[string]string)
	for k, v := range cm {
		for sk, sv := range sr {
			if strings.Contains(v, fmt.Sprintf("{%s}", sk)) {
				cmdata[k] = strings.ReplaceAll(v, fmt.Sprintf("{%s}", sk), strings.TrimSpace(string(sv)))
			} else {
				cmdata[k] = v
			}
		}
	}
	return cmdata
}
