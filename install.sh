#! /bin/bash

function set_iptables {
    # Remove current iptables rules
    iptables -F
    # Default policies
    iptables -P INPUT DROP
    iptables -P OUTPUT ACCEPT
    iptables -P FORWARD DROP
    # Allow active connections
    iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
    # Allow ssh  in our network
    iptables -A INPUT -p tcp --dport 22 -s 163.117.170.0/24 -j ACCEPT
    iptables -A INPUT -p tcp --dport 22 -s 163.117.142.0/24 -j ACCEPT
    # Allow loopback
    iptables -I INPUT 1 -i lo -j ACCEPT
    # Allow ping
    iptables -A INPUT -p icmp --icmp-type 8 -j ACCEPT
}

function install_rpi1 {
    echo "##### Updating hostname, localtime and bashrc..."
    cp install/.bashrc /root/.bashrc
    echo "rpi1" > /etc/hostname
    hostname -F /etc/hostname
    sed -i "s/raspberrypi/rpi1/g" /etc/hosts
    ln -sf /usr/share/zoneinfo/Europe/Madrid /etc/localtime

    echo "##### Installing dependences with apt-get..."
    # Install python libs and zabbix agent
    apt-get -y update && apt-get -y upgrade
    echo iptables-persistent iptables-persistent/autosave_v4 boolean false | debconf-set-selections
    echo iptables-persistent iptables-persistent/autosave_v6 boolean false | debconf-set-selections
    apt-get -y install python-numpy python3-picamera python3-sense-hat zabbix-agent iptables-persistent
    
    echo "##### Enabling components on rpi..."
    # Enable camera and i2c
    raspi-config nonint do_i2c 0    # To use sense hat
    raspi-config nonint do_camera 0 # To use camera

    echo "##### Configuring permissions..."
    usermod -aG input lab # To use sense hat
    usermod -aG i2c lab   # To use sense hat
    usermod -aG video lab # To use camera
    
    echo "##### Preparing daemons and start on boot..."
    # Copy daemons and enable them
    cp install/rpi1/daemons/* /etc/systemd/system
    cd /etc/systemd/system
    for s in `ls cpd_*`; do
        systemctl enable $s
    done
    systemctl enable zabbix-agent
    cd -
    
    echo "##### Configuring zabbix..."
    answ="n"
    while [ $answ != "y" ] && [ $answ != "Y" ]; do
        read -p "Indicate zabbix server address (XXX.XXX.XXX.XXX): " addr
        read -p "Is the IP $addr correct? (Y/n): " answ
        answ=${answ:-Y}
    done
    sed -i "s/Server=127.0.0.1/Server=$addr/g" /etc/zabbix/zabbix_agentd.conf
    sed -i "s/ServerActive=127.0.0.1/ServerActive=$addr/g" /etc/zabbix/zabbix_agentd.conf
    sed -i '/# UserParameter=/a UserParameter=cpd.hum, /bin/cat /tmp/last_hum.txt' /etc/zabbix/zabbix_agentd.conf
    sed -i '/# UserParameter=/a UserParameter=cpd.temp, /bin/cat /tmp/last_temp.txt' /etc/zabbix/zabbix_agentd.conf
    
    echo "##### Setting and saving iptables..."
    set_iptables
    # Allow zabbix in our network
    iptables -A INPUT -p tcp --dport 10050 -s 163.117.170.0/24 -j ACCEPT
    iptables -A INPUT -p tcp --dport 10050 -s 163.117.142.0/24 -j ACCEPT
    # Save rules so they persist after reboot
    iptables-save > /etc/iptables/rules.v4
    
    echo "##### Generating config.json..."
    answ="n"
    while [ $answ != "y" ] && [ $answ != "Y" ]; do
        read -p "Indicate rpi 2 IP address (XXX.XXX.XXX.XXX): " addr
        read -p "Indicate rpi 2 API port (XXXX): " port
        read -p "Indicate rpi 2 API POST Bearer token (XXXXXXXXXXXX): " token
        conf='{\n\t"Rpi2APIAddress" : "'$addr'",\n\t"Rpi2APIPort" : "'$port'",\n\t"Rpi2APIAuthorizedToken" : "Bearer '$token'"\n}'
        echo "config.json generated:"
        echo -e $conf
        read -p "Is that correct? (Y/n): " answ
        answ=${answ:-Y}
    done
    echo -e $conf > config.json
}

function install_rpi2 {
    echo "##### Copying files..."
    # Copy binary and alarm.mp3
    echo "##### Preparing daemons and start on boot..."
    # Copy daemons and enable it
    echo "##### Setting and enabling iptables..."
    # Set and enable iptables
    echo "##### Generating config.json..."
    #create_config_json()
    read -p "##### WARNING: Please, press the philips hue bridge button before continue. Press enter when pressed."
    echo "##### Pairing philips hue bridge..."
}

function install_rpi3 {
    echo "Error: No implemented yet."
    exit 1
}

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   exit 1
fi

echo -e "What rpi are you trying to install?\n"
echo -e "\t(1) rpi1 (inside CPD: monitoring)"
echo -e "\t(2) rpi2 (outside CPD: light control, display CPD info)"
echo -e "\t(3) rpi3 (outside CPD: display classrooms info)\n"
read -p "Please, choose a number. Choose (0) to exit: " opt

case "$opt" in
"1")
    # Install rpi1
    install_rpi1
    ;;
"2")
    # Install rpi2
    install_rpi2
    ;;
"3")
    # Install rpi2
    install_rpi3
    ;;
esac

echo "Done. If you see this message, everything should work after reboot."
echo "Rebooting now..."
#reboot # TODO: enable