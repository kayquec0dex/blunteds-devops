# ── Estágio 1: build ────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copia os arquivos de dependências primeiro
COPY go.mod go.sum ./
RUN go mod download

# Copia o código e compila
COPY . .
RUN go build -o api .

# ── Estágio 2: imagem final ─────────────────────────
FROM alpine:latest

WORKDIR /app

# Copia só o binário compilado do estágio anterior
COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]