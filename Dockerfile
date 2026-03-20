FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go


FROM alpine:3.14

WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD [ "./main" ]