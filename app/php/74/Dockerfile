# 第一阶段：构建阶段
FROM php:7.4-fpm AS builder

# 镜像源加速
# sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
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
        curl \
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
    && mkdir -p /usr/src/php/ext \
    && cd /usr/src/php/ext \
    && curl -o swoole.tgz https://pecl.php.net/get/swoole-4.8.6.tgz \
    && tar -zxvf swoole.tgz \
    && rm swoole.tgz \
    && docker-php-ext-install -j$(nproc) swoole-4.8.6 \
    && curl -o /usr/bin/composer https://mirrors.aliyun.com/composer/composer.phar \
    && chmod +x /usr/bin/composer

# 第二阶段：运行阶段
FROM php:7.4-fpm

# 复制第一阶段安装的PHP扩展
COPY --from=builder /usr/local/lib/php/extensions/ /usr/local/lib/php/extensions/
COPY --from=builder /usr/local/etc/php/conf.d/ /usr/local/etc/php/conf.d/

# 复制composer
COPY --from=builder /usr/bin/composer /usr/bin/composer

# 安装运行时必需的库，不安装构建工具
# 镜像源加速
# sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
RUN apt update \
    && apt install -y --no-install-recommends \
        procps \
        libfreetype6 \
        libjpeg62-turbo \
        libpng16-16 \
        libmagickwand-6.q16-6 \
        libmagickcore-6.q16-6 \
        libzip4 \
        libpq5 \
        librabbitmq4 \
        supervisor \
    && cat /usr/share/zoneinfo/Asia/Shanghai > /etc/localtime \
    && echo 'Asia/Shanghai' > /etc/timezone \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 9000
WORKDIR /www
CMD ["/usr/local/sbin/php-fpm", "-F"]
