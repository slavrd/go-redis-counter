/*
Tests require a functioning Vault instance.
The Vault address and token need to be in VAULT_ADDR and VAULT_TOKEN variables.
A KV secret must exists at path 'kv/redispassword' with key 'pass' and value "myRedisPa$$w0rd"
*/
package main

import "testing"

func TestGetRedisPass(t *testing.T) {

	vc, err := newVaultClient()
	if err != nil {
		t.Fatalf("error initializing Vault client: %v", err)
	}

	r, err := getRedisPass(vc, "kv/redispassword", "pass")
	if err != nil {
		t.Fatalf("error retrieving vault secret: %v", err)
	}

	wpass := "myRedisPa$$w0rd"
	if r != wpass {
		t.Errorf("error retrieving vault secret: got: %s want: %s", r, wpass)
	}

}
