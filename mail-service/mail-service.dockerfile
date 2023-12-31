# base Go image
FROM golang:1.21-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o mailApp ./src/api

RUN chmod +x /app/mailApp

#build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/mailApp /app
COPY --from=builder /app/src/templates /templates

CMD ["/app/mailApp"]