ARG TARGETPLATFORM
FROM golang:1.25.6-alpine AS builder

# ติดตั้ง build-base แทนการแยกแพ็กเกจ เพื่อความครบถ้วนของ CGO
RUN apk add --no-cache build-base git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ปรับ CGO_ENABLED=0 ถ้าไม่จำเป็นต้องใช้ C library (ทำให้รันบน Alpine ได้เสถียรสุด)
# และระบุ path ไปยังไฟล์ main.go ให้ชัดเจน
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/main.go

# ใช้ Alpine เวอร์ชั่นล่าสุด
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary มาจาก stage builder
COPY --from=builder /server .

EXPOSE 8080

# รัน binary
CMD ["./server"]