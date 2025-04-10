FROM golang:1.23-alpine AS builder
WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go mod tidy


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dc-analytics-service ./cmd


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app


COPY --from=builder /app/dc-analytics-service .

COPY .env .

EXPOSE 7002

ENTRYPOINT ["/app/dc-analytics-service"]
