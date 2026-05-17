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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	securityv1alpha1 "github.com/pedrohro1992/Vaultweaver/api/security/v1alpha1"
	"github.com/pedrohro1992/Vaultweaver/internal/vault"
)

var _ = Describe("VaultKubernetesRoleBinding Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}
		vaultkubernetesrolebinding := &securityv1alpha1.VaultKubernetesRoleBinding{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind VaultKubernetesRoleBinding")
			err := k8sClient.Get(ctx, typeNamespacedName, vaultkubernetesrolebinding)
			if err != nil && errors.IsNotFound(err) {
				resource := &securityv1alpha1.VaultKubernetesRoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: securityv1alpha1.VaultKubernetesRoleBindingSpec{
						AuthMount:            "kubernetes",
						RoleName:             "test-role",
						BoundServiceAccounts: []string{"test-sa"},
						BoundNamespaces:      []string{"default"},
						TokenPolicies:        []string{"default"},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &securityv1alpha1.VaultKubernetesRoleBinding{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance VaultKubernetesRoleBinding")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &VaultKubernetesRoleBindingReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
				VaultClient: &vault.MockClient{
					GetKubernetesRoleFunc: func(ctx context.Context, mount, roleName string) (map[string]any, error) {
						return nil, nil // Role does not exist yet
					},
					ApplyKubernetesRoleFunc: func(ctx context.Context, mount, roleName string, config map[string]any) error {
						return nil
					},
					AuthenticateFunc: func(ctx context.Context, role, mount, jwtPath string) error {
						return nil
					},
				},
				VaultConfig: VaultConfig{
					AuthRole: "operator-role",
				},
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			resource := &securityv1alpha1.VaultKubernetesRoleBinding{}
			err = k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())
			Expect(resource.Status.Conditions).NotTo(BeEmpty())
			Expect(resource.Status.Conditions[0].Status).To(Equal(metav1.ConditionTrue))
		})
	})
})
