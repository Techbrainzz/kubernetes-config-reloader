/*
Copyright 2026.

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

package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	configv1alpha1 "github.com/Techbrainzz/kubernetes-config-reloader/api/v1alpha1"
)

// ConfigSyncReconciler reconciles a ConfigSync object
type ConfigSyncReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=config.infra.io,resources=configsyncs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=config.infra.io,resources=configsyncs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=config.infra.io,resources=configsyncs/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

func (r *ConfigSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var configSync configv1alpha1.ConfigSync
	if err := r.Get(ctx, req.NamespacedName, &configSync); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	logger.Info("Reconciling ConfigSync", "Namespace", configSync.Namespace, "Name", configSync.Name)

	desiredConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configSync.Name + "-sync",
			Namespace: configSync.Spec.TargetNamespace,
		},
		Data: configSync.Spec.ConfigData,
	}

	if err := controllerutil.SetControllerReference(&configSync, desiredConfigMap, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	var foundConfigMap corev1.ConfigMap
	err := r.Get(ctx, types.NamespacedName{Name: desiredConfigMap.Name, Namespace: desiredConfigMap.Namespace}, &foundConfigMap)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating target ConfigMap", "Namespace", desiredConfigMap.Namespace, "Name", desiredConfigMap.Name)
		if err := r.Create(ctx, desiredConfigMap); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		foundConfigMap.Data = desiredConfigMap.Data
		logger.Info("Updating target ConfigMap", "Namespace", foundConfigMap.Namespace, "Name", foundConfigMap.Name)
		if err := r.Update(ctx, &foundConfigMap); err != nil {
			return ctrl.Result{}, err
		}
	}

	configSync.Status.Phase = "Synced"
	configSync.Status.Message = "ConfigMap successfully synchronized by operator"
	if err := r.Status().Update(ctx, &configSync); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ConfigSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&configv1alpha1.ConfigSync{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}