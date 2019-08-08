// Command console starts a console client that implements the rediscounter package.
// It will increment the counter value and display the result
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	rediscounter "github.com/slavrd/go-redis-counter"
)

func main() {

	// Command line options
	raddr := flag.String("a", "localhost", "redis server address")
	rport := flag.String("p", "6379", "redis server port")
	rdb := flag.Int("d", 0, "redis server database id")
	rkey := flag.String("k", "count", "redis key to store counter value")
	rpass := flag.String("pass", "", "redis server password")
	useVault := flag.Bool("v-use", false, "retrieve redis password form Vault server")
	vSecretPath := flag.String("v-secret-path", "kv/redispassword", "vault path for redis password secret")
	vSecretKey := flag.String("v-secret-key", "pass", "vault secret key for redis password")
	flag.Parse()

	// build the redis server address
	rFullAddr := strings.Join([]string{*raddr, *rport}, ":")

	// get the redis server password, depending on provided options
	var rPass string
	if *useVault {

		vc, err := newVaultClient()
		if err != nil {
			log.Fatalf("error retrieving redis server password from Vault: %v", err)
		}

		rPass, err = getRedisPass(vc, *vSecretPath, *vSecretKey)
		if err != nil {
			log.Fatalf("error retrieving redis server password from Vault: %v", err)
		}

	} else {
		rPass = *rpass
	}

	// initialize counter
	c, err := rediscounter.NewCounter(rFullAddr, rPass, *rkey, *rdb)
	if err != nil {
		log.Fatalf("error initializing counter: %v", err)
	}

	// increment counter
	r, err := c.IncrBy(1)
	if err != nil {
		log.Fatalf("error incrementing counter: %v", err)
	}

	fmt.Println(r)
}

// getRedisPass gets the redis password from a Vault KV strore
func getRedisPass(vc *api.Client, path string, key string) (string, error) {

	s, err := vc.Logical().Read(path)
	if err != nil {
		return "", err
	}
	pass := s.Data[key].(string) // ignoring type assertion error
	return pass, nil
}

// creates Vault client according to environment variables settings
// VAULT_ADDR is used to set the Vault server address. If not set default is 'https://localhost:8200'
// VAULT_TOKEN is used to set the Vault access token
func newVaultClient() (*api.Client, error) {

	vconf := api.DefaultConfig()

	va := os.Getenv("VAULT_ADDR")
	if va == "" {
		return nil, fmt.Errorf("error initializing vault api client: environment variable 'VAULT_ADDR' is not set")
	}
	vconf.Address = va

	vaultClient, err := api.NewClient(vconf)
	if err != nil {
		return nil, fmt.Errorf("error initializing vault api client: %v", err)
	}

	vt := os.Getenv("VAULT_TOKEN")
	if vt == "" {
		return nil, fmt.Errorf("error initializing vault api client: environment variable 'VAULT_TOKEN' is not set")
	}
	vaultClient.SetToken(vt)

	return vaultClient, nil
}
