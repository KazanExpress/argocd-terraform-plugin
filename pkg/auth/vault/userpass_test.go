package vault_test

import (
	"testing"

	"github.com/KazanExpress/argocd-terraform-plugin/pkg/auth/vault"
	"github.com/KazanExpress/argocd-terraform-plugin/pkg/helpers"
)

func TestUserPassLogin(t *testing.T) {
	cluster, username, password := helpers.CreateTestUserPassVault(t)
	defer cluster.Cleanup()

	userpass := vault.NewUserPassAuth(username, password, "")

	if err := userpass.Authenticate(cluster.Cores[0].Client); err != nil {
		t.Fatalf("expected no errors but got: %s", err)
	}
}
