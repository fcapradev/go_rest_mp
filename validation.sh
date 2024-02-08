#!/bin/bash

set -e

# create a function to print usage message
usage() {
    echo "usage: ./validation.sh <app_name> <group_number:1|2|3|4|5|6|7|8> <env:local|fury> <case:1|2|3|4>"
    exit 1
}

mb_container_name=mountebank
setup() {
    docker run -p 2525:2525 -p 10000:10000 -d --name ${mb_container_name} andyrbell/mountebank
    while ! nc -z localhost 2525 ; do
        echo "waiting for mountebank to start..."
        sleep 0.25
    done
    sleep 2
    curl -XPOST localhost:2525/imposters --data @mb/imposter.json -o /dev/null --silent
    go build -o app cmd/api/main.go
    ./app > app.log 2>&1 &
    while ! nc -z localhost 8080 ; do
        echo "waiting for mountebank to start..."
        sleep 0.25
    done
}

tearDown() {
    docker stop ${mb_container_name} || echo "mountebank container is not running"
    docker rm ${mb_container_name} || echo "mountebank container does not exists"
    while nc -z localhost 8080; do
        lsof -i :8080 | tail -n1 | awk '{print $2}' | xargs kill
        sleep 0.1
    done
    rm -f app
}

if [[ -z $1 ]]; then
    usage
fi

if [[ -z $2 ]]; then
    usage
fi

if [[ -z $3 ]]; then
    usage
fi

if [[ -z $4 ]]; then
    usage
fi


app=$1
group=$2
env=$3
case_number=$4

echo "running validation for"
echo "- app: $app"
echo "- group: $group"
echo "- env: $env"
echo "- case: $case_number"

export APP_NAME=$app
export GROUP_NUMBER=$group
export ENV=$env


# validate if the env is local
if [[ $env == "fury" ]]; then
    fury get-token
    get_token=$(cat $HOME/.fury/.credentials | jq -r '."tiger-token"')
    export FURY_TOKEN=$get_token
    k6 run "validations/case_${case_number}.js"
else
    tearDown
    setup
    trap tearDown SIGINT
    k6 run "validations/case_${case_number}.js"
    tearDown
fi

