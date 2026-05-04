FROM golang:1.26-alpine AS build

WORKDIR /app

COPY shared /app/shared
COPY bff /app/bff

WORKDIR /app/bff

RUN go build -o /bin/bff ./cmd/bff

FROM alpine:3.21

COPY --from=build /bin/bff /usr/local/bin/bff

EXPOSE 8080 9090

CMD ["bff"]
