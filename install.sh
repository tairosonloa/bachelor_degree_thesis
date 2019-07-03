#! /bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root. Try using \"sudo $0\""
   exit 1
fi

INSTALL_USER=lab

function host_configuration {
    echo -e "\t##### Updating hostname, localtime, locale and bashrc..."
    cp install/.bashrc /root/.bashrc
    echo $1 > /etc/hostname
    hostname -F /etc/hostname
    sed -i "s/raspberrypi/$1/g" /etc/hosts
    ln -sf /usr/share/zoneinfo/Europe/Madrid /etc/localtime
    sed -i "s/# es_ES.UTF-8/es_ES.UTF-8/g" /etc/locale.gen
    locale-gen
    sed -i "s/en_GB.UTF-8/es_ES.UTF-8/g" /etc/default/locale
    echo -e "\t##### Creating user $INSTALL_USER..."
    useradd -m -s /bin/bash $INSTALL_USER
    passwd $INSTALL_USER
    while [ $? -ne 0 ]; do
        passwd $INSTALL_USER
    done
    usermod -aG sudo $INSTALL_USER
}

function install_dependencies {
    echo -e "\t##### Installing dependencies with apt-get...\n"
    apt-get -y update && apt-get -y upgrade
    echo iptables-persistent iptables-persistent/autosave_v4 boolean false | debconf-set-selections
    echo iptables-persistent iptables-persistent/autosave_v6 boolean false | debconf-set-selections
    apt-get -y install iptables-persistent $@
}

function daemons_setup  {
    echo -e "\t##### Preparing daemons and start on boot...\n"
    # Copy daemons and enable them
    cd install/$1/daemons/
    for s in `ls .`; do
        sed -i "s/<user>/$INSTALL_USER/g" $s
        sed -i "s/<group>/$INSTALL_USER/g" $s
        cp $s /etc/systemd/system
        systemctl enable $s
    done
    cd -
    systemctl enable ssh
}

function set_config_json {
    echo -e "\n\t##### Generating /etc/$1_conf.json...\n"
    answ="n"
    # Load default config values
    . install/default.config
    echo "NOTE: Default values are shown between parenthesis. If you don't give a value, default value will be taken."
    while [ $answ != "y" ] && [ $answ != "Y" ]; do
        read -p "Indicate rpi 2 IP address ($addr2_default): " addr2
        if [ -z "$addr2" ]; then addr2=$addr2_default; fi
        read -p "Indicate rpi 2 API port ($port2_default): " port2
        if [ -z "$port2" ]; then port2=$port2_default; fi
        case $1 in
        "rpi1")
            read -p "Indicate rpi 2 API POST Bearer token ($token2_default): " token2
            if [ -z "$token2" ]; then token2=$token2_default; fi
            conf='{\n\t"Rpi2APIAddress" : "'$addr2'",\n\t"Rpi2APIPort" : '$port2',\n\t"Rpi2APIAuthorizedToken" : "Bearer '$token2'"\n}'
            ;;
        "rpi2")
            read -p "Indicate rpi 2 API POST Bearer token ($token2_default): " token2
            if [ -z "$token2" ]; then token2=$token2_default; fi
            read -p "Indicate Philips Hue bridge IP address ($hue_default): " hue
            if [ -z "$hue" ]; then hue=$hue_default; fi
            read -p "Indicate Philips Hue bridge secret string ($secret_default): " secret
            if [ -z "$secret" ]; then secret=$secret_default; fi
            echo "Adding the alarm sound file path to the config file..."
            conf='{\n\t"Rpi2APIAddress" : "'$addr2'",\n\t"Rpi2APIPort" : '$port2',\n\t"Rpi2APIAuthorizedToken" : "Bearer '$token2'",\n\t"HueBridgeAddress" : "'$hue'",\n\t"HueBridgeToken" : "'$secret'",\n\t"AlarmSoundPath" : "/usr/local/share/alarm.mp3"\n}'
            ;;
        "rpi3")
            read -p "Indicate rpi 3 IP address ($addr3_default): " addr3
            if [ -z "$addr3" ]; then addr3=$addr3_default; fi
            read -p "Indicate rpi 3 API port ($port3_default): " port3
            if [ -z "$port3" ]; then port3=$port3_default; fi
            read -p "Indicate classrooms control server domain name and ssh port ($server_default): " server
            if [ -z "$server" ]; then server=$server_default; fi
            read -p "Indicate check classrooms occupation command ($cmd_default): " cmd
            if [ -z "$cmd" ]; then cmd=$cmd_default; fi
            read -p "Indicate laboratory reservations web page url ($web_default): " web
            if [ -z "$web" ]; then web=$web_default; fi
            conf='{\n\t"Rpi2APIAddress" : "'$addr2'",\n\t"Rpi2APIPort" : '$port2',\n\t"Rpi3APIAddress" : "'$addr3'",\n\t"Rpi3APIPort" : '$port3',\n\t"ControlServer" : "'$server'",\n\t"OccupationCmd" : "'$cmd'",\n\t"OccupationWeb" : "'$web'"\n}'
            ;;
        esac
        echo "/etc/$1_conf.json generated:"
        echo -e $conf
        read -p "Is that correct? (Y/n): " answ
        answ=${answ:-Y}
    done
    echo -e $conf > /etc/$1_conf.json
    chown 600 /etc/$1_conf.json
    chown $INSTALL_USER:$INSTALL_USER /etc/$1_conf.json
}

