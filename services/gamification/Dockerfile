FROM golang:1.26-alpine AS build

WORKDIR /app

COPY shared /app/shared
COPY services/gamification /app/services/gamification

WORKDIR /app/services/gamification

RUN go build -o /bin/gamification ./cmd/gamification

FROM alpine:3.21

COPY --from=build /bin/gamification /usr/local/bin/gamification

EXPOSE 50052

CMD ["gamification"]
