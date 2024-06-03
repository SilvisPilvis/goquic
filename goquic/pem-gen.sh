#!/usr/bin/env bash
localip=$(ip -4 -o addr show wlp4s0 | awk '{print $4}' | cut -d/ -f1)
# echo "/C=LV/ST=Ropazi/L=Upeslejas/O=Featherworks/OU=Development/CN=${localip}"
# expires in 10 years
openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -sha256 -days 3650 -nodes -subj "/C=LV/ST=Ropazi/L=Upeslejas/O=Featherworks/OU=Development/CN=${localip}"
# openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 3650 -nodes -subj "/C=LV/ST=Ropazi/L=Upeslejas/O=Featherworks/OU=Development/CN=${localip}"