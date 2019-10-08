require 'erb'

Vagrant.configure("2") do |config|

    redis_addr = "192.168.2.11"
    redis_pass = "myRedisPa$$w0rd"

    vault_addr = "192.168.2.12"
    vault_token= "devV@ultRootT0ken"
    vault_rp_path = "kv/redispassword"
    vault_rp_key = "pass"
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

        vault.vm.box = "slavrd/vault"
        vault.vm.network "private_network", ip: vault_addr
        vault.vm.network "forwarded_port", guest: 8200, host: 8200

        vault.vm.provision "shell", inline: "/etc/vault.d/scripts/vault_init.sh"
        vault.vm.provision "shell", inline: "/etc/vault.d/scripts/vault_unseal.sh", run: "always"
        vault.vm.provision "shell", privileged: false, path: "ops/scripts/vault_setup_basic.sh"
        vault.vm.provision "shell", privileged: false, path: "ops/scripts/vault_add_kv_secret.sh", args: "#{vault_rp_path} #{vault_rp_key} \'#{redis_pass}\'"
        vault.vm.provision "shell", privileged: false, path: "ops/scripts/vault_setup_policy.sh", args: ["#{vault_policy_name}", "/vagrant/ops/config/vault-access-policy.hcl"]
        vault.vm.provision "shell", privileged: false, path: "ops/scripts/vault_create_token_local.sh", args: "#{vault_policy_name}"

    end

    config.vm.define 'client' do |c|

        c.vm.box = "slavrd/xenial64-golang"
        c.vm.network "private_network", ip: "192.168.2.21"
        c.vm.synced_folder ".", "/home/vagrant/go/src/github.com/slavrd/go-redis-counter"

        c.vm.provision "shell", inline: "chown -R vagrant:vagrant /home/vagrant/go"
        # set up environment variables for convinience
        c.vm.provision "shell", privileged: false, path: "ops/scripts/provision_client_env.sh", args: "#{redis_addr} '#{redis_pass}' http://#{vault_addr}:8200"

        # TODO: Create a script to setup the Vault token
        c.vm.provision "shell", privileged: false, path: "ops/scripts/client_set_vault_token.sh"

    end

    config.vm.define 'webserver' do |w|

        w.vm.box = "slavrd/go-redis-counter"
        w.vm.network "private_network", ip: "192.168.2.31"
        w.vm.network "forwarded_port", guest: 8000, host: 8000

        ## generate webcounter service environment config from template
        File.write("ops/config/environment.conf", ERB.new(File.read("ops/config/environment.conf.erb")).result(binding))

        ## provision webserver VM, depends on the generated config
        w.vm.provision "file", source: "ops/config/environment.conf", destination: "/tmp/environment.conf"

        ## TODO: run script to set the vault token in /tmp/environment.conf
        w.vm.provision "shell", privileged: false, path: "ops/scripts/webserver_get_vault_token.sh", env: { "VAULT_IP_ADDR": "#{vault_addr}" }

        w.vm.provision "shell", path: "ops/scripts/provision_webserver.sh"
    
    end

end