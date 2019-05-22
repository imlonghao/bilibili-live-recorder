FROM golang:1.12-alpine AS builder
LABEL maintainer="imlonghao <dockerfile@esd.cc>"
WORKDIR /builder
RUN apk add upx && \
    GO111MODULE=on go build -ldflags="-s -w" -o /bilibili-live-recorder
    upx --ultra-brute /bilibili-live-recorder

FROM alpine:latest
LABEL maintainer="imlonghao <dockerfile@esd.cc>"
RUN apk --no-cache add ca-certificates
COPY --from=builder /bilibili-live-recorder .
CMD ["/bilibili-live-recorder"]