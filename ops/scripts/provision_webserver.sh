#!/usr/bin/env bash
# Sets up the environment config for the webcounter service
# depends on systemd dorp-in unit file already being in /tmp/environment.conf

sudo mkdir -p "/etc/systemd/system/webcounter.service.d/"
sudo cp "/tmp/environment.conf" "/etc/systemd/system/webcounter.service.d/."
sudo systemctl daemon-reload
sudo systemctl restart webcounter.service
