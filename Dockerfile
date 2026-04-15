### Build Golang ###
FROM golang:1.26 AS build-backend
WORKDIR /workdir
ENV GO111MODULE="on"
ARG APP_VERSION
ARG APP_COMMIT
## download packages
COPY backend/go.mod backend/go.sum ./
RUN go mod download
## build
COPY backend/ .
RUN GOOS=linux go build -ldflags "-X main.appVersion=${APP_VERSION} -X main.appCommit=${APP_COMMIT}" -o app main.go


### Build Next.js ###
FROM node:24.14.1 AS build-frontend
RUN corepack enable
WORKDIR /workdir
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/root/.local/share/pnpm/store \
    pnpm install --frozen-lockfile
COPY frontend/ .
RUN pnpm run build


### Run ###
FROM gcr.io/distroless/base-debian12:latest
## copy binary
COPY --from=build-backend /workdir/app .
COPY --from=build-frontend /workdir/out frontend/out
## Run
ENTRYPOINT ["./app"]

