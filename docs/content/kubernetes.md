---
title: "Kubernetes"
date: 2022-03-09T00:00:00+00:00
anchor: "kubernetes"
weight: 20
---

## Kubernetes

So far we got the deployment via [Kustomize](https://github.com/kubernetes-sigs/kustomize) to get this service discovery working on Kubernetes. We are already working on a [Helm]() chart to offer more options, dependening on your preferences.

### Kustomize

We won't cover the installation of [Kustomize](https://github.com/kubernetes-sigs/kustomize) or encryption tooling like [KSOPS](https://github.com/viaduct-ai/kustomize-sops) within this guide, to get it installed and working please consult the documentation of these projects. After the installation of [Kustomize](https://github.com/kubernetes-sigs/kustomize) you just need to prepare a `kustomization.yml` wherever you like:

{{< highlight yaml >}}
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: prometheus-vcd-sd

resources:
  - github.com/promhippie/prometheus-vcd-sd?ref=master

configMapGenerator:
  - name: prometheus-vcd-sd
    behavior: merge
    literals:
      - PROMETHEUS_VCD_LOG_LEVEL=info

secretGenerator:
  - name: prometheus-vcd-sd
    behavior: merge
    literals:
      - PROMETHEUS_VCD_URL=https://vdc.example.com/api
      - PROMETHEUS_VCD_USERNAME=username
      - PROMETHEUS_VCD_PASSWORD=p455w0rd
      - PROMETHEUS_VCD_ORG=MY-ORG1
      - PROMETHEUS_VCD_VDC=MY-ORG1-DC1
{{< / highlight >}}

After that you can simply execute `kustomize build | kubectl apply -f -` to get the manifest applied. Generally it's best to use fixed versions of Docker images, this can be done quite easy, you just need to append this block to your `kustomization.yml` to use this specific version:

{{< highlight yaml >}}
images:
  - name: quay.io/promhippie/prometheus-vcd-sd
    newTag: 1.0.0
{{< / highlight >}}

After applying this manifest the metrics of the service discovery should be directly visible within your Prometheus instance if you are using the Prometheus Operator as these manifests are providing a ServiceMonitor. To consume the service discovery you got to apply a custom scrape configuration. For instructions how to apply additional scrape configs to Prometheus Operator please take a look at the [documentation](https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/additional-scrape-config.md).

{{< highlight yaml >}}
scrape_configs:
  - job_name: node
    http_sd_config:
      - url: http://prometheus-vcd-sd.prometheus-vcd-sd.svc.cluster.local:9000/sd
    relabel_configs:
      - source_labels: [__meta_vcd_network_internal]
        replacement: "${1}:9100"
        target_label: __address__
      - source_labels: [__meta_vcd_name]
        target_label: instance
        target_label: instance
{{< / highlight >}}

If you want to use the web-config or a configuration file for configuring the service discovery you could extend the `prometheus-vcd-files` configmap, just add a block similar to this one to the `configMapGenerator` mentioned above, it is referencing to files within the same folder as your `kustomization.yml`, after that you got to set the `PROMETHEUS_VCD_CONFIG` environment variable to `/etc/prometheus-vcd-sd/config.json` within the `prometheus-vcd-sd` configmap:

{{< highlight yaml >}}
  - name: prometheus-vcd-sd-files
    behavior: merge
    files:
      - config.json
{{< / highlight >}}
