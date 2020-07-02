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


