#!/usr/bin/env sh

# 批量编译脚本

echo "开始"

path="./services/"

services=("sing" "speak" "listen")
for service in ${services[*]}
do

echo "编译服务${service}..."

fullpath=$path$service
echo "PATH：${fullpath}"
echo ""

go build -o ./bin $fullpath

done

echo "编译客户端..."
fullpath="./client"

echo "PATH：${fullpath}"
echo ""

go build -o ./bin $fullpath

echo "结束"