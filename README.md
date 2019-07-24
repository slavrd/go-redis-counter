# Go redis counter package

Package `rediscounter` provides a simple counter that stores its value in a redis server.

The package documentation is available as [godoc](https://godoc.org/github.com/slavrd/go-redis-counter)

## Console counter

Command rediscounter/console implements the rediscounter package. When run it will connect to the redis server, increase the value of the counter by one and display the resulting value.

### Building the command

* install [Golang](https://golang.org/dl/) or use the provided Vagrant project which includes also a VM with redis.
* download the package - `go get github.com/slavrd/go-redis-counter`
* build the command - `go build $GOHOME/github.com/slavrd/go-redis-counter/console`

### Running the command

The binnary accepts the following options. All of them have defaults set so it can potentially be run without setting any of them.

* `-a` - the redis server address. Default is `localhost`
* `-p` - the redis server port. Default is `6379`
* `-d` - the redis server db id. Default is `0`
* `-k` - the redis server key in which to store the counter's value. Default is `count`
* `-pass` - the redis server password. Default is empty string, meaning that the server does not require authentication

For example:

Running `./console -a 192.168.0.1 -pass 'mypassword'` will make the app connecto to redis @ 192.168.0.1:6379, db 0, key 'count' using password 'mypassword.'
