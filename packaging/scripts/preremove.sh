#!/bin/sh
set -e

systemctl stop prometheus-vcd-sd.service || true
systemctl disable prometheus-vcd-sdpi.service || true
