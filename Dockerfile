FROM golang:1.19.9-alpine AS builder 
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download 

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

FROM alpine:3.17

ENV APP_ENVIRONMENT=PRODUCTION

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /app/cloud-native-app .

ENTRYPOINT ["./cloud-native-app"]
