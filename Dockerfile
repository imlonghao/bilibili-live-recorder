FROM golang:1.13-alpine AS builder
LABEL maintainer="imlonghao <dockerfile@esd.cc>"
WORKDIR /builder
COPY . /builder
RUN apk add upx && \
    GO111MODULE=on go build -mod=vendor -ldflags="-s -w" -o /bilibili-live-recorder && \
    upx /bilibili-live-recorder

FROM alpine:latest
LABEL maintainer="imlonghao <dockerfile@esd.cc>"
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /bilibili-live-recorder .
CMD ["/bilibili-live-recorder"]
EXPOSE 8080