#!/bin/sh
set -e

chown -R prometheus-vcd-sd:prometheus-vcd-sd /etc/prometheus-vcd-sd
chown -R prometheus-vcd-sd:prometheus-vcd-sd /var/lib/prometheus-vcd-sd
chmod 750 /var/lib/prometheus-vcd-sd

if [ -d /run/systemd/system ]; then
    systemctl daemon-reload

    if systemctl is-enabled --quiet prometheus-vcd-sd.service; then
        systemctl restart prometheus-vcd-sd.service
    fi
fi
