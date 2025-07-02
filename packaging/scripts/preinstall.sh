#!/bin/sh
set -e

if ! getent group prometheus-vcd-sd >/dev/null 2>&1; then
    groupadd --system prometheus-vcd-sd
fi

if ! getent passwd prometheus-vcd-sd >/dev/null 2>&1; then
    useradd --system --create-home --home-dir /var/lib/prometheus-vcd-sd --shell /bin/bash -g prometheus-vcd-sd prometheus-vcd-sd
fi
