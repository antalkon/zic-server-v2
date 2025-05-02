# Стадия сборки
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

RUN apk update && apk add --no-cache \
    ca-certificates git gcc g++ libc-dev binutils

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify && go mod tidy

COPY . .

RUN go build -o bin/application ./cmd/backend

# Стадия рантайма
FROM --platform=$TARGETPLATFORM alpine:3.21 AS runner

RUN apk add --no-cache \
    ca-certificates libc6-compat bash \
    && rm -rf /var/cache/apk/*

WORKDIR /app

# бинарник
COPY --from=builder /app/bin/application ./

# конфиги
COPY --from=builder /app/config/yaml ./config/yaml

ENV APP_ENV=prod

CMD ["./application"]