package vault

import (
	"context"
	"fmt"
)

// ApplyPolicy creates or updates a policy in Vault.
func (v *vaultClient) ApplyPolicy(ctx context.Context, name string, hcl string) error {
	err := v.client.Sys().PutPolicy(name, hcl)
	if err != nil {
		return fmt.Errorf("failed to apply vault policy %s: %w", name, err)
	}
	return nil
}

// DeletePolicy deletes a policy from Vault.
func (v *vaultClient) DeletePolicy(ctx context.Context, name string) error {
	err := v.client.Sys().DeletePolicy(name)
	if err != nil {
		return fmt.Errorf("failed to delete vault policy %s: %w", name, err)
	}
	return nil
}

// GetPolicy retrieves a policy from Vault.
func (v *vaultClient) GetPolicy(ctx context.Context, name string) (string, error) {
	policy, err := v.client.Sys().GetPolicy(name)
	if err != nil {
		return "", fmt.Errorf("failed to read vault policy %s: %w", name, err)
	}
	return policy, nil
}
