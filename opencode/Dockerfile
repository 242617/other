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
		gcc \
		git \
		htop \
		iproute2\
		iputils-ping \
		jq \
		make \
		nano \
		netcat-openbsd \
		openssl \
		python3 \
		python3-pip \
		rsync \
		telnet \
		traceroute \
		tree \
		tzdata \
		unzip \
		wget
RUN apt-get clean
RUN rm -rf /var/lib/apt/lists/*

# Установка Go 1.25.5
RUN wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz && \
	tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz && \
	rm go1.25.5.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin/
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

ENV CGO_ENABLED=0
RUN go install github.com/birdayz/kaf/cmd/kaf@v0.2.13
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.9.3
RUN go install github.com/tsenart/vegeta@latest
RUN go install github.com/mikefarah/yq/v4@v4.45.4

# Установка uv
RUN curl -LsSf https://astral.sh/uv/install.sh | sh

# Установка Node.js 20
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash && \
	apt install -y nodejs

# Установка OpenCode
RUN curl -fsSL https://opencode.ai/install | VERSION=1.1.35 bash

# Конфигурирование
COPY ./build/ /root/.config/opencode

# Рабочая директория
WORKDIR /workspace

ENTRYPOINT ["/root/.opencode/bin/opencode"]
