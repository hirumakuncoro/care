FROM golang:1.25.1-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/care ./cmd/care

FROM alpine:3.21

RUN apk add --no-cache ca-certificates
COPY --from=builder /app/care /app/care

ENTRYPOINT ["/app/care"]