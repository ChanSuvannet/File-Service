# Dockerfile

FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o app .

FROM gcr.io/distroless/base

WORKDIR /app

COPY --from=builder /app/app .

COPY view /app/view

EXPOSE 8080

CMD ["/app/app"]

CMD ["/app/app"]
