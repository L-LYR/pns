#!/bin/sh
apt-get install -y netcat
# wait for mysql
until nc -z mysql 3306 >/dev/null 2>&1; do
    echo "Waiting......"
    sleep 5
done
chown -R mosquitto:mosquitto /mosquitto/data /mosquitto/log
/usr/sbin/mosquitto -c /mosquitto/config/mosquitto.conf
