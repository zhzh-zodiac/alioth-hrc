# 构建阶段（与 go.mod 的 go 版本对齐；若拉取失败可改为 golang:bookworm 并启用 GOTOOLCHAIN=auto）
FROM golang:1.25-bookworm AS build
ENV GOTOOLCHAIN=auto

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /server ./cmd/server

# 运行阶段（带 wget 供 healthcheck）
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata wget

COPY --from=build /server /server

ENV TZ=Asia/Shanghai

EXPOSE 8080

HEALTHCHECK --interval=15s --timeout=5s --start-period=20s --retries=5 \
  CMD wget -q -O /dev/null http://127.0.0.1:8080/healthz || exit 1

ENTRYPOINT ["/server"]
