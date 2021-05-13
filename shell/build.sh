#!/usr/bin/env sh

# Build the Go program script

echo "start"
echo ""

# service code root path
servicePath="./examples/services/"

# client code root path
clientPath="./examples/client"

# services
services=("sing" "speak" "listen")

echo "build service ···"

for service in ${services[*]}
do
fullpath=$servicePath$service
echo "build service - ${service} - PATH：${fullpath}"

go build -o ./bin $fullpath
done

echo ""

echo "build client ···"
echo "PATH：${clientPath}"

go build -o ./bin $clientPath

echo ""

echo "end"