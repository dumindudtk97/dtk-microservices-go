# build the app in a docker image
# base go image
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o listnerApp .

RUN chmod +x /app/listnerApp

# build a smaller docker image by copying only the built app (executable)
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/listnerApp /app

CMD [ "/app/listnerApp" ]

