# 基于 Debian 构建 PHP
FROM php:7.4-fpm

# 替换 apt 镜像源 && 更新
RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list \
    && apt update \
    && apt install -y procps \
    && usermod -u 1000 www-data \
    && groupmod -g 1000 www-data \
    && mkdir -p /var/log/php7 \
    && chmod -R 777 /var/log/php7 \
    && curl -o /usr/bin/composer https://mirrors.aliyun.com/composer/composer.phar \
    && chmod +x /usr/bin/composer

# 自定义启用的扩展
# pecl install or source

# fix: ERROR: `/tmp/pear/temp/amqp/configure --with-php-config=/usr/local/bin/php-config --with-librabbitmq-dir' failed
RUN apt install -y librabbitmq-dev \
    # fix: error: Please reinstall the BZip2 distribution
    libbz2-dev \
    # fix: Package requirements (zlib) were not met:
    libfreetype6-dev \
    libjpeg62-turbo-dev \
    libpng-dev \
    # fix: --with-imagick failed
    libmagickwand-dev \
    libmagickcore-dev \
    # fix: No package 'libzip' found
    libzip-dev \
    && pecl install amqp \
    && docker-php-ext-enable amqp.so \
    && docker-php-ext-install -j$(nproc) bcmath \
    && docker-php-ext-install -j$(nproc) bz2 \
    && docker-php-ext-install -j$(nproc) exif \
    && docker-php-ext-configure gd --with-freetype --with-jpeg \
    && docker-php-ext-install -j$(nproc) gd \
    && docker-php-ext-install -j$(nproc) gettext \
    && pecl install igbinary \
    && docker-php-ext-enable igbinary \
    && pecl install imagick \
    && docker-php-ext-enable imagick.so \
    && docker-php-ext-install -j$(nproc) pcntl \
    && docker-php-ext-install -j$(nproc) pdo_mysql \
    && pecl install redis \
    && docker-php-ext-enable redis.so \
    && docker-php-ext-install -j$(nproc) sockets \
    && docker-php-ext-install -j$(nproc) opcache \
    && docker-php-ext-install -j$(nproc) zip

# 安装 swoole
# version: 4.8.6
COPY ./extension/swoole-4.8.6.tgz /usr/src/php/ext/swoole-4.8.6.tgz

RUN cd /usr/src/php/ext \
    && tar -zxvf /usr/src/php/ext/swoole-4.8.6.tgz \
    && docker-php-ext-install -j$(nproc) swoole-4.8.6 \
    # fix: set timezone Shanghai
    && cat /usr/share/zoneinfo/Asia/Shanghai  > /etc/localtime \
    && echo 'Asia/Shanghai' > /etc/timezone

# 安装 supervisor
RUN apt install -y supervisor

# 增加 pgsql
RUN apt install -y libpq-dev \
    && docker-php-ext-install -j$(nproc) pgsql \
    && docker-php-ext-install -j$(nproc) pdo_pgsql

EXPOSE 9000
# 定义 php-fpm 工作目录
WORKDIR /www

# php-fpm 执行文件目录 /usr/local/sbin/php-fpm
CMD ["/usr/local/sbin/php-fpm", "-F"]
