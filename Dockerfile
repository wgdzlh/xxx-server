FROM golang:1.21.10-alpine3.19 AS build
ENV TZ=Asia/Shanghai

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk --no-cache add tzdata gcc pkgconfig musl-dev gdal-dev \
    && ln -fs /usr/share/zoneinfo/$TZ /etc/localtime

WORKDIR /home/build
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download
COPY . .
ARG APP_VERSION=v1.0.0
RUN go build -ldflags "-s -w -X main.version=$APP_VERSION -X main.buildTime=`date +%Y-%m-%dT%H:%M:%S`"

# runtime image
FROM python:3.10.14-alpine3.19
ENV TZ=Asia/Shanghai

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk --no-cache add tzdata gdal \
    && ln -fs /usr/share/zoneinfo/$TZ /etc/localtime

WORKDIR /app
# COPY scripts/requirements.txt ./
# RUN pip install --no-cache-dir -i https://mirrors.cloud.tencent.com/pypi/simple -r requirements.txt
# COPY scripts/*.py ./
COPY --from=build /home/build/xxx-server .
ENTRYPOINT ["./xxx-server"]
