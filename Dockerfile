# build stage
FROM golang:1.21rc3-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["/app/main"]