FROM alpine:latest

COPY vidispine-monitor.amd64 /usr/local/bin/vidispine-monitor
USER daemon
CMD /usr/local/bin/vidispine-monitor