FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bot .

FROM alpine:latest AS runner
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/bot .
ENTRYPOINT ["./bot"]