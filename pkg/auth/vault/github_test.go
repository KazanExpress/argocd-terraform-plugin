package vault_test

import (
	"testing"

	"github.com/KazanExpress/argocd-terraform-plugin/pkg/auth/vault"
	"github.com/KazanExpress/argocd-terraform-plugin/pkg/helpers"
)

// Need to find a way to mock GitHub Auth within Vault
func TestGithubLogin(t *testing.T) {
	cluster := helpers.CreateTestAuthVault(t)
	defer cluster.Cleanup()

	github := vault.NewGithubAuth("123", "")

	err := github.Authenticate(cluster.Cores[0].Client)
	if err != nil {
		t.Fatalf("expected no errors but got: %s", err)
	}
}
