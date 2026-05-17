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
	"reflect"

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
	vaultKubernetesRoleBindingFinalizer = "security.platform.io/finalizer"
	conditionTypeReady                  = "Ready"
)

// VaultKubernetesRoleBindingReconciler reconciles a VaultKubernetesRoleBinding object
type VaultKubernetesRoleBindingReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	VaultClient vault.Client
	VaultConfig VaultConfig
}

// VaultConfig contains the configuration for the Vault client authentication.
type VaultConfig struct {
	Address     string
	AuthRole    string
	AuthMount   string
	AuthJWTPath string
}

// +kubebuilder:rbac:groups=security.platform.io,resources=vaultkubernetesrolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=security.platform.io,resources=vaultkubernetesrolebindings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=security.platform.io,resources=vaultkubernetesrolebindings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *VaultKubernetesRoleBindingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// 1. Fetch the VaultKubernetesRoleBinding instance
	var vkrb securityv1alpha1.VaultKubernetesRoleBinding
	if err := r.Get(ctx, req.NamespacedName, &vkrb); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Handle deletion
	if !vkrb.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, &vkrb)
	}

	// 3. Add finalizer if not present
	if !controllerutil.ContainsFinalizer(&vkrb, vaultKubernetesRoleBindingFinalizer) {
		controllerutil.AddFinalizer(&vkrb, vaultKubernetesRoleBindingFinalizer)
		if err := r.Update(ctx, &vkrb); err != nil {
			return ctrl.Result{}, err
		}
	}

	// 4. Authenticate to Vault
	if err := r.VaultClient.Authenticate(ctx, r.VaultConfig.AuthRole, r.VaultConfig.AuthMount, r.VaultConfig.AuthJWTPath); err != nil {
		log.Error(err, "Failed to authenticate to Vault")
		return r.updateStatus(ctx, &vkrb, "AuthenticationFailed", err.Error(), metav1.ConditionFalse)
	}

	// 5. Reconcile Vault Role
	if err := r.reconcileVaultRole(ctx, &vkrb); err != nil {
		log.Error(err, "Failed to reconcile Vault role")
		return r.updateStatus(ctx, &vkrb, "ReconciliationFailed", err.Error(), metav1.ConditionFalse)
	}

	// 6. Update Status to Ready
	return r.updateStatus(ctx, &vkrb, "Reconciled", "Vault role is synchronized", metav1.ConditionTrue)
}

func (r *VaultKubernetesRoleBindingReconciler) reconcileVaultRole(ctx context.Context, vkrb *securityv1alpha1.VaultKubernetesRoleBinding) error {
	log := logf.FromContext(ctx)

	desired := map[string]any{
		"bound_service_account_names":      vkrb.Spec.BoundServiceAccounts,
		"bound_service_account_namespaces": vkrb.Spec.BoundNamespaces,
		"token_policies":                   vkrb.Spec.TokenPolicies,
	}

	if vkrb.Spec.TokenTTL != "" {
		desired["token_ttl"] = vkrb.Spec.TokenTTL
	}
	if vkrb.Spec.Audience != "" {
		desired["audience"] = vkrb.Spec.Audience
	}

	actual, err := r.VaultClient.GetKubernetesRole(ctx, vkrb.Spec.AuthMount, vkrb.Spec.RoleName)
	if err != nil {
		return err
	}

	if actual != nil && r.isRoleUpToDate(desired, actual) {
		log.Info("Vault role is already up to date", "role", vkrb.Spec.RoleName)
		return nil
	}

	log.Info("Applying Vault role", "role", vkrb.Spec.RoleName)
	return r.VaultClient.ApplyKubernetesRole(ctx, vkrb.Spec.AuthMount, vkrb.Spec.RoleName, desired)
}

func (r *VaultKubernetesRoleBindingReconciler) isRoleUpToDate(desired, actual map[string]any) bool {
	for k, v := range desired {
		actualVal, exists := actual[k]
		if !exists {
			return false
		}

		// Vault might return slices in different formats or types
		// Simple comparison for now, might need more robust slice comparison
		if !reflect.DeepEqual(v, actualVal) {
			// Special handling for slices as Vault might return them as []interface{}
			if reflect.TypeOf(v).Kind() == reflect.Slice && reflect.TypeOf(actualVal).Kind() == reflect.Slice {
				vd := reflect.ValueOf(v)
				va := reflect.ValueOf(actualVal)
				if vd.Len() != va.Len() {
					return false
				}
				// Convert to string slices for easier comparison if possible
				vSlice := fmt.Sprintf("%v", v)
				aSlice := fmt.Sprintf("%v", actualVal)
				if vSlice != aSlice {
					return false
				}
				continue
			}
			return false
		}
	}
	return true
}

func (r *VaultKubernetesRoleBindingReconciler) reconcileDelete(ctx context.Context, vkrb *securityv1alpha1.VaultKubernetesRoleBinding) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	if controllerutil.ContainsFinalizer(vkrb, vaultKubernetesRoleBindingFinalizer) {
		// Authenticate to Vault
		if err := r.VaultClient.Authenticate(ctx, r.VaultConfig.AuthRole, r.VaultConfig.AuthMount, r.VaultConfig.AuthJWTPath); err != nil {
			log.Error(err, "Failed to authenticate to Vault for deletion")
			return ctrl.Result{}, err
		}

		log.Info("Deleting Vault role", "role", vkrb.Spec.RoleName)
		if err := r.VaultClient.DeleteKubernetesRole(ctx, vkrb.Spec.AuthMount, vkrb.Spec.RoleName); err != nil {
			log.Error(err, "Failed to delete Vault role")
			return ctrl.Result{}, err
		}

		controllerutil.RemoveFinalizer(vkrb, vaultKubernetesRoleBindingFinalizer)
		if err := r.Update(ctx, vkrb); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *VaultKubernetesRoleBindingReconciler) updateStatus(ctx context.Context, vkrb *securityv1alpha1.VaultKubernetesRoleBinding, reason, message string, status metav1.ConditionStatus) (ctrl.Result, error) {
	vkrb.Status.ObservedGeneration = vkrb.Generation

	condition := metav1.Condition{
		Type:               conditionTypeReady,
		Status:             status,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: vkrb.Generation,
	}

	meta.SetStatusCondition(&vkrb.Status.Conditions, condition)

	if err := r.Status().Update(ctx, vkrb); err != nil {
		return ctrl.Result{}, err
	}

	if status == metav1.ConditionFalse {
		return ctrl.Result{}, fmt.Errorf("%s: %s", reason, message)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VaultKubernetesRoleBindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1alpha1.VaultKubernetesRoleBinding{}).
		Named("security-vaultkubernetesrolebinding").
		Complete(r)
}
