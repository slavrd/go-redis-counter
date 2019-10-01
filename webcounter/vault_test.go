/*
Tests require a functioning Vault instance.
The Vault address and token need to be in VAULT_ADDR and VAULT_TOKEN variables.
The Vault secret path and key must be set in VAULT_RP_PATH and VAULT_RP_KEY variables.
*/
package main

import (
	"os"
	"testing"
)

func TestGetRedisPass(t *testing.T) {

	vc, err := newVaultClient()
	if err != nil {
		t.Fatalf("error initializing Vault client: %v", err)
	}

	r, err := getRedisPass(vc, os.Getenv("VAULT_RP_PATH"), os.Getenv("VAULT_RP_KEY"))
	if err != nil {
		t.Fatalf("error retrieving vault secret: %v", err)
	}

	if r == "" {
		t.Errorf("error retrieving vault secret: got: %q", r)
	}

}
