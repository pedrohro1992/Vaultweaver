# Vaultweaver

Vaultweaver is a Kubebuilder-based Kubernetes operator designed to manage the integration between Kubernetes ServiceAccounts and HashiCorp Vault dynamically using Custom Resource Definitions (CRDs).

## Features

- **Declarative Vault Auth Roles**: Manage `auth/kubernetes/role` entries directly from Kubernetes.
- **Declarative Vault Policies**: Manage Vault policies as Kubernetes resources.
- **Dynamic Authentication**: Authenticates to Vault using its own ServiceAccount JWT (Projected Volumes). No static tokens or AppRole secrets required.
- **Idempotent Reconciliation**: Detects drift between Kubernetes and Vault and synchronizes state automatically.
- **Native Kubernetes Experience**: Uses `metav1.Condition` for status tracking and standard controller-runtime idioms.

## Custom Resources

### VaultKubernetesRoleBinding
Manages a Vault Kubernetes Auth role.

```yaml
apiVersion: security.platform.io/v1alpha1
kind: VaultKubernetesRoleBinding
metadata:
  name: cert-manager
spec:
  authMount: kubernetes
  roleName: cert-manager
  boundServiceAccounts:
    - cert-manager
  boundNamespaces:
    - cert-manager
  tokenPolicies:
    - pki-cert-manager
  tokenTTL: 1h
  audience: vault
```

### VaultPolicy
Manages a Vault policy.

```yaml
apiVersion: security.platform.io/v1alpha1
kind: VaultPolicy
metadata:
  name: pki-cert-manager
spec:
  vaultPolicyName: pki-cert-manager
  policy: |
    path "pki/sign/example-dot-com" {
      capabilities = ["update"]
    }
```

## Getting Started

### Vault Setup
The operator requires a Vault instance with the Kubernetes auth method enabled. The operator itself must have a role in Vault that allows it to manage other roles and policies.

**Operator Vault Policy Example:**
```hcl
# Allow managing kubernetes auth roles
path "auth/kubernetes/role/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

# Allow managing policies
path "sys/policy/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
```

### Operator Configuration
The operator is configured via environment variables:

- `VAULT_ADDR`: The address of the Vault server.
- `VAULT_AUTH_KUBERNETES_ROLE`: The Vault role the operator uses to authenticate.
- `VAULT_AUTH_KUBERNETES_MOUNT`: (Optional) The mount path for Kubernetes auth (default: `kubernetes`).
- `VAULT_AUTH_KUBERNETES_JWT_PATH`: (Optional) Path to the ServiceAccount JWT (default: `/var/run/secrets/kubernetes.io/serviceaccount/token`).

### Local Development
To run the operator locally, generate a token for its ServiceAccount and point the operator to it:

1. Generate a token:
   ```bash
   kubectl create token vaultweaver-operator -n vaultweaver-system --audience vault > operator.jwt
   ```
2. Export environment variables:
   ```bash
   export VAULT_ADDR="http://127.0.0.1:8200"
   export VAULT_AUTH_KUBERNETES_ROLE="vaultweaver-operator"
   export VAULT_AUTH_KUBERNETES_JWT_PATH="$(pwd)/operator.jwt"
   ```
3. Run the operator:
   ```bash
   make run
   ```

## Installation
To install the CRDs into your cluster:
```bash
make install
```

To deploy the operator to the cluster:
```bash
make deploy IMG=<some-registry>/vaultweaver:tag
```

## License
Copyright 2026.
Licensed under the Apache License, Version 2.0.
