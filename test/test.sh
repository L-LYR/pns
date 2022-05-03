#!/bin/bash

get_k6() {
    # get xk6
    go install go.k6.io/xk6/cmd/xk6@latest
    mkdir -p ./tool
    # build k6 with prometheus
    xk6 build --with github.com/grafana/xk6-output-prometheus-remote --output ./tool/k6
}

test_base() {
    K6_PROMETHEUS_REMOTE_URL=http://localhost:9090/api/v1/write \
        ./tool/k6 run --summary-export="$1_result.json" "./sample/$1.js" -o output-prometheus-remote
}

$1 "$2"
