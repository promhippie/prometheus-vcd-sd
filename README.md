# Prometheus vCloud Director SD

[![Build Status](http://cloud.drone.io/api/badges/promhippie/prometheus-vcd-sd/status.svg)](http://cloud.drone.io/promhippie/prometheus-vcd-sd)
[![Join the Matrix chat at https://matrix.to/#/#webhippie:matrix.org](https://img.shields.io/badge/matrix-%23webhippie-7bc9a4.svg)](https://matrix.to/#/#webhippie:matrix.org)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/8a5f11b7eb86484eb11ff56253ac20a2)](https://www.codacy.com/gh/promhippie/prometheus-vcd-sd?utm_source=github.com&utm_medium=referral&utm_content=promhippie/prometheus-vcd-sd&utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/promhippie/prometheus-vcd-sd?status.svg)](http://godoc.org/github.com/promhippie/prometheus-vcd-sd)
[![Go Report](http://goreportcard.com/badge/github.com/promhippie/prometheus-vcd-sd)](http://goreportcard.com/report/github.com/promhippie/prometheus-vcd-sd)
[![](https://images.microbadger.com/badges/image/promhippie/prometheus-vcd-sd.svg)](http://microbadger.com/images/promhippie/prometheus-vcd-sd "Get your own image badge on microbadger.com")

This project provides a server to automatically discover nodes within your vCloud Director organization in a Prometheus SD compatible format.

## Install

You can download prebuilt binaries from our [GitHub releases](https://github.com/promhippie/prometheus-vcd-sd/releases), or you can use our Docker images published on [Docker Hub](https://hub.docker.com/r/promhippie/prometheus-vcd-sd/tags/). If you need further guidance how to install this take a look at our [documentation](https://promhippie.github.io/prometheus-vcd-sd/#getting-started).

## Integration

### Available labels

The following list of meta labels can be used to relabel your scrape results entirely. Hopefully the names are self-explaining, that's why I have skipped a description for each label.

-   `__address__`
-   `__meta_vcd_metadata_<name>`
-   `__meta_vcd_name`
-   `__meta_vcd_id`
-   `__meta_vcd_network_<name>`
-   `__meta_vcd_num_cores_per_socket`
-   `__meta_vcd_num_cpus`
-   `__meta_vcd_org`
-   `__meta_vcd_os_type`
-   `__meta_vcd_project`
-   `__meta_vcd_status`
-   `__meta_vcd_storage_profile`
-   `__meta_vcd_vdc`

### Prometheus config

Here you get a snippet for the Prometheus `scrape_config` that configures Prometheus to scrape `node_exporter` assuming that it is deployed on all your servers.

```yml
- job_name: node
  file_sd_configs:
  - files: [ "/etc/prometheus/vcd.json" ]
  relabel_configs:
  - source_labels: [__meta_vcd_network_internal]
    replacement: "${1}:9100"
    target_label: __address__
  - source_labels: [__meta_vcd_org]
    target_label: org
  - source_labels: [__meta_vcd_name]
    target_label: instance
```

## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.11.

```console
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
