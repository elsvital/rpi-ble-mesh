#!/bin/bash
sudo apt update
sudo apt install -y python3-gi python3-dbus libdbus-glib-1-dev
pip install --break-system-packages -r requirements.txt
