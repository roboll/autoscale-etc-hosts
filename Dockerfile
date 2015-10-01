###############################################################################
# https://github.com/roboll/autoscale-etc-hosts
###############################################################################
FROM alpine:3.2

RUN apk --update add ca-certificates && rm -rf /var/cache/apk/*

ADD target/autoscale-etc-hosts-linux-amd64 /autoscale-etc-hosts
ENTRYPOINT ["/autoscale-etc-hosts"]
