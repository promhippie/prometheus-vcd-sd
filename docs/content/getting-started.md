---
title: "Getting Started"
date: 2018-05-02T00:00:00+00:00
anchor: "getting-started"
weight: 10
---

## Installation

We won't cover further details how to properly setup [Prometheus](https://prometheus.io) itself, we will only cover some basic setup based on [docker-compose](https://docs.docker.com/compose/). But if you want to run this service discovery without [docker-compose](https://docs.docker.com/compose/) you should be able to adopt that to your needs.

First of all we need to prepare a configuration for [Prometheus](https://prometheus.io) that includes the service discovery which simply maps to a node exporter.

{{< highlight yaml >}}
global:
  scrape_interval: 1m
  scrape_timeout: 10s
  evaluation_interval: 1m

scrape_configs:
- job_name: node
  file_sd_configs:
  - files: [ "/etc/sd/vcd.json" ]
  relabel_configs:
  - source_labels: [__meta_VCD_network_internal]
    replacement: "${1}:9100"
    target_label: __address__
  - source_labels: [__meta_VCD_name]
    target_label: instance
- job_name: vcd-sd
  static_configs:
  - targets:
    - vcd-sd:9000
{{< / highlight >}}

After preparing the configuration we need to create the `docker-compose.yml` within the same folder, this `docker-compose.yml` starts a simple [Prometheus](https://prometheus.io) instance together with the service discovery. Don't forget to update the envrionment variables with the required credentials. If you are using a different volume for the service discovery you have to make sure that the container user is allowed to write to this volume.

{{< highlight yaml >}}
version: '2'

volumes:
  prometheus:

services:
  prometheus:
    image: prom/prometheus:v2.6.0
    restart: always
    ports:
      - 9090:9090
    volumes:
      - prometheus:/prometheus
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - ./service-discovery:/etc/sd

  vcd-exporter:
    image: promhippie/prometheus-vcd-sd:latest
    restart: always
    environment:
      - PROMETHEUS_VCD_LOG_PRETTY=true
      - PROMETHEUS_VCD_OUTPUT_FILE=/etc/sd/vcd.json
      - PROMETHEUS_VCD_URL=https://vdc.example.com/api
      - PROMETHEUS_VCD_USERNAME=username
      - PROMETHEUS_VCD_PASSWORD=p455w0rd
      - PROMETHEUS_VCD_ORG=MY-ORG1
      - PROMETHEUS_VCD_VDC=MY-ORG1-DC1
    volumes:
      - ./service-discovery:/etc/sd
{{< / highlight >}}

Since our `latest` Docker tag always refers to the `master` branch of the Git repository you should always use some fixed version. You can see all available tags at our [DockerHub repository](https://hub.docker.com/r/promhippie/prometheus-vcd-sd/tags/), there you will see that we also provide a manifest, you can easily start the exporter on various architectures without any change to the image name. You should apply a change like this to the `docker-compose.yml`:

{{< highlight diff >}}
  vcd-exporter:
-   image: promhippie/prometheus-vcd-sd:latest
+   image: promhippie/prometheus-vcd-sd:0.1.0
    restart: always
    environment:
      - PROMETHEUS_VCD_LOG_PRETTY=true
      - PROMETHEUS_VCD_OUTPUT_FILE=/etc/sd/vcd.json
      - PROMETHEUS_VCD_URL=https://vdc.example.com/api
      - PROMETHEUS_VCD_USERNAME=username
      - PROMETHEUS_VCD_PASSWORD=p455w0rd
      - PROMETHEUS_VCD_ORG=MY-ORG1
      - PROMETHEUS_VCD_VDC=MY-ORG1-DC1
    volumes:
      - ./service-discovery:/etc/sd
{{< / highlight >}}

Depending on how you have launched and configured [Prometheus](https://prometheus.io) it's possible that it's running as user `nobody`, in that case you should run the service discovery as this user as well, otherwise [Prometheus](https://prometheus.io) won't be able to read the generated JSON file:

{{< highlight diff >}}
  vcd-exporter:
    image: promhippie/prometheus-vcd-sd:latest
    restart: always
+   user: '65534'
    environment:
      - PROMETHEUS_VCD_LOG_PRETTY=true
      - PROMETHEUS_VCD_OUTPUT_FILE=/etc/sd/vcd.json
      - PROMETHEUS_VCD_URL=https://vdc.example.com/api
      - PROMETHEUS_VCD_USERNAME=username
      - PROMETHEUS_VCD_PASSWORD=p455w0rd
      - PROMETHEUS_VCD_ORG=MY-ORG1
      - PROMETHEUS_VCD_VDC=MY-ORG1-DC1
    volumes:
      - ./service-discovery:/etc/sd
{{< / highlight >}}

Finally the service discovery should be configured fine, let's start this stack with [docker-compose](https://docs.docker.com/compose/), you just need to execute `docker-compose up` within the directory where you have stored `prometheus.yml` and `docker-compose.yml`.

{{< highlight txt >}}
Creating network "vcd-sd_default" with the default driver
Creating volume "vcd-sd_prometheus" with default driver
Creating vcd-sd_vcd-exporter_1 ... done
Creating vcd-sd_prometheus_1   ... done
Attaching to vcd-sd_vcd-exporter_1, vcd-sd_prometheus_1
prometheus_1    | level=info ts=2020-06-20T21:41:32.6738077Z caller=main.go:243 msg="Starting Prometheus" version="(version=2.6.0, branch=HEAD, revision=dbd1d58c894775c0788470944b818cc724f550fb)"
prometheus_1    | level=info ts=2020-06-20T21:41:32.6739242Z caller=main.go:244 build_context="(go=go1.11.3, user=root@bf5760470f13, date=20181217-15:14:46)"
prometheus_1    | level=info ts=2020-06-20T21:41:32.6739553Z caller=main.go:245 host_details="(Linux 4.19.76-linuxkit #1 SMP Tue May 26 11:42:35 UTC 2020 x86_64 c40f124a0e86 (none))"
prometheus_1    | level=info ts=2020-06-20T21:41:32.6743456Z caller=main.go:246 fd_limits="(soft=1048576, hard=1048576)"
prometheus_1    | level=info ts=2020-06-20T21:41:32.6744037Z caller=main.go:247 vm_limits="(soft=unlimited, hard=unlimited)"
prometheus_1    | level=info ts=2020-06-20T21:41:32.6754939Z caller=main.go:561 msg="Starting TSDB ..."
prometheus_1    | level=info ts=2020-06-20T21:41:32.6759495Z caller=web.go:429 component=web msg="Start listening for connections" address=0.0.0.0:9090
prometheus_1    | level=info ts=2020-06-20T21:41:32.6832124Z caller=main.go:571 msg="TSDB started"
prometheus_1    | level=info ts=2020-06-20T21:41:32.6832976Z caller=main.go:631 msg="Loading configuration file" filename=prometheus.yml
prometheus_1    | level=info ts=2020-06-20T21:41:32.6862854Z caller=main.go:657 msg="Completed loading of configuration file" filename=prometheus.yml
prometheus_1    | level=info ts=2020-06-20T21:41:32.6863411Z caller=main.go:530 msg="Server is ready to receive web requests."
vcd-exporter_1  | level=info ts=2020-06-20T21:41:33.1219629Z msg="Launching Prometheus vCloud Director SD" version=1e2c111 revision=1e2c111 date=20200620 go=go1.14.2
vcd-exporter_1  | level=info ts=2020-06-20T21:41:33.4333853Z msg="Starting metrics server" addr=0.0.0.0:9000
{{< / highlight >}}

That's all, the service discovery should be up and running. You can access [Prometheus](https://prometheus.io) at [http://localhost:9090](http://localhost:9090).

{{< figure src="service-discovery.png" title="Prometheus service discovery for vCloud Director" >}}

## Kubernetes

Currently we have not prepared a deployment for Kubernetes, but this is something we will provide for sure. Most interesting will be the integration into the [Prometheus Operator](https://coreos.com/operators/prometheus/docs/latest/), so stay tuned.

## Configuration

### Envrionment variables

If you prefer to configure the service with environment variables you can see the available variables below, in case you want to configure multiple accounts with a single service you are forced to use the configuration file as the environment variables are limited to a single account. As the service is pretty lightweight you can even start an instance per account and configure it entirely by the variables, it's up to you.

PROMETHEUS_VCD_CONFIG
: Path to vCloud Director configuration file, optionally, required for multi credentials

PROMETHEUS_VCD_URL
: URL for the vCloud Director API, required for authentication

PROMETHEUS_VCD_INSECURE
: Insecure access for the vCloud Director API, required for authentication

PROMETHEUS_VCD_USERNAME
: Username for the vCloud Director API, required for authentication

PROMETHEUS_VCD_PASSWORD
: Password for the vCloud Director API, required for authentication

PROMETHEUS_VCD_ORG
: Organization for the vCloud Director API, required for authentication

PROMETHEUS_VCD_VDC
: vDatacenter for the vCloud Director API, required for authentication

PROMETHEUS_VCD_LOG_LEVEL
: Only log messages with given severity, defaults to `info`

PROMETHEUS_VCD_LOG_PRETTY
: Enable pretty messages for logging, defaults to `true`

PROMETHEUS_VCD_WEB_ADDRESS
: Address to bind the metrics server, defaults to `0.0.0.0:9000`

PROMETHEUS_VCD_WEB_PATH
: Path to bind the metrics server, defaults to `/metrics`

PROMETHEUS_VCD_OUTPUT_FILE
: Path to write the file_sd config, defaults to `/etc/prometheus/vcd.json`

PROMETHEUS_VCD_OUTPUT_REFRESH
: Discovery refresh interval in seconds, defaults to `30`

### Configuration file

Especially if you want to configure multiple accounts within a single service discovery you got to use the configuration file. So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/promhippie/prometheus-vcd-sd/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options, they also include the default values.

## Labels

* `__address__`
* `__meta_vcd_metadata_<name>`
* `__meta_vcd_name`
* `__meta_vcd_network_<name>`
* `__meta_vcd_num_cores_per_socket`
* `__meta_vcd_num_cpus`
* `__meta_vcd_org`
* `__meta_vcd_os_type`
* `__meta_vcd_project`
* `__meta_vcd_status`
* `__meta_vcd_storage_profile`
* `__meta_vcd_vdc`

## Metrics

prometheus_vcd_sd_request_duration_seconds
: Histogram of latencies for requests to the vCloud Director API

prometheus_vcd_sd_request_failures_total
: Total number of failed requests to the vCloud Director API
