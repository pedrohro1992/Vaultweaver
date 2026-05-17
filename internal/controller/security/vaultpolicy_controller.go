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

package security

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	securityv1alpha1 "github.com/pedrohro1992/Vaultweaver/api/security/v1alpha1"
	"github.com/pedrohro1992/Vaultweaver/internal/vault"
)

const (
	vaultPolicyFinalizer = "security.platform.io/policy-finalizer"
)

// VaultPolicyReconciler reconciles a VaultPolicy object
type VaultPolicyReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	VaultClient vault.Client
	VaultConfig VaultConfig
}

// +kubebuilder:rbac:groups=security.platform.io,resources=vaultpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=security.platform.io,resources=vaultpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=security.platform.io,resources=vaultpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *VaultPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// 1. Fetch the VaultPolicy instance
	var vp securityv1alpha1.VaultPolicy
	if err := r.Get(ctx, req.NamespacedName, &vp); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Handle deletion
	if !vp.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, &vp)
	}

	// 3. Add finalizer if not present
	if !controllerutil.ContainsFinalizer(&vp, vaultPolicyFinalizer) {
		controllerutil.AddFinalizer(&vp, vaultPolicyFinalizer)
		if err := r.Update(ctx, &vp); err != nil {
			return ctrl.Result{}, err
		}
	}

	// 4. Authenticate to Vault
	if err := r.VaultClient.Authenticate(ctx, r.VaultConfig.AuthRole, r.VaultConfig.AuthMount, r.VaultConfig.AuthJWTPath); err != nil {
		log.Error(err, "Failed to authenticate to Vault")
		return r.updateStatus(ctx, &vp, "AuthenticationFailed", err.Error(), metav1.ConditionFalse)
	}

	// 5. Reconcile Vault Policy
	if err := r.reconcileVaultPolicy(ctx, &vp); err != nil {
		log.Error(err, "Failed to reconcile Vault policy")
		return r.updateStatus(ctx, &vp, "ReconciliationFailed", err.Error(), metav1.ConditionFalse)
	}

	// 6. Update Status to Ready
	return r.updateStatus(ctx, &vp, "Reconciled", "Vault policy is synchronized", metav1.ConditionTrue)
}

func (r *VaultPolicyReconciler) reconcileVaultPolicy(ctx context.Context, vp *securityv1alpha1.VaultPolicy) error {
	log := logf.FromContext(ctx)

	actual, err := r.VaultClient.GetPolicy(ctx, vp.Spec.VaultPolicyName)
	if err != nil {
		return err
	}

	desired := strings.TrimSpace(vp.Spec.Policy)
	current := strings.TrimSpace(actual)

	if actual != "" && desired == current {
		log.Info("Vault policy is already up to date", "policy", vp.Spec.VaultPolicyName)
		return nil
	}

	log.Info("Applying Vault policy", "policy", vp.Spec.VaultPolicyName)
	return r.VaultClient.ApplyPolicy(ctx, vp.Spec.VaultPolicyName, desired)
}

func (r *VaultPolicyReconciler) reconcileDelete(ctx context.Context, vp *securityv1alpha1.VaultPolicy) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if controllerutil.ContainsFinalizer(vp, vaultPolicyFinalizer) {
		// Authenticate to Vault
		if err := r.VaultClient.Authenticate(ctx, r.VaultConfig.AuthRole, r.VaultConfig.AuthMount, r.VaultConfig.AuthJWTPath); err != nil {
			log.Error(err, "Failed to authenticate to Vault for deletion")
			return ctrl.Result{}, err
		}

		log.Info("Deleting Vault policy", "policy", vp.Spec.VaultPolicyName)
		if err := r.VaultClient.DeletePolicy(ctx, vp.Spec.VaultPolicyName); err != nil {
			log.Error(err, "Failed to delete Vault policy")
			return ctrl.Result{}, err
		}

		controllerutil.RemoveFinalizer(vp, vaultPolicyFinalizer)
		if err := r.Update(ctx, vp); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *VaultPolicyReconciler) updateStatus(ctx context.Context, vp *securityv1alpha1.VaultPolicy, reason, message string, status metav1.ConditionStatus) (ctrl.Result, error) {
	vp.Status.ObservedGeneration = vp.Generation

	condition := metav1.Condition{
		Type:               conditionTypeReady,
		Status:             status,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: vp.Generation,
	}

	meta.SetStatusCondition(&vp.Status.Conditions, condition)

	if err := r.Status().Update(ctx, vp); err != nil {
		return ctrl.Result{}, err
	}

	if status == metav1.ConditionFalse {
		return ctrl.Result{}, fmt.Errorf("%s: %s", reason, message)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VaultPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1alpha1.VaultPolicy{}).
		Named("security-vaultpolicy").
		Complete(r)
}
