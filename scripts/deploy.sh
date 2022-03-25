#!/bin/bash

# TODO: check dependencies
# TODO: run as superdoer

# initialize basic global variables
root_path=$(realpath "$0" | xargs dirname | xargs dirname)
if [ -z "${DEBUG}" ]; then
    deploy_path="$root_path/.deploy"
    pns_dockerfile="${deploy_path}/dockerfile"
else
    deploy_path="$root_path/.test_deploy"
    pns_dockerfile="${deploy_path}/debug.dockerfile"
fi
deploy_config_path="$root_path/deploy"
build_path="$root_path/build"
TZ="$(cat /etc/timezone)" || exit
source_dir=("cmd" "assets" "config" "docs" "internal" "web")

# define and export env variables
export root_path
export deploy_path
export build_path
export pns_dockerfile
export db_root_pass="pns_root"
export pns_image_name="hammerli/pns:v1"
export pns_log_volume="$deploy_path/pns_log"
export pns_mongo_volume="$deploy_path/pns_mongo"
export pns_mysql_volume="$deploy_path/pns_mysql"
export pns_redis_volume="$deploy_path/pns_redis"
export pns_prometheus_volume="$deploy_path/pns_prometheus"
export pns_mosquitto_volume="$deploy_path/pns_mosquitto"
export pns_grafana_volume="$deploy_path/pns_grafana"
export pns_grafana_provisioning="$deploy_path/pns_grafana_provisioning"
export TZ

# print basic infomation
printf "Working directory: %s\n" "$(pwd)"
printf "Root path of pns: %s\n" "$root_path"
printf "Timezone: %s\n" "$TZ"

goto_deploy_directory() {
    printf "Change to deploy directory\n"
    cd "$deploy_path" || exit
}

update_configs() {
    printf "Copying all deployment configs...\n"
    cp -a "$deploy_config_path/." "$deploy_path" || exit
}

up() {
    printf "Making directory 'deploy' in working directory...\n"
    mkdir -p "$deploy_path" || exit
    update_configs
    goto_deploy_directory

    printf "Deploying...\n"
    printf "Making directories for database volumes...\n"
    mkdir -p "$pns_log_volume" || exit
    mkdir -p "$pns_mongo_volume" || exit
    mkdir -p "$pns_mysql_volume" || exit
    mkdir -p "$pns_redis_volume" || exit
    mkdir -p "$pns_prometheus_volume" || exit
    mkdir -p "$pns_mosquitto_volume" || exit
    sudo chown -R 65534:65534 "$pns_prometheus_volume" || exit
    mkdir -p "$pns_grafana_volume" || exit
    sudo chown -R 472:472 "$pns_grafana_volume" || exit

    mkdir -p "$pns_grafana_provisioning" || exit
    printf "Bootstrap...\n"
    if [ -z "$(docker images -q $pns_image_name)" ]; then
        docker-compose build || exit
    fi
    docker-compose up -d || exit
    printf "Serving...\n"
    printf "Monitor: localhost:3000\n"
    # TODO: add other urls
}

down() {
    goto_deploy_directory

    printf "Shutdown...\n"
    docker-compose down
    docker rmi "$pns_image_name"
    printf "Change to working directory\n"
    if [ "${DEBUG}" ]; then
        printf "Cleanup...\n"
        sudo rm -rf "$deploy_path"
    fi
}

stop() {
    goto_deploy_directory

    printf "Shutdown...\n"
    docker-compose stop || exit
}

start() {
    update_configs
    goto_deploy_directory

    printf "Starting...\n"
    docker-compose start || exit
}

update() {
    if [ -z "${DEBUG}" ]; then
        exit
    fi

    goto_deploy_directory
    printf "Updating...\n"
    if [ -z "$(docker container ls --format "{{.Status}} {{.Names}}" |
        awk '{if ($NF == "pns" && $1 == "Up") {print "pns is working\n"}}')" ]; then
        printf "pns is not working\n"
        printf "Try to remove container and image, and then rebuild and restart...\n"
        docker stop pns
        docker rm pns
        docker rmi ${pns_image_name}
        docker-compose build || exit
        docker-compose up -d || exit
        exit
    fi
    for dir in "${source_dir[@]}"; do
        docker cp "$root_path"/"$dir" pns:/pns/
    done
    docker exec -it pns /bin/sh -c make all
    docker exec -it pns mv -f /pns/build/pns -t /pns
    docker restart pns
}

$1
