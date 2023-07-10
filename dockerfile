# syntax=docker/dockerfile:1

FROM onepill/golang-python:latest AS build

##build
WORKDIR /usr/src

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN go build -o /start-server

##deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=build /start-server ./start-server
COPY /src/ ./src
COPY main.go ./main.go

EXPOSE 80

USER nonroot:nonroot

ENTRYPOINT ["/app/start-server"]