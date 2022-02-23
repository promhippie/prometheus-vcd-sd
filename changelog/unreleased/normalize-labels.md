Bugfix: Properly normalize label names

We have added more character replacements for generating the label names as it
could contain bad characters depending on the definitions within a vCloud
Director instance. Now we are replacing `-`, `.` and `,` by `_`.

https://github.com/promhippie/prometheus-vcd-sd/issues/86
