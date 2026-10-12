#!/bin/sh
# Boots a doppler sporocarp: Garage in the background, a bucket and key for it,
# sample posts and photos (cmd/seed; the author is pinwheel / pinwheel), then
# the app on :1323. Everything lives in the container and goes with it.
set -eu

/garage -c /etc/garage.toml server >/tmp/garage.log 2>&1 &

/app/garage-init
set -a
. /shared/s3-credentials.env
set +a

DOPPLER_SESSION_SECRET="$(head -c 32 /dev/urandom | base64)"
export DOPPLER_SESSION_SECRET

/app/seed-binary
exec /app/app-binary
