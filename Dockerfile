ARG ARCH=
FROM ${ARCH}golang:1.21-alpine3.18 as builder

WORKDIR /app/build
RUN apk add -U make git

COPY .git .
COPY src .

# RUN ls . && exit 2

RUN make build

FROM ${ARCH}alpine:3.18 as container

WORKDIR /app
RUN apk add -U bash gettext

COPY --from=builder /app/build/idrac_exporter /app/bin/

COPY idrac.yml.template /etc/prometheus/
COPY entrypoint.sh /app

ENTRYPOINT ["/app/entrypoint.sh"]
