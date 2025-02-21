FROM golang:1.24.0-alpine3.21
ARG PROTOBUF_VERSION=29.3
ARG MOCKERY_VERSION=2.52.1
ARG MIGRATE=4.18.2
EXPOSE 8080

RUN apk update \
&& apk add alpine-sdk \
&& apk --update add tzdata \
&& cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
&& apk del tzdata \
&& rm -rf /var/cache/apk/*

WORKDIR /tmp/protoc
RUN apk add --no-cache unzip curl \
  && curl -L https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip -o protoc.zip \
  && unzip protoc.zip \
  && mv /tmp/protoc/bin/* /usr/local/bin/ \
  && mv /tmp/protoc/include /usr/local/include \
  && rm -rf protoc protoc.zip \
  && apk add --no-cache protobuf-dev

RUN curl -L https://github.com/vektra/mockery/releases/download/v${MOCKERY_VERSION}/mockery_${MOCKERY_VERSION}_Linux_x86_64.tar.gz | tar xvz && mv ./mockery /usr/bin/mockery \
  && curl -L https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE}/migrate.linux-amd64.tar.gz | tar xvz && mv ./migrate /usr/bin/migrate

WORKDIR /opt/task_manage
COPY ./backend /opt/task_manage
RUN go mod tidy
RUN go install github.com/air-verse/air@latest
RUN go install github.com/bufbuild/buf/cmd/buf@latest \
  && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
  && go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest