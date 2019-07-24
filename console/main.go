// Command console starts a console client that implements the rediscounter package.
// It will increment the counter value and display the result.package main
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	rediscounter "github.com/slavrd/go-redis-counter"
)

func main() {

	// Command line options
	raddr := flag.String("a", "localhost", "redis server address")
	rport := flag.String("p", "6379", "redis server port")
	rdb := flag.Int("d", 0, "redis server database id")
	rkey := flag.String("k", "count", "redis key to store counter value")
	rpass := flag.String("pass", "", "redis server password")
	flag.Parse()

	rFullAddr := strings.Join([]string{*raddr, *rport}, ":")

	// initialize counter
	c, err := rediscounter.NewCounter(rFullAddr, *rpass, *rkey, *rdb)
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
