package vault

import (
	"context"
	"fmt"
)

// ApplyKubernetesRole creates or updates a Kubernetes auth role in Vault.
func (v *vaultClient) ApplyKubernetesRole(ctx context.Context, mount string, roleName string, config map[string]any) error {
	path := fmt.Sprintf("auth/%s/role/%s", mount, roleName)
	_, err := v.client.Logical().WriteWithContext(ctx, path, config)
	if err != nil {
		return fmt.Errorf("failed to apply kubernetes role %s on mount %s: %w", roleName, mount, err)
	}
	return nil
}

// DeleteKubernetesRole deletes a Kubernetes auth role from Vault.
func (v *vaultClient) DeleteKubernetesRole(ctx context.Context, mount string, roleName string) error {
	path := fmt.Sprintf("auth/%s/role/%s", mount, roleName)
	_, err := v.client.Logical().DeleteWithContext(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete kubernetes role %s on mount %s: %w", roleName, mount, err)
	}
	return nil
}

// GetKubernetesRole retrieves a Kubernetes auth role from Vault.
func (v *vaultClient) GetKubernetesRole(ctx context.Context, mount string, roleName string) (map[string]any, error) {
	path := fmt.Sprintf("auth/%s/role/%s", mount, roleName)
	secret, err := v.client.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read kubernetes role %s on mount %s: %w", roleName, mount, err)
	}
	if secret == nil {
		return nil, nil
	}
	return secret.Data, nil
}
