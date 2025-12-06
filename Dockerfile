FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY tapline /usr/local/bin/tapline

RUN chmod +x /usr/local/bin/tapline

USER nobody:nobody

ENTRYPOINT ["tapline"]
CMD ["--help"]