function auto_login_gui {
    echo -e "\n\t##### Enabling auto-login and auto start chromium"
    cp install/$1/autostart /etc/xdg/openbox/autostart
    chmod +x /etc/xdg/openbox/autostart
    raspi-config nonint do_boot_behaviour B4 # Auto login with GUI
    sed -i "s/#xserver-command=X/xserver-command=X -nocursor/g" /etc/lightdm/lightdm.conf # Disable mouse on screen
}

function auto_power_monitor {
    echo -e "\t##### Preparing monitor auto on/off on working hours..."
    cp install/raspi-monitor /usr/local/sbin/raspi-monitor
    chmod +x /usr/local/sbin/raspi-monitor
    # Set cron jobs
    (crontab -l 2>/dev/null; echo "# Enable the monitor every weekday morning at 8:10") | crontab -
    (crontab -l 2>/dev/null; echo "10 8 * * 1-5 /usr/local/sbin/raspi-monitor on > /dev/null 2>&1") | crontab -
    (crontab -l 2>/dev/null; echo "# Disable the monitor every weekday evening at 21:10") | crontab -
    (crontab -l 2>/dev/null; echo "10 21 * * 1-5 /usr/local/sbin/raspi-monitor off > /dev/null 2>&1") | crontab -
}

function set_monitor_resolution {
    echo -e "\n\t##### Setting monitor resolution..."
    sed -i "s/#disable_overscan=1/disable_overscan=1/g" /boot/config.txt
    sed -i "s/#overscan_left=16/overscan_left=0/g" /boot/config.txt
    sed -i "s/#overscan_right=16/overscan_right=0/g" /boot/config.txt
    sed -i "s/#overscan_top=16/overscan_top=0/g" /boot/config.txt
    sed -i "s/#overscan_bottom=16/overscan_bottom=0/g" /boot/config.txt
    sed -i "s/#framebuffer_width=1280/framebuffer_width=1920/g" /boot/config.txt
    sed -i "s/#framebuffer_height=720/framebuffer_height=1080/g" /boot/config.txt
}

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
    # Basic host configuration
    host_configuration "rpi1"

    # Install python libs and zabbix agent
    packages="python-numpy python3-picamera python3-sense-hat zabbix-agent"
    install_dependencies $packages
    
    # Install the core
    echo -e "\n\t##### Installing scripts in /usr/local/bin/rpi1_cpd/..."
    cp -r $1/scripts /usr/local/bin/$1_cpd
    chown -R $INSTALL_USER:$INSTALL_USER /usr/local/bin/$1_cpd

    echo -e "\t##### Enabling components on rpi..."
    # Enable camera and i2c
    raspi-config nonint do_i2c 0    # To use sense hat
    raspi-config nonint do_camera 0 # To use camera

    echo -e "\t##### Configuring permissions..."
    usermod -aG input $INSTALL_USER # To use sense hat
    usermod -aG i2c $INSTALL_USER   # To use sense hat
    usermod -aG video $INSTALL_USER # To use camera
    
    # Setup daemons and start on boot
    daemons_setup "rpi1"
    systemctl enable zabbix-agent
    
    # Configure zabbix
    echo -e "\n\t##### Configuring zabbix...\n"
    answ="n"
    . install/default.config
    while [ $answ != "y" ] && [ $answ != "Y" ]; do
        read -p "Indicate zabbix server address ($addrZ_default): " addrZ
        if [ -z "$addrZ" ]; then addrZ=$addrZ_default; fi
        read -p "Is the IP $addrZ correct? (Y/n): " answ
        answ=${answ:-Y}
    done
    sed -i "s/Server=127.0.0.1/Server=$addrZ/g" /etc/zabbix/zabbix_agentd.conf
    sed -i "s/ServerActive=127.0.0.1/ServerActive=$addrZ/g" /etc/zabbix/zabbix_agentd.conf
    sed -i '/# UserParameter=/a UserParameter=cpd.hum, /bin/cat /tmp/last_hum.txt' /etc/zabbix/zabbix_agentd.conf
    sed -i '/# UserParameter=/a UserParameter=cpd.temp, /bin/cat /tmp/last_temp.txt' /etc/zabbix/zabbix_agentd.conf

    # Configure core
    set_config_json "rpi1"

    echo -e "\n\t##### Setting and saving iptables..."
    set_iptables
    # Allow zabbix in our network
    iptables -A INPUT -p tcp --dport 10050 -s 163.117.170.0/24 -j ACCEPT
    iptables -A INPUT -p tcp --dport 10050 -s 163.117.142.0/24 -j ACCEPT
    # Save rules so they persist after reboot
    iptables-save > /etc/iptables/rules.v4
}

