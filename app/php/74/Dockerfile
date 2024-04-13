FROM php:7.4-fpm

RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list \
    && apt update \
    && apt install -y \
        procps \
        librabbitmq-dev \
        libbz2-dev \
        libfreetype6-dev \
        libjpeg62-turbo-dev \
        libpng-dev \
        libmagickwand-dev \
        libmagickcore-dev \
        libzip-dev \
        libpq-dev \
        supervisor \
    && pecl install amqp igbinary imagick redis \
    && docker-php-ext-enable amqp igbinary imagick redis \
    && docker-php-ext-install -j$(nproc) \
        bcmath \
        bz2 \
        exif \
        gd \
        gettext \
        pcntl \
        pdo_mysql \
        sockets \
        opcache \
        zip \
        pgsql \
        pdo_pgsql \
    && cd /usr/src/php/ext \
    && curl -o swoole.tgz https://pecl.php.net/get/swoole-4.8.6.tgz \
    && tar -zxvf swoole.tgz \
    && rm swoole.tgz \
    && docker-php-ext-install -j$(nproc) swoole-4.8.6 \
    && cat /usr/share/zoneinfo/Asia/Shanghai  > /etc/localtime \
    && echo 'Asia/Shanghai' > /etc/timezone \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 9000

WORKDIR /www

CMD ["/usr/local/sbin/php-fpm", "-F"]

