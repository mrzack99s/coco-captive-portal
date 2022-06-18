#!/bin/bash

DIR="/var/coco-captive-portal"

command_exists() {
	command -v "$@" > /dev/null 2>&1
}

if (( $EUID != 0 )); then
    echo "this script needs the ability to run commands as root. We are unable to find either \"sudo\" or \"su\" available to make this happen."
    exit
fi

# Detect OS
lsb_dist=""
if [ -r /etc/os-release ]; then
	lsb_dist="$(. /etc/os-release && echo "$ID")"
fi

if [[ "$lsb_dist" -ne "ubuntu" ]]; then
    echo "this script only support ubuntu"
    exit
fi

# Detect Version
if command_exists lsb_release; then
	dist_version="$(lsb_release --codename | cut -f2)"
fi
if [ -z "$dist_version" ] && [ -r /etc/lsb-release ]; then
	dist_version="$(. /etc/lsb-release && echo "$DISTRIB_CODENAME")"
fi

if [[ "$dist_version" -ne "focal" || "$lsb_dist" -ne "bionic" ]]; then
    echo "this script only support ubuntu focal(20.04.x) and ubuntu bionic(18.04.x)"
    exit
fi

if [ ! -d "$DIR" ] 
then
    mkdir -p $DIR
    chmod 744 $DIR
fi

# Install
apt install -y libpcap0.8 
curl -fsSL https://packages.redis.io/gpg | gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg

echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/redis.list

apt-get update
apt-get install -y redis


# curl -L https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/coco_linux_amd64 -o $DIR/coco
# curl -L https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/coco-captive-portal.service -o /etc/systemd/system/coco-captive-portal.service
# curl -L https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/config.yaml.sample -o $DIR/config.yaml.sample
# curl -L https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/ui.tar.gz -o /tmp/install-coco-captive-portal/ui.tar.gz

