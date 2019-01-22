#! /bin/bash

function install_rpi1 {
    echo "##### Installing dependences with apt-get..."
    # Install python libs and zabbix agent
    apt-get -y update && apt-get -y upgrade
    apt-get -y install python-numpy python3-picamera python3-sense-hat zabbix-agent iptables-persistent
    echo "##### Enabling components on rpi..."
    # Enable camera and i2c
    . raspi-config nonint
    do_i2c 0
    do_camera 0
    echo "##### Preparing daemons and start on boot..."
    # Copy daemons and enable them
    cp rpi1/daemons/* /etc/systemd/system
    cd /etc/systemd/system
    for s in `ls cpd_*`; do
        systemctl enable $s
    done
    systemctl enable zabbix-agent
    cd -
    echo "##### Configuring zabbix..."
    # TODO:
    echo "##### Setting and saving iptables..."
    # Remove current iptables rules
    iptables -F
    # Default policies
    iptables -P INPUT DROP
    iptables -P OUTPUT ACCEPT
    iptables -P FORWARD DROP
    # Allow active connections
    iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
    # Allow ssh and zabbix in our network
    iptables -A INPUT -p tcp --dport 22 -s 163.117.170.0/24 -j ACCEPT
    iptables -A INPUT -p tcp --dport 10050 -s 163.117.170.0/24 -j ACCEPT
    iptables -A INPUT -p tcp --dport 22 -s 163.117.142.0/24 -j ACCEPT
    iptables -A INPUT -p tcp --dport 10050 -s 163.117.142.0/24 -j ACCEPT
    # Allow loopback
    iptables -I INPUT 1 -i lo -j ACCEPT
    # Allow ping
    iptables -A INPUT -p icmp --icmp-type 8 -j ACCEPT
    # Save rules so they persist after reboot
    iptables-save > /etc/iptables/rules.v4
    echo "##### Generating config.json..."
    #create_config_json()
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

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   exit 1
fi

echo -e "What rpi are you trying to install?\n\n\t(1) rpi1\n\t(2) rpi2\n"
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
esac

echo "Done. If you see this message, everything should work after reboot."
echo "Rebooting now..."
reboot