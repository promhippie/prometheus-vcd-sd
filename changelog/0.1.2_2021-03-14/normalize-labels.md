Bugfix: Normalize labels for networks

We have applied a fix to properly normalize names of networks attached to VMs,
before this patch the labels could include dashes, which is an invalid label for
prometheus.

https://github.com/promhippie/prometheus-vcd-sd/issues/5
