// Contains the test setup
package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

// redis addres to use for the tests
var raddr string

// redis password to use for tests
var rpass string

// redis key to use for the tests
var rkey string

// redis db to use for the tests. This cannot be overridden by setting the CLI flag.
var rdb = 10

// a redis.Client to use for setting up redis test values
var c *redis.Client

// commnad line flags
var clearRKey = flag.Bool("d", false, "clear redis key if already present")

func init() {

	flag.Parse()

	// normalize redis address
	envRAddr := os.Getenv("REDIS_ADDR")
	if envRAddr != "" {
		raddr = envRAddr
	} else {
		raddr = *redisHost
	}
	if !strings.ContainsRune(raddr, ':') {
		raddr = strings.Join([]string{raddr, strconv.Itoa(*redisPort)}, ":")
	}

	// set redis password to environment variable if present
	envRPass := os.Getenv("REDIS_PASS")
	if envRPass != "" {
		rpass = envRPass
	} else {
		rpass = *redisPass
	}

	rkey = *redisKey

	// create a redis.Client to interract with redis
	c = redis.NewClient(&redis.Options{
		Addr:     raddr,
		DB:       rdb,
		Password: rpass,
	})

	// confirm the client is working
	_, err := c.Ping().Result()
	if err != nil {
		log.Fatalf("error creating redis.Client: %v", err)
	}

	// check if key is already present
	_, err = c.Get(rkey).Result()
	if err != redis.Nil {
		log.Printf("WARN: key %q already exist in redis db %v", rkey, rdb)
		if !*clearRKey {
			log.Fatal("aborting tests: can pass -d flag to allow existing key modification")
		}
	}

	// set up tmlCounterCtx
	htmlCounterCtx, err = newHTMLCounter(raddr, rpass, rkey, rdb)
	if err != nil {
		log.Fatalf("error createting htmlCounter instance for htmlCounterCtx: %v", err)
	}

	// set up htmlCounterTpl
	htmlCounterTpl, err = loadTemplate(*tplPath)
	if err != nil {
		log.Fatalf("error loading html template: %v", err)
	}

}
