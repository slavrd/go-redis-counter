# Console redis counter

[![Build Status](https://travis-ci.com/slavrd/go-redis-counter.svg?branch=master)](https://travis-ci.com/slavrd/go-redis-counter)

A console application implementing the rediscounter package. When it is run it will connect to redis, increase the counter's value and display the result.

Redis connection options can be specified using different flags.

## Building the command

* install [Golang](https://golang.org/dl/) or use the provided Vagrant [project](../README.md#vagrant-environment) which includes also a VM with redis.
* download the package - `go get github.com/slavrd/go-redis-counter`
* build the command - `go build $(go env GOPATH)/src/github.com/slavrd/go-redis-counter/console`

## Running the command

The binary accepts the following options. All of them have defaults set so it can potentially be run without using any of them.

* `-a` - the redis server address. Default is `localhost`.
* `-p` - the redis server port. Default is `6379`.
* `-d` - the redis server db id. Default is `0`.
* `-k` - the redis server key in which to store the counter's value. Default is `count`.

There are two ways to specify the password for the redis server (if needed):

1. Pass it directly using the `-pass` command line flag. Default is `""` which will result in connection without authentication.
2. Retrieve it from a Vault server KV secrets engine. To use this option use `-v-use` flag.
    * Vault address must be set in `VAULT_ADDR` environment variable.
    * Vault access token must be set in `VAULT_TOKEN` environment variable.
    * `-v-secret-path` - sets the vault KV secret path. Default is `kv/redispassword`.
    * `-v-secret-key` - sets tha vault KV secret key. Default is `pass`.

Examples:

`./console -a 192.168.0.1 -pass 'mypassword'` will make the app connect to redis @ 192.168.0.1:6379, db 0, key 'count' using password 'mypassword'.

`./console -a 192.168.0.1 -v-use` will make the app connect to redis @ 192.168.0.1:6379, db 0, key 'count' using a password retrieved from Vault.
