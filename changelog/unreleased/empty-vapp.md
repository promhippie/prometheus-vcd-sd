Bugfix: Fix nil pointer if vApp doesn't provide a VM

We have added a check if a vApp really provides children to avoid a panic
because of nil pointer dereference within the vCD client SDK. Without this fix
the service discovery panics on empty vApp.

https://github.com/promhippie/prometheus-vcd-sd/issues/10
