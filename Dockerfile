FROM --platform=$BUILDPLATFORM golang:1.24.5-alpine3.21@sha256:72ff633a5298088a576d505c51630257cf1f681fc64cecddfb5234837eb4a747 AS builder

RUN apk add --no-cache -U git curl
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

WORKDIR /go/src/prometheus-vcd-sd
COPY . /go/src/prometheus-vcd-sd/

RUN --mount=type=cache,target=/go/pkg \
    go mod download -x

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    task generate build GOOS=${TARGETOS} GOARCH=${TARGETARCH}

FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1

RUN apk add --no-cache ca-certificates mailcap && \
    addgroup -g 1337 prometheus-vcd-sd && \
    adduser -D -u 1337 -h /var/lib/prometheus-vcd-sd -G prometheus-vcd-sd prometheus-vcd-sd

EXPOSE 9000
VOLUME ["/var/lib/prometheus-vcd-sd"]
ENTRYPOINT ["/usr/bin/prometheus-vcd-sd"]
CMD ["server"]
HEALTHCHECK CMD ["/usr/bin/prometheus-vcd-sd", "health"]

ENV PROMETHEUS_VCD_OUTPUT_ENGINE="http"
ENV PROMETHEUS_VCD_OUTPUT_FILE="/var/lib/prometheus-vcd-sd/output.json"

COPY --from=builder /go/src/prometheus-vcd-sd/bin/prometheus-vcd-sd /usr/bin/prometheus-vcd-sd
WORKDIR /var/lib/prometheus-vcd-sd
USER prometheus-vcd-sd
