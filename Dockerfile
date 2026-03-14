# Estágio 1: Builder
FROM golang:1.21-alpine AS builder

# 1. Instala as ferramentas necessárias, com versão fixada para garantir consistência
RUN apk add --no-cache git
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12

WORKDIR /app

# 2. Copia todo o código-fonte
COPY . .

# 3. Remove a pasta 'docs' antiga para evitar conflitos de versão
RUN rm -rf docs

# 4. Gera uma nova pasta 'docs' usando a ferramenta com versão fixada
RUN /go/bin/swag init -g cmd/api/main.go

# 5. Adiciona a diretiva 'replace' para que o Go encontre o pacote 'docs' localmente
RUN go mod edit -replace=github.com/fiap-x/video-processor=./

# 6. Sincroniza as dependências (go.sum) com base no código completo.
RUN go mod tidy

# 7. Compila os binários.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /worker ./cmd/worker/main.go

# --- Estágio 2: Imagem Final ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /api .
COPY --from=builder /worker .

EXPOSE 8080

CMD ["./api"]
