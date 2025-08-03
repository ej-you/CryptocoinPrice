# -----
# BUILD
# -----

FROM golang:1.24.4-alpine AS build

# set up workdir
WORKDIR /go/src

# install dependences
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

# compile migrator
COPY ./cmd/migrator ./cmd/migrator
COPY ./config ./config
COPY ./internal/pkg/migrate ./internal/pkg/migrate
RUN go build -o ./migrator ./cmd/migrator/main.go

# compile app
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN go build -o ./app ./cmd/app/main.go

# ---
# RUN
# ---

FROM alpine:3.21 AS run

WORKDIR /app

# copy compiled app and migrator files
COPY --from=build /go/src/app .
COPY --from=build /go/src/migrator .
# copy migrations and files for swagger
COPY ./migrations ./migrations
COPY ./docs ./docs

# run app
CMD ["/bin/sh", "-c", "/app/app"]
