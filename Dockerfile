FROM golang:1.19-alpine3.16 AS build

WORKDIR /build

COPY ./ ./
RUN go mod download
RUN go build ./cmd/main

FROM alpine:3.16

WORKDIR /app

COPY --from=build /build/main server

CMD adduser daemon

USER daemon

ENV API_PORT 8080

EXPOSE ${API_PORT}

ENTRYPOINT [ "/app/server" ]