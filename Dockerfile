FROM golang:1.19-alpine AS builder

# Instalação de dependências necessárias
RUN apk add --no-cache git

# Configuração do diretório de trabalho
WORKDIR /app

# Copia o arquivo go.mod primeiro para aproveitar o cache
COPY go.mod ./
RUN go mod download

# Copia o restante dos arquivos do projeto
COPY . .

# Compila o projeto
RUN go build -o hardshell ./cmd/hardshell

# Imagem final
FROM alpine:3.16

# Instalação de ferramentas necessárias
RUN apk add --no-cache bash

# Copia o binário compilado
COPY --from=builder /app/hardshell /usr/local/bin/hardshell

# Copia os arquivos de configuração
COPY --from=builder /app/configs /etc/hardshell/configs

# Define o diretório de trabalho
WORKDIR /data

# Define o ponto de entrada
ENTRYPOINT ["hardshell"]

# Comando padrão (pode ser sobrescrito)
CMD ["--help"]
