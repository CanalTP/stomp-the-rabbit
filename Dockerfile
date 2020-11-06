FROM golang:1.15 AS builder
WORKDIR /src
RUN apt update && apt install -y make
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN make STATIC=1 build


FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /src/stomptherabbit /
ENTRYPOINT ["/stomptherabbit"]
