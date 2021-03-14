# Changelog for unreleased

The following sections list the changes for unreleased.

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