function install_rpi2 {
    # Basic host configuration
    host_configuration "rpi2"

    # Install openbox and chromium to display grafana dashboard, lightdm to autologin on GUI, omxplayer to play alarm sound
    packages="openbox chromium-browser lightdm omxplayer"
    install_dependencies $packages
    
    # Install the core
    echo -e "\n\t##### Installing binary in /usr/local/bin/ and resources in /usr/local/share/..."
    cp install/rpi2/rpi2_api_arm /usr/local/bin/rpi2_api_arm # Rpi2 API binary
    chmod 755 /usr/local/bin/rpi2_api_arm
    cp install/rpi2/alarm.mp3 /usr/local/share/alarm.mp3     # Rpi2 API alarm sound file

    echo -e "\t##### Configuring permissions..."
    usermod -aG video $INSTALL_USER # To use audio jack

    # Setup daemons and start on boot
    daemons_setup "rpi2"

    # Enable auto login and chromium start
    auto_login_gui "rpi2"

    # Enable monitor auto power off/on
    auto_power_monitor

    # Configure core
    set_config_json "rpi2"

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
        echo -e "\t      Pairing failed! Manual pairing required"
        exit 2
    fi
}

function install_rpi3 {
    # Basic host configuration
    host_configuration "rpi3"

    # Install openbox and chromium to display website, lightdm to autologin on GUI, npm to build
    packages="openbox chromium-browser lightdm npm"
    install_dependencies $packages

    echo -e "\n\t##### Installing binary in /usr/local/bin/ and web in /srv/rpi3/..."
    cp install/rpi3/rpi3_api_arm /usr/local/bin/rpi3_api_arm # Rpi3 API binary
    chmod 755 /usr/local/bin/rpi3_api_arm
    cp install/rpi3/web_server_arm /usr/local/bin/web_server_arm # Custom web server
    chmod 755 /usr/local/bin/web_server_arm
    cd rpi3/GUI/
    npm install
    npm run build
    mv public/ /srv/rpi3 # Website files
    cd -
    chown -R $INSTALL_USER:$INSTALL_USER /srv/rpi3

    # Setup daemons and start on boot
    daemons_setup "rpi3"

    # Enable auto login and chromium start
    auto_login_gui "rpi3"

    # Enable monitor auto power off/on
    auto_power_monitor

    # Set monitor resolution
    set_monitor_resolution

    # Configure core
    set_config_json "rpi3"

    echo -e "\n\t##### Setting and enabling iptables..."
    set_iptables
    # Allow API requests only on the university network
    iptables -A INPUT -p tcp --dport $port3 -s 163.117.0.0/16 -j ACCEPT
    # Allow web request only on our subnet
    iptables -A INPUT -p tcp --dport 9000 -s 163.117.142.0/24 -j ACCEPT
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
    # Install rpi3
    install_rpi3
    .install/default.config
    echo "Please. Generate an ssh key to rpi3 and copy pub key into $server_default:/root/.ssh/authorized_keys"
    ;;
esac
echo "Done. If you see this message, everything should work after reboot."
echo "REMEMBER TO REMOVE PI USER FOR SECURITY."
answ="n"
read -p "Do you want to reboot now? (Y/n): " answ
answ=${answ:-Y}
if [ $answ == "y" ] || [ $answ == "Y" ]; then
    reboot
fi