# Changelog for 0.2.0

The following sections list the changes for 0.2.0.

## Summary

 * Chg #33: Improvements for automated documentation
 * Chg #34: Integrate new HTTP service discovery handler
 * Chg #30: Add new label for virtual machine ID
 * Chg #35: Integrate standard web config

## Details

 * Change #33: Improvements for automated documentation

   We have added some simple scripts that gets executed by Drone to keep moving documentation
   parts like the available labels or the available environment variables always up to date. No
   need to update the docs related to that manually anymore.

   https://github.com/promhippie/prometheus-vcd-sd/pull/33

 * Change #34: Integrate new HTTP service discovery handler

   We integrated the new HTTP service discovery which have been introduced by Prometheus
   starting with version 2.28. With this new service discovery you can deploy this service
   whereever you want and you are not tied to the Prometheus filesystem anymore.

   https://github.com/promhippie/prometheus-vcd-sd/issues/34

 * Change #30: Add new label for virtual machine ID

   We've added a new label to get the current ID of a virtual machine in the format provided by vCloud
   Director in the form of `urn:vcloud:vm:807b799e-c72f-4592-9acd-ccefefe92720`.

   https://github.com/promhippie/prometheus-vcd-sd/issues/30

 * Change #35: Integrate standard web config

   We integrated the new web config from the Prometheus toolkit which provides a configuration
   for TLS support and also some basic builtin authentication. For the detailed configuration
   you check out the documentation.

   https://github.com/promhippie/prometheus-vcd-sd/issues/35


# Changelog for 0.1.2

The following sections list the changes for 0.1.2.

## Summary

 * Fix #10: Fix nil pointer if vApp doesn't provide a VM
 * Fix #5: Normalize labels for networks
 * Chg #6: Use bingo for development tooling
 * Chg #7: Update Go version and dependencies
 * Chg #8: Drop dariwn/386 release builds

## Details

 * Bugfix #10: Fix nil pointer if vApp doesn't provide a VM

   We have added a check if a vApp really provides children to avoid a panic because of nil pointer
   dereference within the vCD client SDK. Without this fix the service discovery panics on empty
   vApp.

   https://github.com/promhippie/prometheus-vcd-sd/issues/10

 * Bugfix #5: Normalize labels for networks

   We have applied a fix to properly normalize names of networks attached to VMs, before this patch
   the labels could include dashes, which is an invalid label for prometheus.

   https://github.com/promhippie/prometheus-vcd-sd/issues/5

 * Change #6: Use bingo for development tooling

   We switched to use [bingo](github.com/bwplotka/bingo) for fetching development and build
   tools based on fixed defined versions to reduce the dependencies listed within the regular
   go.mod file within this project.

   https://github.com/promhippie/prometheus-vcd-sd/issues/6

 * Change #7: Update Go version and dependencies

   We updated the Go version used to build the binaries within the CI system and beside that in the
   same step we have updated all dependencies ti keep everything up to date.

   https://github.com/promhippie/prometheus-vcd-sd/issues/7

 * Change #8: Drop dariwn/386 release builds

   We dropped the build of 386 builds on Darwin as this architecture is not supported by current Go
   versions anymore.

   https://github.com/promhippie/prometheus-vcd-sd/issues/8


# Changelog for 0.1.1

The following sections list the changes for 0.1.1.

## Summary

 * Fix #2: Fix authentication timeout/invalidation

## Details

 * Bugfix #2: Fix authentication timeout/invalidation

   When running the service it happened that the authentication had been invalidated or simply
   timed out, this should be fixed by simply authenticating the defined user before looping
   through all the results. At the end also the disconnect function from the used library gets
   executed.

   https://github.com/promhippie/prometheus-vcd-sd/issues/2


# Changelog for 0.1.0

The following sections list the changes for 0.1.0.

## Summary

 * Chg #1: Initial release of basic version

## Details

 * Change #1: Initial release of basic version

   Just prepared an initial basic version which could be released to the public.

   https://github.com/promhippie/prometheus-vcd-sd/issues/1


