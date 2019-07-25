require 'erb'

Vagrant.configure("2") do |config|

    redis_addr = "192.168.2.11"
    redis_pass = "myRedisPa$$w0rd"

    config.vm.define 'redis' do |r|

        r.vm.box = "slavrd/redis64"
        r.vm.network "private_network", ip: redis_addr
        r.vm.network "forwarded_port", guest: 6379, host: 6379

        ## generate redis config from template
        File.write("ops/config/redis.conf", ERB.new(File.read("ops/config/redis.conf.erb")).result(binding))

        ## provision redis VM, depends on the generated config
        r.vm.provision "file", source: "ops/config/redis.conf", destination: "/tmp/redis.conf"
        r.vm.provision "shell", inline: "sudo cp /tmp/redis.conf /etc/redis/redis.conf && sudo systemctl restart redis-server.service"

    end

    config.vm.define 'client' do |c|

        c.vm.box = "slavrd/xenial64"
        c.vm.network "private_network", ip: "192.168.2.21"

        c.vm.provision "shell", privileged: false, path: "ops/scripts/provision.sh"

        c.vm.provision "shell", inline: "echo export REDIS_ADDR='#{redis_addr}' | sudo tee -a /home/vagrant/.profile"
        c.vm.provision "shell", inline: "echo export REDIS_PASS=\\''#{redis_pass}'\\' | sudo tee -a /home/vagrant/.profile"

    end

end