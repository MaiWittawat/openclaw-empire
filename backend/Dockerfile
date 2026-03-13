FROM golang:1.25.6-alpine AS builder

RUN apk add --no-cache gcc musl-dev git

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /server ./cmd

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /server .
COPY --from=builder /app/config/config.yaml ./config/

EXPOSE 8080

CMD ["./server"]
