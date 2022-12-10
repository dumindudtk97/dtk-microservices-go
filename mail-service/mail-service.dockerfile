# build a smaller docker image by copying only the built app (executable)
FROM alpine:latest

RUN mkdir /app

COPY mailApp /app
COPY templates /templates

CMD [ "/app/mailApp" ]