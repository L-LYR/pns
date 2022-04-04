#!/bin/bash
function gen() {
    protoc --go_out="./pkg/" "$1" || exit
}

for file in ./proto/*.proto; do
    gen "$file"
done
