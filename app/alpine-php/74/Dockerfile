# ==========================================
# 第一阶段：编译阶段 (builder)
# ==========================================
FROM php:7.4.33-fpm-alpine3.16 AS builder

# 安装所有编译所需的工具和 -dev 开发库
RUN apk add --no-cache \
    autoconf g++ make libtool \
    libpng-dev libjpeg-turbo-dev freetype-dev \
    libzip-dev bzip2-dev gettext-dev \
    postgresql-dev rabbitmq-c-dev imagemagick-dev

# 编译并安装扩展
# 1. 先通过 pecl 安装，但不直接用 enable，因为 redis 容易报错
RUN pecl install amqp-1.11.0 \
    && pecl install igbinary \
    && pecl install imagick-3.7.0 \
    && pecl install redis-5.3.7

# 2. 统一启用 PECL 扩展 (注意这里明确指定 .so)
RUN docker-php-ext-enable amqp igbinary imagick redis

# 3. 安装内置扩展 (保持不变)
RUN docker-php-ext-configure gd --with-freetype --with-jpeg \
    && docker-php-ext-install -j$(nproc) \
        bcmath bz2 exif gd gettext pcntl pdo_mysql sockets opcache zip pgsql pdo_pgsql

# ==========================================
# 第二阶段：最终运行阶段 (final)
# ==========================================
FROM php:7.4.33-fpm-alpine3.16

# 1. 只安装运行时必须的动态库（非 -dev 库）
# 这些库通常比 -dev 库小得多
RUN apk add --no-cache \
    bash \
    supervisor \
    libpng libjpeg-turbo freetype \
    libzip bzip2 gettext \
    libpq \
    rabbitmq-c \
    imagemagick \
    tzdata
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && printf "[Date]\ndate.timezone = Asia/Shanghai\n" > /usr/local/etc/php/conf.d/timezone.ini

# (可选) 设置环境变量，某些程序会读取此变量
ENV TZ=Asia/Shanghai

# 2. 从 builder 阶段直接复制编译好的扩展文件 (.so) 和配置文件
COPY --from=builder /usr/local/lib/php/extensions /usr/local/lib/php/extensions
COPY --from=builder /usr/local/etc/php/conf.d /usr/local/etc/php/conf.d

# 3. 复制 Composer (官方多阶段构建示例)
COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

# 4. 配置 Supervisor
RUN mkdir -p /var/log/supervisor
COPY supervisord.conf /etc/supervisord.conf

WORKDIR /www

EXPOSE 9000

CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisord.conf"]
