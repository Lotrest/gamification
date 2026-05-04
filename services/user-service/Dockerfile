FROM golang:1.26-alpine AS build

WORKDIR /app

COPY shared /app/shared
COPY services/user-service /app/services/user-service

WORKDIR /app/services/user-service

RUN go build -o /bin/user-service ./cmd/user-service

FROM alpine:3.21

COPY --from=build /bin/user-service /usr/local/bin/user-service

EXPOSE 50051

CMD ["user-service"]
