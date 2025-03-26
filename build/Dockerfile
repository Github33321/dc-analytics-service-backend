FROM golang:1.23-alpine AS builder
WORKDIR /app


COPY .. .


RUN ls -laR .


RUN CGO_ENABLED=0 GOOS=linux go build -o dc-analytics-service-backend ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/dc-analytics-service-backend .

EXPOSE 8081

CMD ["./dc-analytics-service-backend"]
