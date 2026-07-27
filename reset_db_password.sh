#!/bin/bash
sudo -u postgres psql -c "ALTER USER starehe_user WITH PASSWORD 'changeme';"
