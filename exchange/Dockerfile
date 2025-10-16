FROM golang:1.25-alpine AS builder
WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod downlaod
COPY go.mod ./
COPY . .
EXPOSE 8080
RUN go build -o app ./main.go

FROM alpine:latest
COPY --from=builder /app/app /app
ENTRYPOINT [ "/app" ]
