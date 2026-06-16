FROM --platform=$BUILDPLATFORM golang:1.26.4-alpine@sha256:f1ddd9fe14fffc091dd98cb4bfa999f32c5fc77d2f2305ea9f0e2595c5437c14 AS builder

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

FROM alpine:3.24@sha256:f5064d3e5f88c467c714509f491853ab2d951932c5cad699c0cb969dcec6f3b4

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
