package vault

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

// Client defines the interface for interacting with Vault.
type Client interface {
	// Authenticate authenticates the client using Kubernetes auth.
	Authenticate(ctx context.Context, role string, mount string, jwtPath string) error
	// ApplyKubernetesRole creates or updates a Kubernetes auth role in Vault.
	ApplyKubernetesRole(ctx context.Context, mount string, roleName string, config map[string]any) error
	// DeleteKubernetesRole deletes a Kubernetes auth role from Vault.
	DeleteKubernetesRole(ctx context.Context, mount string, roleName string) error
	// GetKubernetesRole retrieves a Kubernetes auth role from Vault.
	GetKubernetesRole(ctx context.Context, mount string, roleName string) (map[string]any, error)

	// ApplyPolicy creates or updates a policy in Vault.
	ApplyPolicy(ctx context.Context, name string, hcl string) error
	// DeletePolicy deletes a policy from Vault.
	DeletePolicy(ctx context.Context, name string) error
	// GetPolicy retrieves a policy from Vault.
	GetPolicy(ctx context.Context, name string) (string, error)
}

type vaultClient struct {
	client *vault.Client
}

// NewClient creates a new Vault client.
func NewClient(address string) (Client, error) {
	config := vault.DefaultConfig()
	config.Address = address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	return &vaultClient{
		client: client,
	}, nil
}
