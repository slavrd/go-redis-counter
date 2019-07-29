# Go redis counter package

[![Build Status](https://travis-ci.com/slavrd/go-redis-counter.svg?branch=master)](https://travis-ci.com/slavrd/go-redis-counter)

Package `rediscounter` provides a simple counter that stores its value in a redis server.

The package documentation is available as [godoc](https://godoc.org/github.com/slavrd/go-redis-counter)

## Console counter

Command rediscounter/console implements the rediscounter package. When run it will connect to the redis server, increase the value of the counter by one and display the resulting value.

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

For example:

Running `./console -a 192.168.0.1 -pass 'mypassword'` will make the app connect to redis @ 192.168.0.1:6379, db 0, key 'count' using password 'mypassword.'

## Vagrant environment

A Vagrant project that builds a VM with a redis server and a VM with Golang.

In case you are not familiar with Vagrant, a getting started guide can be found [here](https://www.vagrantup.com/intro/index.html).

### Using the project

The `redis` VM will have redis server running on port `6379` which will be mapped to the host as well. The redis server will have password authentication configured as well.

The `client` VM will have Golang installed so the application can be built and run. The redis server IP and password will be set in the environment variables `$REDIS_ADDR` and `$REDIS_PASS`

Example:

```bash
vagrant up # build the VMs
vagrant ssh client # login to the client VM

# commands below are executed on the client VM
go get github.com/slavrd/go-redis-counter # download the counter app source code
go build $(go env GOPATH)/src/github.com/slavrd/go-redis-counter/console # build the console counter app
./console -a $REDIS_ADDR -pass $REDIS_PASS # run the app
exit # exit from the client VM to the host

vagrant destroy # destroy the vagrant VMs
```
