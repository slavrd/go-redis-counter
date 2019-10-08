# Go redis counter package

[![Build Status](https://travis-ci.com/slavrd/go-redis-counter.svg?branch=master)](https://travis-ci.com/slavrd/go-redis-counter)
[![Godoc Ref](https://godoc.org/github.com/slavrd/go-redis-counter?status.svg)](https://https://godoc.org/github.com/slavrd/go-redis-counter)

Package `rediscounter` provides a simple counter that stores its value in a redis server.

The package documentation is available as [godoc](https://godoc.org/github.com/slavrd/go-redis-counter)

The project includes basic implementations of the package:

* a console application in folder `console/` - [readme](console/README.md).
* a web application in folder `webcounter/`- [readme](webcounter/README.md).

## Vagrant environment

The repository includes a Vagrant project that builds 3 VMs with redis server, vault server and Golang installed.

In case you are not familiar with Vagrant, a getting started guide can be found [here](https://www.vagrantup.com/intro/index.html).

### Using the project

The `redis` VM will have redis server running on port `6379` which will be mapped to the host as well. The redis server will have password authentication configured as well.

The `vault` VM will have Vault installed as a service. There is a script set to unseal vault each time the VM boots. The vault root token is set-up as environment variable for the `vagrant` user.

The `client` VM will have Golang installed so the application can be built and run. The redis server IP and password will be set in the environment variables `$REDIS_ADDR` and `$REDIS_PASS`. The Vault server address and access token will be set as well in `$VAULT_ADDR` and `$VAULT_TOKEN` respectively.

The `webserver` VM has the latest version of the webcounter app installed and running as a service.

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
## TODO

- [x] Console implementation
- [x] `console`: add integration with Vault KV secrets engine. Application should be able to retrieve redis password from Vault.
- [x] Web implementation
- [x] Make Vagrant mount the project in the `$GOPATH` of the client VM.
- [x] Make Vagrant use a "golang" box instead of provisioning the client VM each time.
- [x] Make go tests warn and stop if redis key is already set. Make it possible to force running the tests with a command flag.
- [x] `webcounter`: add a `/health` check method that tests the redis server connection.
- [x] `webcounter`: add a `/metrics` method which will report how many times each path has been requested.
- [x] `webcounter`: add a `/crash` method which will stop the webserver
- [x] `webcounter`: add a `/reset` method which will reset the counter
- [x] `webcounter`: add a `/decr` method which will decrease the counter by `1`. Should be guarded from going below `0`
- [x] `webcounter`: start even if redis connection is unavailable.
- [ ] `webcounter`: implement mutual exclusion lock on redis level, so multiple app instances can be used.
- [X] Box with webcounter running as a service
- [ ] Update box with webcounter service running as a non privileged user
- [ ] `webcounter`: add integration with Vault KV secrets engine.
    - [x] add the integration to the app code
    - [x] use a vagrant box with vault
    - [ ] use an AWS AMI with vault in terraform config.
- [x] `webcounter`: update UI to call the methods with buttons.
- [x] packer project that creates AWS AMI with `webcounter` app installed as a service
- [x] terraform project that deploys `webcounter` app and its redis server in AWS
- [x] add network module to the terraform project
- [ ] terraform create AWS subnets in different AZs
- [ ] replace publicly accessible webcounter EC2 instances with load balancer
- [ ] terraform remove redis password from EC2 instances user data

## Related projects

* Vagrant [counter box](https://github.com/slavrd/packer-go-redis-counter-vagrant) - a Packer project that builds a Vagrant Virtualbox box with webcounter installed as a service.
* AWS [counter AMI](https://github.com/slavrd/packer-go-redis-counter-aws) - a Packer project that builds an AWS AMI with webcounter installed as a service.
* AWS [redis AMI](https://github.com/slavrd/packer-aws-redis64) - a Packer project that build an AWS AMI with Redis server installed.
* AWS [Terraform deployment](https://github.com/slavrd/terraform-go-redis-counter) - a Terraform configuration to deploy the webcounter and redis AMIs in AWS. 
