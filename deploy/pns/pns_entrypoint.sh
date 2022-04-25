#!/bin/sh
apk add netcat-openbsd
# wait for mysql
until nc -z mysql 3306 >/dev/null 2>&1; do
    echo "Waiting......"
    sleep 5
done
/pns/build/pns --gf.gcfg.file=config.container.toml
