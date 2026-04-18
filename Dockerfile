FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o blog .

# ---

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/blog .
COPY templates/ templates/
COPY static/ static/
COPY posts/ posts/
COPY images/ images/

EXPOSE 8080

CMD ["./blog"]
