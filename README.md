# Go redis counter package

[![Build Status](https://travis-ci.com/slavrd/go-redis-counter.svg?branch=master)](https://travis-ci.com/slavrd/go-redis-counter)

Package `rediscounter` provides a simple counter that stores its value in a redis server.

The package documentation is available as [godoc](https://godoc.org/github.com/slavrd/go-redis-counter)

The project includes a console application in folder `console/` which implements the `rediscounter` package. Detailed description is in its [readme](console/README.md)

## Vagrant environment

The repository includes a Vagrant project that builds 3 VMs with redis server, vault server and Golang installed.

In case you are not familiar with Vagrant, a getting started guide can be found [here](https://www.vagrantup.com/intro/index.html).

### Using the project

The `redis` VM will have redis server running on port `6379` which will be mapped to the host as well. The redis server will have password authentication configured as well.

The `vault` VM will have Vault installed and running in `dev` mode on port `8200` which will be mapped to the host as well.

The `client` VM will have Golang installed so the application can be built and run. The redis server IP and password will be set in the environment variables `$REDIS_ADDR` and `$REDIS_PASS`. The Vault server address and root token will be set as well in `$VAULT_ADDR` and `$VAULT_TOKEN` respectively.

Example:

```bash
vagrant up # build the VMs
vagrant ssh client # login to the client VM

# commands below are executed on the client VM

# download the counter app source code and all dependencies
go get github.com/slavrd/go-redis-counter/...

# build the console counter app
go build $(go env GOPATH)/src/github.com/slavrd/go-redis-counter/console

# run the app passing the redis password directly
./console -a $REDIS_ADDR -pass $REDIS_PASS

# run the app retrieving the redis password from Vault
./console -a $REDIS_ADDR -v-use

exit # exit from the client VM to the host

vagrant destroy # destroy the vagrant VMs
```
