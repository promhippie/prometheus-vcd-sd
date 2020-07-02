Bugfix: Fix authentication timeout/invalidation

When running the service it happened that the authentication had been
invalidated or simply timed out, this should be fixed by simply authenticating
the defined user before looping through all the results. At the end also the
disconnect function from the used library gets executed.

https://github.com/promhippie/prometheus-vcd-sd/issues/2
