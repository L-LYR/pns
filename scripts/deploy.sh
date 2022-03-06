#!/bin/sh

# TODO: check dependencies
root_path=$(realpath "$0" | xargs dirname | xargs dirname)
export db_root_pass="pns_root"
printf "Working directory: %s\n" "$(pwd)"
printf "Root path of pns: %s\n" "$root_path"
deploy_path="$root_path/deploy"
docker_file_dir="$root_path/manifest/docker"
export pns_mongo_volume="$deploy_path/pns_mongo"
export pns_mysql_volume="$deploy_path/pns_mysql"
export pns_redis_volume="$deploy_path/pns_redis"
TZ="$(cat /etc/timezone)" || exit
printf "Timezone: %s\n" "$TZ"
export TZ

up() {
    printf "Making directory 'deploy' in working directory...\n"
    mkdir -p "$deploy_path" || exit
    printf "Copying all deployment configs...\n"
    cp -a "$docker_file_dir/." "$deploy_path" || exit
    printf "Change to deploy directory\n"
    cd "$deploy_path" || exit
    printf "Deploying...\n"
    printf "Making directories for database volumes...\n"
    mkdir -p "$pns_mongo_volume" || exit
    mkdir -p "$pns_mysql_volume" || exit
    mkdir -p "$pns_redis_volume" || exit
    printf "Bootstrap...\n"
    docker-compose up -d || exit
    printf "Serving...\n"
}

down() {
    printf "Change to deploy directory...\n"
    cd "$deploy_path" || exit
    printf "Shutdown...\n"
    docker-compose down || exit
    printf "Change to working directory\n"
    printf "Cleanup...\n"
    sudo rm -rf "$deploy_path"
}

$1
