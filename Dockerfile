FROM alpine:3.23
RUN apk add --no-cache iperf3
USER iperf3:iperf3
WORKDIR exporter
COPY iperf3_exporter .
ENTRYPOINT ["/exporter/iperf3_exporter"]
