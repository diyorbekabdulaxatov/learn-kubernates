FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /server .

FROM alpine:3.24
COPY --from=builder /server /server

EXPOSE 8080
ENTRYPOINT [ "/server" ]