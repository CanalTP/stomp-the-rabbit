FROM golang:1.15 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags="-w -s" ./cmd/stomptherabbit


FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /src/stomptherabbit /
ENTRYPOINT ["/stomptherabbit"]
