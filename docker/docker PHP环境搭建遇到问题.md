# docker PHP环境搭建遇到问题

mysql8.0 版本的容器运行后如果报这种错误：

Authentication plugin 'caching_sha2_password' cannot be loaded: 找不到指定的模块

解决方法：

```
docker exec -it CONTAINER_ID bash
```

`mysql -h localhost -u root -p`

```
ALTER USER 'username' IDENTIFIED WITH mysql_native_password BY 'password';
```

####  docker内安装PHP扩展

`docker-php-ext-install 扩展名`

在PHP代码中访问数据应该这样连接

`$pdo = new PDO('mysql:host=mysql;dbname=xxx','root','123456');`

这是因为PHP代码是在FPM容器中，FPM容器启动时会自动在`/etc/hosts`中加上：

`172.17.0.2 mysql 11e55f91c4c3 dlnmp_mysql_1`（这里的ip 和容器id 容器名称不固定）

就是说，MySQL自动指向了MySQL容器动态生成的IP

Redis使用和MySQL类似，将上面的hosts参数用redis，端口用6379

#### 让DNMP随系统启动

在Ubuntu中的Launcher中搜索Startup Applications（启动应用程序）

然后Add一项，名字：Dnmp，命令填：

`docker compose -f /home/dnmp/docker-compose.yml up -d`（-f 指定docker-compose.yml文件位置）

其他Linux系统可以在`/etc/rc.local`文件，加上上面的命令

docker 访问composer

`docker exec -it phpcontainerID /bin/bash`进入到容器里面使用composer

没有读写权限的问题进入到php容器内 将`/var/www/html`的用户和用户组改变，`/var/www/html`这是我挂载的网站目录不固定：

`docker exec -it phpcontainerID /bin/bash`

`chmod -R www-data:www-data /var/www/html`

#### docker安装PHP扩展

PHP源码安装,官方提供了docker-php-source快捷脚本，用于对源文件压缩包的解压（extract）及解压后文件进行删除（delete）的操作：

```dockerfile
FROM php:latest
RUN docker-php-source extract \
	# 此处开始执行需要的操作 \
	&& docker-php-source delete
```

**注意：一定要记得删除，否则解压出来的文件会大大增加镜像的文件大小**

**核心扩展**

官方提供了docker-php-ext-configure和docker-php-ext-install快捷脚本：

```dockerfile
FROM php:latest
RUN apt-get update \
	# 相关依赖必须手动安装
	&& apt-get install -y \
	libfreetype6-dev \
	libjpeg62-turbo-dev \
	libmcrypt-dev \
	libpng-dev \
	# 安装扩展
	&& docker-php-ext-install -j$(nproc) iconv mcrypt \
	# 如果安装的扩展需要自定义配置时
	&& docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/ \
	&& docker-php-ext-install -j$(nproc) gd
```

**注意：这里的docker-php-ext-configure和docker-php-ext-install已经包含了docker-php-source的操作，所以不需要再手动去执行**

**PECL扩展**

有一些扩展并不包含在PHP源码文件中，所以使用pecl install 安装扩展，然后使用官方提供的docker-php-ext-enable快捷脚本启用扩展，如下：

```dockerfile
FROM php:latest
RUN apt-get update \
	# 手动安装依赖
	&& apt-getinstall -y libmemcached-dev zlib1g-dev \
	# 安装需要的扩展
	&& pecl install redis-4.0.2 \
	# 启用扩展
	&& docker-php-ext-enable redis
```

**其他扩展**

一些既不在PHP源码包，也不在pecl扩展仓库中的扩展，可以通过下载扩展程序源码，编译安装的方式安装：

```dockerfile
FROM php:latest
RUN curl -fsSL 'https://xcache.lighttpd.net/pub/Releases/3.2.0/xcache-3.2.0.tar.gz' -o 		xcache.tar.gz \
	&& mkdir -p xcache \
	&& tar -xf xcache.tar.gz -C xcache --strip-components=1 \
	&& rm xcache.tar.gz \
	&& ( \
		cd xcache \
		&& phpize \ 
		&& ./configure --enable-xcache \
		&& make -j$(nproc) \
		&& make install \
	) \
	&& rm -r xcache \
	&& docker-php-ext-enable xcache
```

**注意：官方提的docker-php-ext-*脚本接受任意的绝对路径（不支持相对路径，以便与系统内置的扩展程序进行区分）**：

```dockerfile
FROM php:latest
RUN curl -fsSL 'https://xcache.lighttpd.net/pub/Releases/3.2.0/xcache-3.2.0.tar.gz' -o xcache.tar.gz \
	&& mkdir -p /tmp/xcache \
	&& tar -xf xcache.tar.gz -C /tmp/xcache --strip-components=1 \
	&& rm xcache.tar.gz \
	&& docker-php-ext-configure /tmp/xcache --enable-xcache \
	&& docker-php-ext-install /tmp/xcache \
	&& rm -r /tmp/xcache
```

docker`exited(0)`在`docker-compose.yml`文件里添加`tty:true`

#### 修改 docker-compose 文件后

```
docker stop 容器名
docker rm 容器名
docker-compose up -d --no-deps --build 镜像名 // 后台启动，不启动关联容器，构建镜像
```



