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

# Configurar variável de ambiente para DevContainer
ENV DEVCONTAINER=true
