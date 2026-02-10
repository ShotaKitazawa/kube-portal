FROM golang:1.26

COPY development/air.toml /

ENV AIR_VERSION v1.40.4
RUN go install github.com/cosmtrek/air@$AIR_VERSION

RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

WORKDIR /app

COPY go.* /app/
RUN go mod download

ENTRYPOINT ["air", "-c", "/air.toml", "--"]
