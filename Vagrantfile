
Vagrant.configure("2") do |config| 
    
    
    # ---------------------------------------------- 
    # Set up dev environment locally
    # ---------------------------------------------- 
    config.vm.define "local" do |local_config|
        local_config.vm.box = "phusion/ubuntu-14.04-amd64"
        
        local_config.vm.provision "shell", path: "scripts/provision.sh"
            
        local_config.vm.network :private_network, ip: "10.3.3.3"

        local_config.vm.network :forwarded_port, guest:80,   host:80,   auto_correct: true   
        
        local_config.vm.synced_folder "./", "/opt/gopath/src/github.com/hoodiez/kapal/"

        local_config.vm.provider :virtualbox do |vb|            
             vb.customize ["modifyvm", :id, "--memory", "1024"]
        end
    end

end