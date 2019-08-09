require 'erb'

Vagrant.configure("2") do |config|

    redis_addr = "192.168.2.11"
    redis_pass = "myRedisPa$$w0rd"

    vault_addr = "192.168.2.12"
    vault_token="devV@ultRootT0ken"
    vault_policy_name="rediscounter" # name for vault policy which restrics access to redis password secret only 

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

    config.vm.define 'vault' do |vault|

        vault.vm.box = "slavrd/xenial64"
        vault.vm.network "private_network", ip: vault_addr
        vault.vm.network "forwarded_port", guest: 8200, host: 8200

        vault.vm.provision "shell", privileged: false, path: "ops/scripts/provision_vault.sh"
        vault.vm.provision "shell", privileged: false, path: "ops/scripts/vault_setup_basic.sh", args: [vault_token]
        vault.vm.provision "shell", privileged: false, path: "ops/scripts/vault_add_kv_secret.sh", args: "kv/redispassword pass \'#{redis_pass}\'"
        vault.vm.provision "shell", privileged: false, path: "ops/scripts/vault_setup_policy.sh", args: ["#{vault_policy_name}", "/vagrant/ops/config/vault-access-policy.hcl"]


    end

    config.vm.define 'client' do |c|

        c.vm.box = "slavrd/xenial64"
        c.vm.network "private_network", ip: "192.168.2.21"

        c.vm.provision "shell", privileged: false, path: "ops/scripts/provision_client.sh"

        # set up environment variables for convinience
        c.vm.provision "shell", inline: "echo export REDIS_ADDR='#{redis_addr}' | sudo tee -a /home/vagrant/.profile"
        c.vm.provision "shell", inline: "echo export REDIS_PASS=\\''#{redis_pass}'\\' | sudo tee -a /home/vagrant/.profile"
        c.vm.provision "shell", inline: "echo export VAULT_ADDR='http://#{vault_addr }:8200' | sudo tee -a /home/vagrant/.profile"
        c.vm.provision "shell", privileged: false, path: "ops/scripts/vault_create_token.sh", args: ["http://#{vault_addr }:8200", "#{vault_token}","#{vault_policy_name}"]

    end

end