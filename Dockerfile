# build stage
FROM golang:alpine AS builder
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache build-base
WORKDIR /src
COPY . .
# ENV GOPROXY=https://goproxy.io,direct
ENV GOPROXY=https://mirrors.aliyun.com/goproxy
ENV GO111MODULE=on
RUN go build

RUN cp /src/shiori /usr/local/bin

# server image
# FROM golang:alpine
# COPY --from=builder /src/shiori /usr/local/bin/
ENV SHIORI_DIR /srv/shiori/
EXPOSE 8080
CMD ["/usr/local/bin/shiori", "serve"]
