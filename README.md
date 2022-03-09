# Prometheus vCloud Director SD

[![Current Tag](https://img.shields.io/github/v/tag/promhippie/prometheus-vcd-sd?sort=semver)](https://github.com/promhippie/prometheus-vcd-sd) [![Build Status](https://github.com/promhippie/prometheus-vcd-sd/actions/workflows/general.yml/badge.svg)](https://github.com/promhippie/prometheus-vcd-sd/actions) [![Join the Matrix chat at https://matrix.to/#/#webhippie:matrix.org](https://img.shields.io/badge/matrix-%23webhippie-7bc9a4.svg)](https://matrix.to/#/#webhippie:matrix.org) [![Docker Size](https://img.shields.io/docker/image-size/promhippie/prometheus-vcd-sd/latest)](https://hub.docker.com/r/promhippie/prometheus-vcd-sd) [![Docker Pulls](https://img.shields.io/docker/pulls/promhippie/prometheus-vcd-sd)](https://hub.docker.com/r/promhippie/prometheus-vcd-sd) [![Go Reference](https://pkg.go.dev/badge/github.com/promhippie/prometheus-vcd-sd.svg)](https://pkg.go.dev/github.com/promhippie/prometheus-vcd-sd) [![Go Report Card](https://goreportcard.com/badge/github.com/promhippie/prometheus-vcd-sd)](https://goreportcard.com/report/github.com/promhippie/prometheus-vcd-sd) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/8a5f11b7eb86484eb11ff56253ac20a2)](https://www.codacy.com/gh/promhippie/prometheus-vcd-sd?utm_source=github.com&utm_medium=referral&utm_content=promhippie/prometheus-vcd-sd&utm_campaign=Badge_Grade)

This project provides a server to automatically discover nodes within your vCloud Director organization in a Prometheus SD compatible format.

## Install

You can download prebuilt binaries from our [GitHub releases](https://github.com/promhippie/prometheus-vcd-sd/releases), or you can use our Docker images published on [Docker Hub](https://hub.docker.com/r/promhippie/prometheus-vcd-sd/tags/) or [Quay](https://quay.io/repository/promhippie/prometheus-vcd-sd?tab=tags). If you need further guidance how to install this take a look at our [documentation](https://promhippie.github.io/prometheus-vcd-sd/#getting-started).

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.11.

```bash
git clone https://github.com/promhippie/prometheus-vcd-sd.git
cd prometheus-vcd-sd

make generate build

./bin/prometheus-vcd-sd -h
```

## Security

If you find a security issue please contact [thomas@webhippie.de](mailto:thomas@webhippie.de) first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

-   [Thomas Boerger](https://github.com/tboerger)

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2020 Thomas Boerger <thomas@webhippie.de>
```
