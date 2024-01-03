FROM golang:1.21

COPY development/air.toml /development/air.toml

#install mariadb-client
RUN apt-get update \
    && apt-get install -y default-mysql-client

ENV DOCKERIZE_VERSION v0.6.1
ENV AIR_VERSION v1.40.4
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz
RUN go install github.com/cosmtrek/air@$AIR_VERSION

RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

WORKDIR /server

COPY go.* /server/
RUN go mod download

ENTRYPOINT ["air", "-c", "/development/air.toml"]
