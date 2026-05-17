package vault

import (
	"context"
)

// MockClient is a mock implementation of the Vault Client interface.
type MockClient struct {
	AuthenticateFunc         func(ctx context.Context, role string, mount string, jwtPath string) error
	ApplyKubernetesRoleFunc  func(ctx context.Context, mount string, roleName string, config map[string]any) error
	DeleteKubernetesRoleFunc func(ctx context.Context, mount string, roleName string) error
	GetKubernetesRoleFunc    func(ctx context.Context, mount string, roleName string) (map[string]any, error)
	ApplyPolicyFunc          func(ctx context.Context, name string, hcl string) error
	DeletePolicyFunc         func(ctx context.Context, name string) error
	GetPolicyFunc            func(ctx context.Context, name string) (string, error)
}

func (m *MockClient) Authenticate(ctx context.Context, role string, mount string, jwtPath string) error {
	if m.AuthenticateFunc != nil {
		return m.AuthenticateFunc(ctx, role, mount, jwtPath)
	}
	return nil
}

func (m *MockClient) ApplyKubernetesRole(ctx context.Context, mount string, roleName string, config map[string]any) error {
	if m.ApplyKubernetesRoleFunc != nil {
		return m.ApplyKubernetesRoleFunc(ctx, mount, roleName, config)
	}
	return nil
}

func (m *MockClient) DeleteKubernetesRole(ctx context.Context, mount string, roleName string) error {
	if m.DeleteKubernetesRoleFunc != nil {
		return m.DeleteKubernetesRoleFunc(ctx, mount, roleName)
	}
	return nil
}

func (m *MockClient) GetKubernetesRole(ctx context.Context, mount string, roleName string) (map[string]any, error) {
	if m.GetKubernetesRoleFunc != nil {
		return m.GetKubernetesRoleFunc(ctx, mount, roleName)
	}
	return nil, nil
}

func (m *MockClient) ApplyPolicy(ctx context.Context, name string, hcl string) error {
	if m.ApplyPolicyFunc != nil {
		return m.ApplyPolicyFunc(ctx, name, hcl)
	}
	return nil
}

func (m *MockClient) DeletePolicy(ctx context.Context, name string) error {
	if m.DeletePolicyFunc != nil {
		return m.DeletePolicyFunc(ctx, name)
	}
	return nil
}

func (m *MockClient) GetPolicy(ctx context.Context, name string) (string, error) {
	if m.GetPolicyFunc != nil {
		return m.GetPolicyFunc(ctx, name)
	}
	return "", nil
}
