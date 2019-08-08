# Go redis counter package

[![Build Status](https://travis-ci.com/slavrd/go-redis-counter.svg?branch=master)](https://travis-ci.com/slavrd/go-redis-counter)

Package `rediscounter` provides a simple counter that stores its value in a redis server.

The package documentation is available as [godoc](https://godoc.org/github.com/slavrd/go-redis-counter)

## Console counter

Command rediscounter/console implements the rediscounter package. When run it will connect to the redis server, increase the value of the counter by one and display the resulting value.

There are two ways to specify the password for the redis server (if needed):

1. Pass it directly using the `-pass` command line flag.
2. Retrieve it from a Vault server KV secrets engine. To use this option use `-v-use` flag.
    * Vault address must be set in `VAULT_ADDR` environment variable.
    * Vault access token must be set in `VAULT_TOKEN` environment variable.
    * Vault secret's path and key can be specified using `-v-secret-path` and `-v-secret-key` flags.

### Building the command

* install [Golang](https://golang.org/dl/) or use the provided Vagrant [project](#vagrant-environment) which includes also a VM with redis.
* download the package - `go get github.com/slavrd/go-redis-counter`
* build the command - `go build $(go env GOPATH)/src/github.com/slavrd/go-redis-counter/console`

### Running the command

The binary accepts the following options. All of them have defaults set so it can potentially be run without setting any of them.

* `-a` - the redis server address. Default is `localhost`
* `-p` - the redis server port. Default is `6379`
* `-d` - the redis server db id. Default is `0`
* `-k` - the redis server key in which to store the counter's value. Default is `count`
* `-pass` - the redis server password. Default is empty string, meaning that the server does not require authentication
* `-v-use` - if redis server password should be retrieved from Vault. Default is `false`
* `-v-secret-path` - the Vault secret path. Default is `kv/redispassword`
* `-v-secret-key` - - the Vault secret path. Default is `pass`

For example:

Running `./console -a 192.168.0.1 -pass 'mypassword'` will make the app connect to redis @ 192.168.0.1:6379, db 0, key 'count' using password 'mypassword.'

## Vagrant environment

A Vagrant project that builds a VM with a redis server and a VM with Golang.

In case you are not familiar with Vagrant, a getting started guide can be found [here](https://www.vagrantup.com/intro/index.html).

### Using the project

The `redis` VM will have redis server running on port `6379` which will be mapped to the host as well. The redis server will have password authentication configured as well.

The `client` VM will have Golang installed so the application can be built and run. The redis server IP and password will be set in the environment variables `$REDIS_ADDR` and `$REDIS_PASS`. The Vault server address and token will be set as well in `$VAULT_ADDR` and `$VAULT_TOKEN`.

Example:

```bash
vagrant up # build the VMs
vagrant ssh client # login to the client VM

# commands below are executed on the client VM
go get github.com/slavrd/go-redis-counter # download the counter app source code
go build $(go env GOPATH)/src/github.com/slavrd/go-redis-counter/console # build the console counter app
./console -a $REDIS_ADDR -pass $REDIS_PASS # run the app passing the redis password directly
./console -a $REDIS_ADDR -v-use # run the app retrieving the redis password from Vault
exit # exit from the client VM to the host

vagrant destroy # destroy the vagrant VMs
```
