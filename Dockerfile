### Build Golang ###
FROM golang:1.21 AS build-backend
WORKDIR /workdir
ENV GO111MODULE="on"
ARG APP_VERSION
ARG APP_COMMIT
## download packages
COPY go.mod go.sum ./
RUN go mod download
## build
COPY cmd cmd
COPY server server
RUN GOOS=linux go build -ldflags "-X main.appVersion=${APP_VERSION} -X main.appCommit=${APP_COMMIT}" -o app cmd/kubeportal/main.go


### Build Next.js ###
FROM node:20.10.0 AS build-frontend
WORKDIR /workdir
COPY client/ ./
RUN yarn install
RUN yarn build


### Run ###
FROM gcr.io/distroless/base
## copy binary
COPY --from=build-backend /workdir/app .
COPY --from=build-frontend /workdir/out client/out
## Run
ENTRYPOINT ["./app"]

