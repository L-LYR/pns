#!/bin/sh

# TODO: check dependencies
# TODO: run as superdoer

# initialize basic global variables
root_path=$(realpath "$0" | xargs dirname | xargs dirname)
if [ -z "${DEBUG}" ]; then
    deploy_path="$root_path/deploy"
else
    deploy_path="$root_path/.test_deploy"
fi
docker_file_dir="$root_path/manifest/docker"
TZ="$(cat /etc/timezone)" || exit

# define and export env variables
export db_root_pass="pns_root"
export pns_mongo_volume="$deploy_path/pns_mongo"
export pns_mysql_volume="$deploy_path/pns_mysql"
export pns_redis_volume="$deploy_path/pns_redis"
export pns_prometheus_volume="$deploy_path/pns_prometheus"
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
    cp -a "$docker_file_dir/." "$deploy_path" || exit
}

up() {
    printf "Making directory 'deploy' in working directory...\n"
    mkdir -p "$deploy_path" || exit
    update_configs
    goto_deploy_directory

    printf "Deploying...\n"
    printf "Making directories for database volumes...\n"
    mkdir -p "$pns_mongo_volume" || exit
    mkdir -p "$pns_mysql_volume" || exit
    mkdir -p "$pns_redis_volume" || exit
    mkdir -p "$pns_prometheus_volume" || exit
    sudo chown -R 65534:65534 "$pns_prometheus_volume" || exit
    mkdir -p "$pns_grafana_volume" || exit
    sudo chown -R 472:472 "$pns_grafana_volume" || exit

    mkdir -p "$pns_grafana_provisioning" || exit
    printf "Bootstrap...\n"
    docker-compose up -d || exit
    printf "Serving...\n"
    # TODO: add serving url
}

down() {
    goto_deploy_directory

    printf "Shutdown...\n"
    docker-compose down || exit
    printf "Change to working directory\n"
    printf "Cleanup...\n"
    sudo rm -rf "$deploy_path"
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

$1
