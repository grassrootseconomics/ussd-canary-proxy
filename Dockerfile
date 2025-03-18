FROM golang:1.23-bullseye as build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

ENV BUILD_COMMIT=${BUILD_COMMIT}

WORKDIR /build

COPY go.* .
RUN go mod download

COPY . .
RUN go build -o ussd-canary-proxy -ldflags="-X main.build=${BUILD_COMMIT} -s -w" cmd/service/*


FROM debian:bullseye-slim

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /service

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/ussd-canary-proxy .
COPY migrations migrations/
COPY config.toml .
COPY queries.sql .
COPY LICENSE .

EXPOSE 5000

CMD ["./ussd-canary-proxy"]