FROM node:22

# Atualizar e instalar dependências necessárias
RUN apt-get update && apt-get install -y \
    openssl procps wget inetutils-ping && \
    wget https://downloads.mongodb.com/compass/mongodb-mongosh_2.3.4_amd64.deb && \
    dpkg -i mongodb-mongosh_2.3.4_amd64.deb && \
    apt-get install -f -y && \
    rm -f mongodb-mongosh_2.3.4_amd64.deb && \
    mongosh --version

# Instalar CLI do Nest.js
RUN npm install -g @nestjs/cli

# Instale o Go 1.23.4
RUN wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz && \
    rm go1.23.4.linux-amd64.tar.gz

# Configure variáveis de ambiente para o Go
ENV PATH="/usr/local/go/bin:${PATH}" \
    GOPATH="/go" \
    GOROOT="/usr/local/go"

# Ajustar permissões para o diretório /go
RUN mkdir -p /go && chown node:node /go

# Configurar variável de ambiente para DevContainer
ENV DEVCONTAINER=true
