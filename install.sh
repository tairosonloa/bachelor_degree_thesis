#! /bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root. Try using \"sudo $0\""
   exit 1
fi

if [[ $# -ne 1 ]]; then
    echo "This script must be called with the user who will run the software as argument."
    echo -e "\tExample: $0 lab"
    exit 1
fi
INSTALL_USER=$1

function set_iptables {
    # Remove current iptables rules
    iptables -F
    # Default policies
    iptables -P INPUT DROP
    iptables -P OUTPUT ACCEPT
    iptables -P FORWARD DROP
    # Allow active connections
    iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
    # Allow ssh in our network
    iptables -A INPUT -p tcp --dport 22 -s 163.117.170.0/24 -j ACCEPT
    iptables -A INPUT -p tcp --dport 22 -s 163.117.142.0/24 -j ACCEPT
    # Allow loopback
    iptables -I INPUT 1 -i lo -j ACCEPT
    # Allow ping
    iptables -A INPUT -p icmp --icmp-type 8 -j ACCEPT
}

function install_rpi1 {
    echo -e "\t##### Updating hostname, localtime and bashrc..."
    cp install/.bashrc /root/.bashrc
    echo "rpi1" > /etc/hostname
    hostname -F /etc/hostname
    sed -i "s/raspberrypi/rpi1/g" /etc/hosts
    ln -sf /usr/share/zoneinfo/Europe/Madrid /etc/localtime

    echo -e "\t##### Installing dependences with apt-get...\n"
    # Install python libs and zabbix agent
    apt-get -y update && apt-get -y upgrade
    echo iptables-persistent iptables-persistent/autosave_v4 boolean false | debconf-set-selections
    echo iptables-persistent iptables-persistent/autosave_v6 boolean false | debconf-set-selections
    apt-get -y install python-numpy python3-picamera python3-sense-hat zabbix-agent iptables-persistent
    
    echo -e "\n\t##### Installing scripts in /usr/local/bin/rpi1_cpd/..."
    cp -r rpi1/scripts /usr/local/bin/rpi1_cpd
    chown -R $INSTALL_USER:$INSTALL_USER /usr/local/bin/rpi1_cpd

    echo -e "\t##### Enabling components on rpi..."
    # Enable camera and i2c
    raspi-config nonint do_i2c 0    # To use sense hat
    raspi-config nonint do_camera 0 # To use camera

    echo -e "\t##### Configuring permissions..."
    usermod -aG input $INSTALL_USER # To use sense hat
    usermod -aG i2c $INSTALL_USER   # To use sense hat
    usermod -aG video $INSTALL_USER # To use camera
    
    echo -e "\t##### Preparing daemons and start on boot...\n"
    # Copy daemons and enable them
    cp install/rpi1/daemons/* /etc/systemd/system
    cd /etc/systemd/system
    for s in `ls cpd_*`; do
        systemctl enable $s
    done
    systemctl enable zabbix-agent
    cd -
    
    echo -e "\n\t##### Configuring zabbix...\n"
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
    
    echo -e "\n\t##### Generating /etc/rpi1_conf.json...\n"
    answ="n"
    while [ $answ != "y" ] && [ $answ != "Y" ]; do
        read -p "Indicate rpi 2 IP address (XXX.XXX.XXX.XXX): " addr
        read -p "Indicate rpi 2 API port (XXXX): " port
        read -p "Indicate rpi 2 API POST Bearer token (XXXXXXXXXXXX): " token
        conf='{\n\t"Rpi2APIAddress" : "'$addr'",\n\t"Rpi2APIPort" : '$port',\n\t"Rpi2APIAuthorizedToken" : "Bearer '$token'"\n}'
        echo "/etc/rpi1_conf.json generated:"
        echo -e $conf
        read -p "Is that correct? (Y/n): " answ
        answ=${answ:-Y}
    done
    echo -e $conf > /etc/rpi1_conf.json
    chown 600 /etc/rpi1_conf.json
    chown $INSTALL_USER:$INSTALL_USER /etc/rpi1_conf.json

    echo -e "\n\t##### Setting and saving iptables..."
    set_iptables
    # Allow zabbix in our network
    iptables -A INPUT -p tcp --dport 10050 -s 163.117.170.0/24 -j ACCEPT
    iptables -A INPUT -p tcp --dport 10050 -s 163.117.142.0/24 -j ACCEPT
    # Save rules so they persist after reboot
    iptables-save > /etc/iptables/rules.v4
}

function install_rpi2 {
    echo -e "\t##### Updating hostname, localtime and bashrc..."
    cp install/.bashrc /root/.bashrc
    echo "rpi2" > /etc/hostname
    hostname -F /etc/hostname
    sed -i "s/raspberrypi/rpi2/g" /etc/hosts
    ln -sf /usr/share/zoneinfo/Europe/Madrid /etc/localtime

    echo -e "\t##### Installing dependences with apt-get...\n"
    # Install openbox and chromium to display grafana dashboard, lightdm to autologin on GUI, omxplayer to play alarm sound
    apt-get -y update && apt-get -y upgrade
    echo iptables-persistent iptables-persistent/autosave_v4 boolean false | debconf-set-selections
    echo iptables-persistent iptables-persistent/autosave_v6 boolean false | debconf-set-selections
    apt-get -y install iptables-persistent openbox chromium-browser lightdm omxplayer

    echo -e "\n\t##### Installing binary in /usr/local/bin/ and resources in /usr/local/share/..."
    cp install/rpi2/rpi2_api_arm /usr/local/bin/rpi2_api_arm # Rpi2 API binary
    chmod 755 /usr/local/bin/rpi2_api_arm
    cp install/rpi2/alarm.mp3 /usr/local/share/alarm.mp3     # Rpi2 API alarm sound file

    echo -e "\t##### Configuring permissions..."
    usermod -aG video $INSTALL_USER # To use audio jack

    echo -e "\t##### Preparing daemons and start on boot...\n"
    # Copy daemons and enable them
    cp install/rpi2/daemons/* /etc/systemd/system
    cd install/rpi2/daemons/
    for s in `ls .`; do
        systemctl enable $s
    done
    cd -

    echo -e "\n\t##### Enabling auto-login and auto start chromium"
    cp install/rpi2/autostart /etc/xdg/openbox/autostart
    chmod +x /etc/xdg/openbox/autostart
    raspi-config nonint do_boot_behaviour B4 # Auto login with GUI
    sed -i "s/#xserver-command=X/xserver-command=X -nocursor/g" /etc/lightdm/lightdm.conf # Disable mouse on screen

    echo -e "\t##### Preparing monitor auto on/off on working hours..."
    cp install/rpi2/raspi-monitor /usr/local/sbin/raspi-monitor
    chmod +x /usr/local/sbin/raspi-monitor
    # Set cron jobs
    (crontab -l 2>/dev/null; echo "# Enable the monitor every weekday morning at 8:10") | crontab -
    (crontab -l 2>/dev/null; echo "10 8 * * 1-5 /usr/local/sbin/raspi-monitor on > /dev/null 2>&1") | crontab -
    (crontab -l 2>/dev/null; echo "# Disable the monitor every weekday evening at 21:10") | crontab -
    (crontab -l 2>/dev/null; echo "10 21 * * 1-5 /usr/local/sbin/raspi-monitor off > /dev/null 2>&1") | crontab -

    echo -e "\t##### Generating /etc/rpi2_conf.json...\n"
    answ="n"
    while [ $answ != "y" ] && [ $answ != "Y" ]; do
        read -p "Indicate rpi 2 IP address (XXX.XXX.XXX.XXX): " addr
        read -p "Indicate rpi 2 API port (XXXX): " port
        read -p "Indicate rpi 2 API POST Bearer token (XXXXXXXXXXXX): " token
        read -p "Indicate Philips Hue bridge IP address (XXX.XXX.XXX.XXX): " hue
        read -p "Indicate Philips Hue bridge secret string (XXXXXXXXXXXX): " secret
        echo "Adding the alarm sound file path to the config file..."
        conf='{\n\t"Rpi2APIAddress" : "'$addr'",\n\t"Rpi2APIPort" : '$port',\n\t"Rpi2APIAuthorizedToken" : "Bearer '$token'",\n\t"HueBridgeAddress" : "'$hue'",\n\t"HueBridgeToken" : "'$secret'",\n\t"AlarmSoundPath" : "/usr/local/share/alarm.mp3"\n}'
        echo "/etc/rpi2_conf.json generated:"
        echo -e $conf
        read -p "Is that correct? (Y/n): " answ
        answ=${answ:-Y}
    done
    echo -e $conf > /etc/rpi2_conf.json
    chown 600 /etc/rpi2_conf.json
    chown $INSTALL_USER:$INSTALL_USER /etc/rpi2_conf.json

    echo -e "\n\t##### Setting and enabling iptables..."
    set_iptables
    # Allow API requests only on the university network
    iptables -A INPUT -p tcp --dport $port -s 163.117.0.0/16 -j ACCEPT
    # Save rules so they persist after reboot
    iptables-save > /etc/iptables/rules.v4

    echo -e "\t##### WARNING: Please, press the Philips Hue bridge button before continue. Press enter when pressed."
    read
    echo -e "\t##### Launching api to pair Philips Hue bridge..."
    systemctl start rpi2_api.service
    echo -e "\t      Waiting 10 seconds"
    sleep 10
    systemctl is-active --quiet rpi2_api.service
    if [[ $? -eq 0 ]]; then
        echo -e "\t      Pairing successful!"
    else
        echo -e "\t      Pairing failed!"
        exit 2
    fi
}

function install_rpi3 {
    echo -e "\t##### Updating hostname, localtime and bashrc..."
    cp install/.bashrc /root/.bashrc
    echo "rpi3" > /etc/hostname
    hostname -F /etc/hostname
    sed -i "s/raspberrypi/rpi3/g" /etc/hosts
    ln -sf /usr/share/zoneinfo/Europe/Madrid /etc/localtime

    echo -e "\t##### Installing dependences with apt-get...\n"
    # Install openbox and chromium to display grafana dashboard, lightdm to autologin on GUI, omxplayer to play alarm sound
    apt-get -y update && apt-get -y upgrade
    echo iptables-persistent iptables-persistent/autosave_v4 boolean false | debconf-set-selections
    echo iptables-persistent iptables-persistent/autosave_v6 boolean false | debconf-set-selections
    apt-get -y install iptables-persistent

    echo -e "\n\t##### Installing binary in /usr/local/bin/ and resources in /usr/local/share/..."
    cp install/rpi3/rpi3_api_arm /usr/local/bin/rpi3_api_arm # Rpi3 API binary
    chmod 755 /usr/local/bin/rpi3_api_arm

    echo -e "\t##### Preparing daemons and start on boot...\n"
    # Copy daemons and enable them
    cp install/rpi3/daemons/* /etc/systemd/system
    cd install/rpi3/daemons/
    for s in `ls .`; do
        systemctl enable $s
    done
    cd -

    echo -e "\t##### Generating /etc/rpi3_conf.json...\n"
    answ="n"
    while [ $answ != "y" ] && [ $answ != "Y" ]; do
        read -p "Indicate rpi 3 IP address (XXX.XXX.XXX.XXX): " addr
        read -p "Indicate rpi 3 API port (XXXX): " port
        conf='{\n\t"Rpi3APIAddress" : "'$addr'",\n\t"Rpi3APIPort" : '$port'\n}'
        echo "/etc/rpi3_conf.json generated:"
        echo -e $conf
        read -p "Is that correct? (Y/n): " answ
        answ=${answ:-Y}
    done
    echo -e $conf > /etc/rpi3_conf.json
    chown 600 /etc/rpi3_conf.json
    chown $INSTALL_USER:$INSTALL_USER /etc/rpi3_conf.json

    echo -e "\n\t##### Setting and enabling iptables..."
    set_iptables
    # Allow API requests only on the university network
    iptables -A INPUT -p tcp --dport $port -s 163.117.0.0/16 -j ACCEPT
    # Save rules so they persist after reboot
    iptables-save > /etc/iptables/rules.v4
}

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
reboot