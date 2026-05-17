package vault

import (
	"context"
	"fmt"
	"os"
)

// Authenticate authenticates the client using Kubernetes auth.
func (v *vaultClient) Authenticate(ctx context.Context, role string, mount string, jwtPath string) error {
	if jwtPath == "" {
		jwtPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	}

	jwt, err := os.ReadFile(jwtPath)
	if err != nil {
		return fmt.Errorf("failed to read service account jwt: %w", err)
	}

	params := map[string]any{
		"jwt":  string(jwt),
		"role": role,
	}

	loginPath := fmt.Sprintf("auth/%s/login", mount)
	secret, err := v.client.Logical().WriteWithContext(ctx, loginPath, params)
	if err != nil {
		return fmt.Errorf("failed to login to vault: %w", err)
	}

	if secret == nil || secret.Auth == nil {
		return fmt.Errorf("vault login returned no auth info")
	}

	v.client.SetToken(secret.Auth.ClientToken)
	return nil
}
