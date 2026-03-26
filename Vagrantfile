Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/jammy64"

  config.vm.define "app" do |app|
    app.vm.hostname = "vm-app"
    app.vm.network "private_network", ip: "192.168.56.10"
    app.vm.network "forwarded_port", guest: 8080, host: 8080

    app.vm.provider "virtualbox" do |vb|
      vb.memory = "2048"
      vb.cpus = 2
    end

    app.vm.provision "shell", inline: <<-SHELL
    apt-get update -q
    echo "vm-app pronta!"
    SHELL
  end

  config.vm.define "observ" do |observ|
    observ.vm.hostname = "vm-observ"
    observ.vm.network "private_network", ip: "192.168.56.20"

    observ.vm.network "forwarded_port", guest: 3000, host: 3000
    observ.vm.network "forwarded_port", guest: 9090, host: 9090
    observ.vm.network "forwarded_port", guest: 3100, host: 3100

    observ.vm.provider "virtualbox" do |vb|
      vb.name   = "vm-observ"
      vb.memory = 3072
      vb.cpus   = 2
    end

    observ.vm.provision "shell", inline: <<-SHELL
      apt-get update -q
      echo "vm-observ pronta!"
    SHELL
  end
end
