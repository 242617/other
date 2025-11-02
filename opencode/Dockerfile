FROM ubuntu:22.04

# Настройка локали для кириллицы
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

# Настройка временной зоны
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
	echo $TZ > /etc/timezone

# Установка зависимостей
RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		build-essential \
		curl \
		dnsutils \
		git \
		iproute2\
		iputils-ping \
		jq \
		nano \
		netcat-openbsd \
		openssl \
		python3 \
		python3-pip \
		rsync \
		htop \
		telnet \
		traceroute \
		tree \
		tzdata \
		unzip \
		wget
RUN apt-get clean
RUN rm -rf /var/lib/apt/lists/*

# Установка Node.js 20
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash && \
	apt install -y nodejs

# Установка OpenCode
RUN curl -fsSL https://opencode.ai/install | VERSION=0.15.19 bash

# Конфигурирование
COPY ./build/ /root/.config/opencode

# Рабочая директория
WORKDIR /workspace

ENTRYPOINT ["/root/.opencode/bin/opencode"]
