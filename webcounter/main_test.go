// Contains the tests setup
package main

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/go-redis/redis"
	rediscounter "github.com/slavrd/go-redis-counter"
)

// commnad line flags
var clearRKey = flag.Bool("d", false, "clear redis key if already present")

func TestMain(m *testing.M) {

	// parse command line flags
	wcInit()

	// create a redis.Client to interract with redis
	c := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		DB:       *redisDB,
		Password: *redisPass,
	})

	// confirm the client is working
	_, err := c.Ping().Result()
	if err != nil {
		log.Fatalf("error creating redis.Client: %v", err)
	}

	// check if key is already present
	_, err = c.Get(*redisKey).Result()
	if err != redis.Nil {
		log.Printf("WARN: key %q already exist in redis db %v", *redisKey, *redisDB)
		if !*clearRKey {
			log.Fatal("aborting tests: can pass:\n  -d flag to allow modifications\n  -redis-db N to use a different database")
		}
	}

	// initialize global redis counter
	counter, err = rediscounter.NewCounter(redisAddr, *redisPass, *redisKey, *redisDB)
	if err != nil {
		log.Fatalf("error loading metrics html template: %v", err)
	}

	os.Exit(m.Run())
}
