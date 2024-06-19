FROM php:8.2.20-fpm

# RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list \
RUN apt update \
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
    && pecl install amqp igbinary imagick redis xlswriter \
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
    && cat /usr/share/zoneinfo/Asia/Shanghai  > /etc/localtime \
    && echo 'Asia/Shanghai' > /etc/timezone \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && curl -o /usr/bin/composer https://mirrors.aliyun.com/composer/composer.phar \
    && chmod +x /usr/bin/composer

EXPOSE 9000

WORKDIR /www

CMD ["/usr/local/sbin/php-fpm", "-F"]
