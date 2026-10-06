package controller

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	configv1alpha1 "github.com/Techbrainzz/kubernetes-config-reloader/api/v1alpha1"
)

type ConfigSyncReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=config.infra.io,resources=configsyncs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=config.infra.io,resources=configsyncs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=config.infra.io,resources=configsyncs/finalizers,verbs=get;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update;patch

func (r *ConfigSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling ConfigSync", "name", req.NamespacedName)

	var configSync configv1alpha1.ConfigSync
	if err := r.Get(ctx, req.NamespacedName, &configSync); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Helper function to update status
	updateStatus := func(phase, message string) error {
		configSync.Status.Phase = phase
		configSync.Status.Message = message
		return r.Status().Update(ctx, &configSync)
	}

	// 1. Sync ConfigMap
	desiredConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configSync.Name + "-synced",
			Namespace: configSync.Namespace,
		},
		Data: configSync.Spec.ConfigData,
	}

	if err := ctrl.SetControllerReference(&configSync, desiredConfigMap, r.Scheme); err != nil {
		_ = updateStatus("Failed", err.Error())
		return ctrl.Result{}, err
	}

	var existingConfigMap corev1.ConfigMap
	err := r.Get(ctx, types.NamespacedName{Name: desiredConfigMap.Name, Namespace: desiredConfigMap.Namespace}, &existingConfigMap)
	if err != nil && apierrors.IsNotFound(err) {
		logger.Info("Creating target ConfigMap", "name", desiredConfigMap.Name)
		if err := r.Create(ctx, desiredConfigMap); err != nil {
			_ = updateStatus("Failed", err.Error())
			return ctrl.Result{}, err
		}
	} else if err == nil {
		existingConfigMap.Data = desiredConfigMap.Data
		if err := r.Update(ctx, &existingConfigMap); err != nil {
			_ = updateStatus("Failed", err.Error())
			return ctrl.Result{}, err
		}
	} else {
		_ = updateStatus("Failed", err.Error())
		return ctrl.Result{}, err
	}

	// 2. Trigger Rolling Restarts
	configHash := calculateConfigHash(configSync.Spec.ConfigData)
	for _, deployName := range configSync.Spec.TargetDeployments {
		var deployment appsv1.Deployment
		deployKey := types.NamespacedName{Name: deployName, Namespace: configSync.Namespace}

		if err := r.Get(ctx, deployKey, &deployment); err != nil {
			logger.Error(err, "Failed to fetch deployment", "deployment", deployName)
			continue
		}

		if deployment.Spec.Template.Annotations == nil {
			deployment.Spec.Template.Annotations = make(map[string]string)
		}

		if deployment.Spec.Template.Annotations["config.infra.io/config-hash"] != configHash {
			deployment.Spec.Template.Annotations["config.infra.io/config-hash"] = configHash
			logger.Info("Triggering rolling restart", "deployment", deployName)
			if err := r.Update(ctx, &deployment); err != nil {
				_ = updateStatus("Failed", err.Error())
				return ctrl.Result{}, err
			}
		}
	}

	// 3. Mark as successfully synced
	if err := updateStatus("Synced", "Successfully synchronized config and updated deployments at "+time.Now().Format(time.RFC3339)); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func calculateConfigHash(data map[string]string) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write([]byte(data[k]))
	}
	return hex.EncodeToString(h.Sum(nil))[:12]
}

func (r *ConfigSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&configv1alpha1.ConfigSync{}).
		Complete(r)
}