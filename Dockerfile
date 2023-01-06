FROM golang:1.19 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go build -o outboundForwarder .

FROM scratch
COPY --from=builder /build/outboundForwarder /
COPY --from=builder /build/outboundForwarder.log /

ENTRYPOINT ["/outboundForwarder"]