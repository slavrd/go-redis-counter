package main

import (
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// getRedisPass gets the redis password from a Vault KV strore
func getRedisPass(vc *vaultapi.Client, path string, key string) (string, error) {

	s, err := vc.Logical().Read(path)
	if err != nil {
		return "", err
	}
	if s == nil {
		return "", fmt.Errorf("vault retrutned <nil> secret")
	}
	pass := s.Data[key].(string) // ignoring type assertion error
	return pass, nil

}

// newVaultClient creates Vault client according to environment variables settings.
// VAULT_ADDR is used to set the Vault server address. If not set default is 'https://localhost:8200'
// VAULT_TOKEN is used to set the Vault access token.
// It is currently just a wrapper for the NewClient() so just a placeholder for more detailed init when needed.
func newVaultClient() (*vaultapi.Client, error) {

	vconf := vaultapi.DefaultConfig()

	vaultClient, err := vaultapi.NewClient(vconf)
	if err != nil {
		return nil, fmt.Errorf("error initializing vault api client: %v", err)
	}

	return vaultClient, nil
}
