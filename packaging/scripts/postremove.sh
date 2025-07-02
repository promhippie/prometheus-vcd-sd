#!/bin/sh
set -e

if [ ! -d /var/lib/prometheus-vcd-sd ] && [ ! -d /etc/prometheus-vcd-sd ]; then
    userdel prometheus-vcd-sd 2>/dev/null || true
    groupdel prometheus-vcd-sd 2>/dev/null || true
fi
