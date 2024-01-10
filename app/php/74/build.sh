#!/bin/bash
docker build --network=host \
 --build-arg PHP_VERSION=php:7.4-fpm \
 --build-arg DEBIAN_ORIGIN_REPOSITORIES=deb.debian.org \
 --build-arg DEBIAN_TARGET_REPOSITORIES=mirrors.ustc.edu.cn \
 -t yinshiyi/php74:3.1 .
