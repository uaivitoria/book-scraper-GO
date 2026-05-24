# ============================================================
# Stage 1: builder — compila o binário Go
# ============================================================
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copia os arquivos de dependências primeiro (cache layer)
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o código fonte
COPY main.go .

# Compila o binário estático (sem dependências externas)
RUN CGO_ENABLED=0 GOOS=linux go build -o scraper main.go

# ============================================================
# Stage 2: runtime — imagem final minimalista
# ============================================================
FROM alpine:3.19

WORKDIR /app

# Copia apenas o binário compilado do stage anterior
COPY --from=builder /app/scraper .

# Cria usuário não-root por segurança
RUN adduser -D scraper
USER scraper

# Pasta onde os arquivos gerados serão salvos
VOLUME ["/app/output"]

# Comando padrão ao rodar o container
CMD ["./scraper"]