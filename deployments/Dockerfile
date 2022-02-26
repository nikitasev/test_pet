FROM alpine:3.12.0

RUN addgroup -S --gid 8080 webadmin && adduser -S -h /app -u 8080 -G webadmin webadmin
RUN apk add libc6-compat
ADD /bin/app /app/app
WORKDIR /app
USER webadmin

ENTRYPOINT ["/app/app"]