FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/bot

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Europe/Moscow
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["/root/main"]