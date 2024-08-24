# Changelog for 1.0.0

The following sections list the changes for 1.0.0.

## Summary

 * Chg #223: Read secrets form files
 * Enh #223: Update all releated dependencies

## Details

 * Change #223: Read secrets form files

   We have added proper support to load secrets like passwords from files or from
   base64-encoded strings. Just provide the flags or environment variables with a
   DSN formatted string like `file://path/to/file` or `base64://Zm9vYmFy`.

   https://github.com/promhippie/prometheus-vcd-sd/issues/223

 * Enhancement #223: Update all releated dependencies

   We've updated all dependencies to the latest available versions, including more
   current versions of build tools and used Go version to build the binaries. It's
   time to mark a stable release.

   https://github.com/promhippie/prometheus-vcd-sd/issues/223


# Changelog for 0.4.1

The following sections list the changes for 0.4.1.

## Summary

 * Fix #142: Network names can include a space

## Details

 * Bugfix #142: Network names can include a space

   Nnetwork names can include a space, with this fix they are properly sanitized
   and you are able to use these labels.

   https://github.com/promhippie/prometheus-vcd-sd/issues/142


# Changelog for 0.4.0

The following sections list the changes for 0.4.0.

## Summary

 * Enh #127: Improve doucmentation and repo structure

## Details

 * Enhancement #127: Improve doucmentation and repo structure

   We have improved the available documentation pretty hard and we also added
   documentation how to install this service discovery via Helm or Kustomize on
   Kubernetes. Beside that we are testing to build the bundled Kustomize manifests
   now.

   https://github.com/promhippie/prometheus-vcd-sd/pull/127


# Changelog for 0.3.0

The following sections list the changes for 0.3.0.

## Summary

 * Fix #86: Properly normalize label names
 * Chg #90: Refactor build tools and project structure

## Details

 * Bugfix #86: Properly normalize label names

   We have added more character replacements for generating the label names as it
   could contain bad characters depending on the definitions within a vCloud
   Director instance. Now we are replacing `-`, `.` and `,` by `_`.

   https://github.com/promhippie/prometheus-vcd-sd/issues/86

 * Change #90: Refactor build tools and project structure

   To have a unified project structure and build tooling we have integrated the
   same structure we already got within our GitHub exporter.

   https://github.com/promhippie/prometheus-vcd-sd/issues/90


# Changelog for 0.2.0

The following sections list the changes for 0.2.0.

## Summary

 * Chg #30: Add new label for virtual machine ID
 * Chg #33: Improvements for automated documentation
 * Chg #34: Integrate new HTTP service discovery handler
 * Chg #35: Integrate standard web config

## Details

 * Change #30: Add new label for virtual machine ID

   We've added a new label to get the current ID of a virtual machine in the format
   provided by vCloud Director in the form of
   `urn:vcloud:vm:807b799e-c72f-4592-9acd-ccefefe92720`.

   https://github.com/promhippie/prometheus-vcd-sd/issues/30

 * Change #33: Improvements for automated documentation

   We have added some simple scripts that gets executed by Drone to keep moving
   documentation parts like the available labels or the available environment
   variables always up to date. No need to update the docs related to that manually
   anymore.

   https://github.com/promhippie/prometheus-vcd-sd/pull/33

 * Change #34: Integrate new HTTP service discovery handler

   We integrated the new HTTP service discovery which have been introduced by
   Prometheus starting with version 2.28. With this new service discovery you can
   deploy this service whereever you want and you are not tied to the Prometheus
   filesystem anymore.

   https://github.com/promhippie/prometheus-vcd-sd/issues/34

 * Change #35: Integrate standard web config

   We integrated the new web config from the Prometheus toolkit which provides a
   configuration for TLS support and also some basic builtin authentication. For
   the detailed configuration you check out the documentation.

   https://github.com/promhippie/prometheus-vcd-sd/issues/35


# Changelog for 0.1.2

The following sections list the changes for 0.1.2.

## Summary

 * Fix #5: Normalize labels for networks
 * Fix #10: Fix nil pointer if vApp doesn't provide a VM
 * Chg #6: Use bingo for development tooling
 * Chg #7: Update Go version and dependencies
 * Chg #8: Drop dariwn/386 release builds

## Details

 * Bugfix #5: Normalize labels for networks

   We have applied a fix to properly normalize names of networks attached to VMs,
   before this patch the labels could include dashes, which is an invalid label for
   prometheus.

   https://github.com/promhippie/prometheus-vcd-sd/issues/5

 * Bugfix #10: Fix nil pointer if vApp doesn't provide a VM

   We have added a check if a vApp really provides children to avoid a panic
   because of nil pointer dereference within the vCD client SDK. Without this fix
   the service discovery panics on empty vApp.

   https://github.com/promhippie/prometheus-vcd-sd/issues/10

 * Change #6: Use bingo for development tooling

   We switched to use [bingo](github.com/bwplotka/bingo) for fetching development
   and build tools based on fixed defined versions to reduce the dependencies
   listed within the regular go.mod file within this project.

   https://github.com/promhippie/prometheus-vcd-sd/issues/6

 * Change #7: Update Go version and dependencies

   We updated the Go version used to build the binaries within the CI system and
   beside that in the same step we have updated all dependencies ti keep everything
   up to date.

   https://github.com/promhippie/prometheus-vcd-sd/issues/7

 * Change #8: Drop dariwn/386 release builds

   We dropped the build of 386 builds on Darwin as this architecture is not
   supported by current Go versions anymore.

   https://github.com/promhippie/prometheus-vcd-sd/issues/8


# Changelog for 0.1.1

The following sections list the changes for 0.1.1.

## Summary

 * Fix #2: Fix authentication timeout/invalidation

## Details

 * Bugfix #2: Fix authentication timeout/invalidation

   When running the service it happened that the authentication had been
   invalidated or simply timed out, this should be fixed by simply authenticating
   the defined user before looping through all the results. At the end also the
   disconnect function from the used library gets executed.

   https://github.com/promhippie/prometheus-vcd-sd/issues/2


# Changelog for 0.1.0

The following sections list the changes for 0.1.0.

## Summary

 * Chg #1: Initial release of basic version

## Details

 * Change #1: Initial release of basic version

   Just prepared an initial basic version which could be released to the public.

   https://github.com/promhippie/prometheus-vcd-sd/issues/1


